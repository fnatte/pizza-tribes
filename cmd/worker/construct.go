package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func (h *handler) handleConstructBuilding(ctx context.Context, senderId string, m *internal.ClientMessage_ConstructBuilding) {
	gsKey := fmt.Sprintf("user:%s:gamestate", senderId)

	var gs internal.GameState

	txf := func(tx *redis.Tx) error {
		// Get current game state
		b, err := internal.RedisJsonGet(tx, ctx, gsKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		err = gs.LoadProtoJson([]byte(b))
		if err != nil {
			return err
		}

		buildingInfo := internal.FullGameData.Buildings[int32(m.Building)]
		buildingCount := internal.CountBuildings(&gs)
		buildingConstrCount := internal.CountBuildingsUnderConstruction(&gs)

		cost := buildingInfo.Cost
		constructionTime := buildingInfo.ConstructionTime

		// The first building of each type is free and built 100 times faster
		if buildingCount[int32(m.Building)]+buildingConstrCount[int32(m.Building)] == 0 {
			cost = 0
			constructionTime = int32(float64(constructionTime)/100.0) + 1
		}

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
				log.Error().Err(err).Msg("Failed to decrease coins")
				return err
			}

			construction := internal.Construction{
				CompleteAt: completeAt,
				LotId:      m.LotId,
				Building:   m.Building,
			}

			b, err := protojson.Marshal(&construction)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal training")
				return err
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
		log.Error().Err(err).Msg("Failed to place on construction queue")
		return
	}

	internal.SetNextUpdate(h.rdb, ctx, senderId, &gs)

	h.sendFullStateUpdate(ctx, senderId)
}
