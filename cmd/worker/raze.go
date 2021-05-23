package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
)

func (h *handler) handleRazeBuilding(ctx context.Context, senderId string, m *models.ClientMessage_RazeBuilding) error {
	gsKey := fmt.Sprintf("user:%s:gamestate", senderId)

	var gs models.GameState

	txf := func(tx *redis.Tx) error {
		// Get current game state
		s, err := internal.RedisJsonGet(tx, ctx, gsKey, ".").Result()
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

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			_, err := internal.RedisJsonNumIncrBy(
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

			err = internal.RedisJsonArrAppend(
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

	err := h.rdb.Watch(ctx, txf, gsKey)
	if err != nil {
		return fmt.Errorf("failed to place on construction queue: %w", err)
	}

	internal.SetNextUpdate(h.rdb, ctx, senderId, &gs)

	h.sendFullStateUpdate(ctx, senderId)

	return nil
}