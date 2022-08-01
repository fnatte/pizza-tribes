package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/rs/zerolog/log"
)

func (h *handler) handleSteal(ctx context.Context, senderId string, m *models.ClientMessage_Steal) error {
	gsKeyThief := fmt.Sprintf("user:%s:gamestate", senderId)

	var gsThief models.GameState

	// Validate target town
	worldEntry, err := h.world.GetEntryXY(ctx, int(m.X), int(m.Y))
	if err != nil {
		return err
	}
	town := worldEntry.GetTown()
	if town == nil {
		return fmt.Errorf("no town at %d, %d", m.X, m.Y)
	}
	if town.UserId == senderId {
		return errors.New("can't steal from own town")
	}

	txf := func() error {
		// Get game state of thief
		s, err := redis.RedisJsonGet(h.rdb, ctx, gsKeyThief, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), &gsThief); err != nil {
			return err
		}

		// Validate game state of thief
		e := internal.CountTownPopulationEducations(&gsThief)
		if e[models.Education_THIEF] < m.Amount {
			return errors.New("not enough thieves")
		}

		arrivalAt := internal.CalculateArrivalTime(
			gsThief.TownX, gsThief.TownY,
			m.X, m.Y,
			internal.ThiefSpeed)

		travel := models.Travel{
			ArrivalAt:    arrivalAt,
			DestinationX: m.X,
			DestinationY: m.Y,
			Returning:    false,
			Thieves:      m.Amount,
			Coins:        0,
		}

		b, err := protojson.Marshal(&travel)
		if err != nil {
			return fmt.Errorf("failed to marshal travel: %w", err)
		}

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if err = redis.RedisJsonArrAppend(pipe, ctx, gsKeyThief,
				".travelQueue", b).Err(); err != nil {
				return err
			}

			log.Info().
				Int32("thieves", travel.Thieves).
				Time("arrivalAt", time.Unix(0, travel.ArrivalAt)).
				Msg("Steal dispatched")

			return nil
		})

		return err
	}

	mutex := h.rdb.NewMutex("lock:" + gsKeyThief)
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to obtain lock: %w", err)
	}
	err2 := txf()
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}
	if err2 != nil {
		return fmt.Errorf("failed to handle steal: %w", err2)
	}

	h.sendFullStateUpdate(ctx, senderId)

	_, err = internal.SchedulePushNotification(ctx, h.rdb, &messaging.Message{
		Data: map[string]string{
			"userId": town.UserId,
		},
		Notification: &messaging.Notification{
			Title: "Incoming heist",
			Body:  "Boss! We have heard rumours about a planned heist on our town.",
		},
	}, time.Now())
	if err != nil {
		log.Error().Err(err).Msg("Failed to schedule push notification for a heist")
	}

	return nil
}
