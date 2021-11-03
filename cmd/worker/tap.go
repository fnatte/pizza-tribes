package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const MAX_TAP_STREAK = 12

func (h *handler) handleTap(ctx context.Context, userId string, m *models.ClientMessage_Tap) error {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)
	lotPath := fmt.Sprintf(".lots[\"%s\"]", m.LotId)
	tappedAtPath := fmt.Sprintf("%s.tapped_at", lotPath)
	tapsPath := fmt.Sprintf("%s.taps", lotPath)
	streakPath := fmt.Sprintf("%s.streak", lotPath)
	now := time.Now().UnixNano()

	var lot models.GameState_Lot

	txf := func() error {
		// Get lot
		str, err := internal.RedisJsonGet(h.rdb, ctx, gsKey, lotPath).Result()
		if err != nil {
			return fmt.Errorf("failed to get lot: %w", err)
		}
		if err = protojson.Unmarshal([]byte(str), &lot); err != nil {
			return fmt.Errorf("failed to unmarshal lot: %w", err)
		}

		// Check tap interval
		nextTapAt := lot.TappedAt + (500 * time.Millisecond).Nanoseconds()
		if nextTapAt > now {
			return fmt.Errorf("tapped to soon, next tap at %d", nextTapAt)
		}

		// Set reset time to the beginning of the next hour after TappedAt,
		// and streak time to the hour after that.
		resetTime := time.Unix(0, lot.TappedAt).Add(1 * time.Hour).Truncate(1 * time.Hour)
		resetStreakTime := resetTime.Add(1 * time.Hour)

		// Reset streak if we are past the reset streak time
		if time.Now().After(resetStreakTime) {
			lot.Streak = 0
		}

		// Reset taps if next hour
		if time.Now().After(resetTime) {
			lot.Taps = 0
		}

		if lot.Taps >= 10 {
			return fmt.Errorf("tap is maxed out this hour")
		}

		// Determine what resource to increase and how much
		var incrPath string
		var incrAmount int64
		factor := math.Sqrt(float64(lot.Level+1) * float64(lot.Streak+1))
		switch lot.Building {
		case models.Building_KITCHEN:
			incrPath = ".resources.pizzas"
			incrAmount = int64(math.Round(80*factor/5) * 5)
		case models.Building_SHOP:
			incrPath = ".resources.coins"
			incrAmount = int64(math.Round(35*factor/5) * 5)
		default:
			return fmt.Errorf("this building cannot be tapped")
		}

		lot.Taps = lot.Taps + 1

		if lot.Taps == 10 {
			lot.Streak = lot.Streak + 1
		}

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// Update tapped_at to now
			err = internal.RedisJsonSet(pipe, ctx, gsKey, tappedAtPath, int64(now)).Err()
			if err != nil {
				return fmt.Errorf("failed to set tapped_at: %w", err)
			}

			// Update taps
			err = internal.RedisJsonSet(pipe, ctx, gsKey, tapsPath, lot.Taps).Err()
			if err != nil {
				return fmt.Errorf("failed to set taps: %w", err)
			}

			// Update streak
			err = internal.RedisJsonSet(pipe, ctx, gsKey, streakPath, lot.Streak).Err()
			if err != nil {
				return fmt.Errorf("failed to set streak: %w", err)
			}

			// Increase the resource
			err = internal.RedisJsonNumIncrBy(
				pipe, ctx, gsKey, incrPath, incrAmount).Err()
			if err != nil {
				return fmt.Errorf("failed to increase resource: %w", err)
			}

			return nil
		})

		return err
	}

	mutex := h.rdb.NewMutex("lock:" + gsKey)
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to obtain lock: %w", err)
	}
	err2 := txf()
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}
	if err2 != nil {
		return fmt.Errorf("failed to handle tap: %w", err2)
	}

	if err := h.sendTapUpdate(ctx, userId, m.LotId, &lot, now); err != nil {
		return fmt.Errorf("failed to send tap update: %w", err)
	}

	return nil
}

func (h *handler) sendTapUpdate(ctx context.Context, userId string, lotId string, lot *models.GameState_Lot, tappedAt int64) error {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)
	path := ".resources"

	s, err := h.rdb.JsonGet(ctx, gsKey, path).Result()
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	res := models.GameState_Resources{}
	protojson.Unmarshal([]byte(s), &res)
	if err != nil {
		return fmt.Errorf("failed to unmarshal resources: %w", err)
	}

	lotsPatch := map[string]*models.GameStatePatch_LotPatch{}
	lotsPatch[lotId] = &models.GameStatePatch_LotPatch{
		Building: lot.Building,
		TappedAt: tappedAt,
		Level:    lot.Level,
		Taps:     lot.Taps,
		Streak:   lot.Streak,
	}

	return h.send(ctx, userId, &models.ServerMessage{
		Id: xid.New().String(),
		Payload: &models.ServerMessage_StateChange{
			StateChange: &models.GameStatePatch{
				Lots: lotsPatch,
				Resources: &models.GameStatePatch_ResourcesPatch{
					Coins:  &wrapperspb.Int32Value{Value: res.Coins},
					Pizzas: &wrapperspb.Int32Value{Value: res.Pizzas},
				},
			},
		},
	})
}
