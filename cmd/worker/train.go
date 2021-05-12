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

func (h *handler) handleTrain(ctx context.Context, senderId string, m *internal.ClientMessage_Train) {
	log.Info().
		Str("senderId", senderId).
		Interface("Education", m.Education).
		Int32("Amount", m.Amount).
		Msg("Received train message")

	gameStateKey := fmt.Sprintf("user:%s:gamestate", senderId)

	var gs internal.GameState

	txf := func(tx *redis.Tx) error {
		// Get current game state
		b, err := internal.RedisJsonGet(tx, ctx, gameStateKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		err = gs.LoadProtoJson([]byte(b))
		if err != nil {
			return err
		}

		eduInfo := internal.FullGameData.Educations[int32(m.Education)]
		trainTime := int64(eduInfo.TrainTime)
		cost := eduInfo.Cost

		if gs.Population.Uneducated < m.Amount {
			return errors.New("Too few uneducated")
		}

		if gs.Resources.Coins < cost {
			return errors.New("Not enough coins")
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			_, err := internal.RedisJsonNumIncrBy(
				pipe,
				ctx,
				gameStateKey,
				".population.uneducated",
				int64(-m.Amount)).Result()
			if err != nil {
				log.Error().Err(err).Msg("Failed to decrease uneducated")
				return err
			}

			_, err = internal.RedisJsonNumIncrBy(
				pipe, ctx, gameStateKey,
				".resources.coins",
				int64(-cost)).Result()
			if err != nil {
				log.Error().Err(err).Msg("Failed to decrease coins")
				return err
			}

			training := internal.Training{
				CompleteAt: time.Now().UnixNano() + trainTime*1e9,
				Education:  m.Education,
				Amount:     m.Amount,
			}

			b, err := protojson.Marshal(&training)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal training")
				return err
			}

			internal.RedisJsonArrAppend(
				pipe,
				ctx,
				fmt.Sprintf("user:%s:gamestate", senderId),
				".trainingQueue",
				b,
			)

			return nil
		})
		return err
	}

	err := h.rdb.Watch(ctx, txf, gameStateKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to train")
		return
	}

	h.fetchAndUpdateTimestamp(ctx, senderId)
	h.sendFullStateUpdate(ctx, senderId)
}


