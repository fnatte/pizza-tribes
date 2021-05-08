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

func (h *handler) handleSteal(ctx context.Context, senderId string, m *internal.ClientMessage_Steal) error {
	gsKeyThief := fmt.Sprintf("user:%s:gamestate", senderId)

	var gsThief internal.GameState

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

	txf := func(tx *redis.Tx) error {
		// Get game state of thief
		str, err := internal.RedisJsonGet(tx, ctx, gsKeyThief, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		err = gsThief.LoadProtoJson([]byte(str))
		if err != nil {
			return err
		}

		// Validate game state of thief
		if gsThief.Population == nil || gsThief.Population.Thieves < m.Amount {
			return errors.New("no enough thieves")
		}

		arrivalAt := internal.CalculateArrivalTime(
			gsThief.TownX, gsThief.TownY,
			m.X, m.Y,
			internal.ThiefSpeed)

		travel := internal.Travel{
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

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// Decrease thieves in town population of sending town
			_, err := internal.RedisJsonNumIncrBy(
				pipe, ctx, gsKeyThief,
				".population.thieves",
				int64(-travel.Thieves)).Result()
			if err != nil {
				return fmt.Errorf("failed to decrease thieves of sender: %w", err)
			}

			if err = internal.RedisJsonArrAppend(pipe, ctx, gsKeyThief,
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

	if err = h.rdb.Watch(ctx, txf, gsKeyThief); err != nil {
		return err
	}

	h.sendFullStateUpdate(ctx, senderId)

	return nil
}
