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
	"github.com/rs/zerolog/log"
)

func (h *handler) handleUpgrade(ctx context.Context, senderId string, m *models.ClientMessage_UpgradeBuilding) error {
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

		for _, constr := range gs.ConstructionQueue {
			if constr.LotId == m.LotId {
				return errors.New("the building is under construction")
			}
		}

		lot := gs.Lots[m.LotId]
		if lot == nil {
			return errors.New("no building in lot")
		}

		buildingInfo := internal.FullGameData.Buildings[int32(lot.Building)]
		if int(lot.Level)+1 >= len(buildingInfo.LevelInfos) {
			return errors.New("building already max level")
		}
		cost := buildingInfo.LevelInfos[lot.Level+1].Cost
		constructionTime := buildingInfo.LevelInfos[lot.Level+1].ConstructionTime

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
				log.Error().Err(err).Msg("Failed to decrease coins")
				return err
			}

			construction := models.Construction{
				CompleteAt: completeAt,
				LotId:      m.LotId,
				Building:   lot.Building,
				Level:      lot.Level + 1,
			}

			b, err := protojson.Marshal(&construction)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal construction")
				return err
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
		return fmt.Errorf("failed to handle upgrade: %w", err2)
	}

	internal.SetNextUpdate(h.rdb, ctx, senderId, &gs)

	h.sendFullStateUpdate(ctx, senderId)

	return nil
}
