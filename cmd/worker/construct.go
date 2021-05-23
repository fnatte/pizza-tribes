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
	"github.com/rs/zerolog/log"
)

func (h *handler) handleConstructBuilding(ctx context.Context, senderId string, m *models.ClientMessage_ConstructBuilding) {
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

		buildingInfo := internal.FullGameData.Buildings[int32(m.Building)]
		buildingCount := internal.CountBuildings(&gs)
		buildingConstrCount := internal.CountBuildingsUnderConstruction(&gs)

		cost := buildingInfo.LevelInfos[0].Cost
		constructionTime := buildingInfo.LevelInfos[0].ConstructionTime

		// Can only build at empty lot
		if gs.Lots[m.LotId] != nil {
			return errors.New("Lot must be empty")
		}
		for _, constr := range gs.ConstructionQueue {
			if constr.LotId == m.LotId {
				return errors.New("Already constructing at lot")
			}
		}

		// The first building of all types except markting hq and research institute
		// are free and built 100 times faster
		if m.Building != models.Building_MARKETINGHQ &&
		   m.Building != models.Building_RESEARCH_INSTITUTE &&
			buildingCount[int32(m.Building)]+buildingConstrCount[int32(m.Building)] == 0 {
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

			construction := models.Construction{
				CompleteAt: completeAt,
				LotId:      m.LotId,
				Building:   m.Building,
			}

			b, err := protojson.Marshal(&construction)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal construction")
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
