package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (h *handler) handleTap(ctx context.Context, userId string, m *models.ClientMessage_Tap) error {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)
	lotPath := fmt.Sprintf(".lots[\"%s\"]", m.LotId)
	popPath := ".population"
	tappedAtPath := fmt.Sprintf("%s.tapped_at", lotPath)
	now := time.Now().UnixNano()

	var lot models.GameState_Lot
	var population models.GameState_Population

	txf := func(tx *redis.Tx) error {
		// Get lot
		str, err := internal.RedisJsonGet(h.rdb, ctx, gsKey, lotPath).Result()
		if err != nil {
			return fmt.Errorf("failed to get lot: %w", err)
		}
		if err = protojson.Unmarshal([]byte(str), &lot); err != nil {
			return fmt.Errorf("failed to unmarshal lot: %w", err)
		}

		// Get population
		str, err = internal.RedisJsonGet(h.rdb, ctx, gsKey, popPath).Result()
		if err != nil {
			return fmt.Errorf("failed to get population: %w", err)
		}
		if err = protojson.Unmarshal([]byte(str), &population); err != nil {
			return fmt.Errorf("failed to unmarshal population: %w", err)
		}

		// Determine what resource to increase and how much
		var incrPath string
		var incrAmount int64
		switch lot.Building {
		case models.Building_KITCHEN:
			incrPath = ".resources.pizzas"
			incrAmount = 80 * int64(internal.CountTownPopulation(&population))
		case models.Building_SHOP:
			incrPath = ".resources.coins"
			incrAmount = 35 * int64(internal.CountTownPopulation(&population))
		default:
			return fmt.Errorf("this building cannot be tapped")
		}

		nextTapAt := lot.TappedAt + (15 * time.Minute).Nanoseconds()

		if nextTapAt > now {
			return fmt.Errorf("tapped to soon, next tap at %d", nextTapAt)
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// Update tapped_at to now
			err = internal.RedisJsonSet(pipe, ctx, gsKey, tappedAtPath, int64(now)).Err()
			if err != nil {
				return fmt.Errorf("failed to set tapped_at: %w", err)
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

	if err := h.rdb.Watch(ctx, txf, gsKey); err != nil {
		return err
	}

	if err := h.sendTapUpdate(ctx, userId, m.LotId, lot.Building, lot.Level, now); err != nil {
		return fmt.Errorf("failed to send tap update: %w", err)
	}

	return nil
}

func (h *handler) sendTapUpdate(ctx context.Context, userId string, lotId string, building models.Building, level int32, tappedAt int64) error {
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
		Building: building,
		TappedAt: tappedAt,
		Level: level,
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
