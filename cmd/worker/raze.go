package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
)

func (h *handler) handleRazeBuilding(ctx context.Context, senderId string, m *models.ClientMessage_RazeBuilding) error {
	if !internal.IsValidLotId(m.LotId) {
		return errors.New("Invalid lot id")
	}

	gsKey := fmt.Sprintf("user:%s:gamestate", senderId)

	var gs models.GameState

	txf := func() error {
		// Get current game state
		s, err := redis.RedisJsonGet(h.rdb, ctx, gsKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
			return err
		}

		// Can only raze existing buildings
		if gs.Lots[m.LotId] == nil {
			return errors.New("Lot was already empty")
		}
		for _, constr := range gs.ConstructionQueue {
			if constr.LotId == m.LotId {
				return errors.New("Already constructing at lot")
			}
		}

		lot := gs.Lots[m.LotId]

		buildingInfo := internal.FullGameData.Buildings[int32(lot.Building)]
		constructionTime := buildingInfo.LevelInfos[lot.Level].ConstructionTime * 2
		cost := buildingInfo.LevelInfos[lot.Level].Cost / 2

		if gs.Resources.Coins < cost {
			return errors.New("Not enough coins")
		}

		// Calculate when this construction will be completed.
		// If there's already already something being constructed, this building will
		// be started at the end of previous one. If there's nothing in queue, it can
		// be started immediately (time.Now()).
		timeOffset := time.Now().UnixNano()
		if n := len(gs.ConstructionQueue); n > 0 {
			timeOffset = gs.ConstructionQueue[n-1].CompleteAt
		}
		completeAt := timeOffset + int64(constructionTime)*1e9

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			_, err := redis.RedisJsonNumIncrBy(
				pipe, ctx, gsKey,
				".resources.coins",
				int64(-cost)).Result()
			if err != nil {
				return fmt.Errorf("failed to decrease coins: %w", err)
			}

			construction := models.Construction{
				CompleteAt: completeAt,
				LotId:      m.LotId,
				Building:   lot.Building,
				Level:      lot.Level,
				Razing:     true,
			}

			b, err := protojson.Marshal(&construction)
			if err != nil {
				return fmt.Errorf("failed to marshal construction: %w", err)
			}

			err = redis.RedisJsonArrAppend(
				pipe,
				ctx,
				fmt.Sprintf("user:%s:gamestate", senderId),
				".constructionQueue",
				b,
			).Err()
			if err != nil {
				return err
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
		return fmt.Errorf("failed to place on construction queue: %w", err2)
	}

	internal.SetNextUpdate(h.rdb, ctx, senderId, &gs)

	h.sendFullStateUpdate(ctx, senderId)

	return nil
}

func (h *handler) handleCancelRazeBuilding(ctx context.Context, senderId string, m *models.ClientMessage_CancelRazeBuilding) error {
	if !internal.IsValidLotId(m.LotId) {
		return errors.New("Invalid lot id")
	}

	gsKey := fmt.Sprintf("user:%s:gamestate", senderId)

	var gs models.GameState

	// TODO: adjust the constructions times of succeeding items

	txf := func() error {
		// Get current game state
		s, err := redis.RedisJsonGet(h.rdb, ctx, gsKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
			return err
		}

		// Make sure there is a "raze building" to cancel
		index := -1
		for i, constr := range gs.ConstructionQueue {
			if constr.LotId == m.LotId {
				index = i
				break
			}
		}
		if index == -1 {
			return fmt.Errorf("could not find a raze building at that lot to cancel")
		}

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if err = redis.RedisJsonArrPop(pipe, ctx, gsKey,
				".constructionQueue", index).Err(); err != nil {
				return err
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
		return fmt.Errorf("failed to cancel raze building: %w", err2)
	}

	internal.SetNextUpdate(h.rdb, ctx, senderId, &gs)

	h.sendFullStateUpdate(ctx, senderId)

	return nil
}
