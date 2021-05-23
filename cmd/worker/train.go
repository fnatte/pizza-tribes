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

func (h *handler) handleTrain(ctx context.Context, senderId string, m *models.ClientMessage_Train) error {
	log.Info().
		Str("senderId", senderId).
		Interface("Education", m.Education).
		Int32("Amount", m.Amount).
		Msg("Received train message")

	gameStateKey := fmt.Sprintf("user:%s:gamestate", senderId)

	var gs models.GameState

	txf := func() error {
		// Get current game state
		s, err := internal.RedisJsonGet(h.rdb, ctx, gameStateKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
			return err
		}

		if m.Amount <= 0 {
			return errors.New("Amount must be greater than 0")
		}

		if gs.Population.Uneducated < m.Amount {
			return errors.New("Too few uneducated")
		}

		eduInfo := internal.FullGameData.Educations[int32(m.Education)]
		trainTime := int64(eduInfo.TrainTime)
		cost := eduInfo.Cost * m.Amount

		if gs.Resources.Coins < cost {
			return errors.New("Not enough coins")
		}

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
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

			training := models.Training{
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

	mutex := h.rdb.NewMutex("lock:" + gameStateKey)
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to obtain lock: %w", err)
	}
	err2 := txf()
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}
	if err2 != nil {
		return fmt.Errorf("failed to handle train: %w", err2)
	}

	h.fetchAndUpdateTimestamp(ctx, senderId)
	h.sendFullStateUpdate(ctx, senderId)

	return nil
}
