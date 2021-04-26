package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fnatte/pizza-mouse/internal"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

var trainTime = map[internal.Education]int64{
	internal.Education_CHEF: 10,
	internal.Education_SALESMOUSE: 15,
	internal.Education_GUARD: 20,
	internal.Education_THIEF: 30,
}

type handler struct {
	rdb internal.RedisClient
}

func (h *handler) Handle(ctx context.Context, senderId string, m *internal.ClientMessage) {
	switch x := m.Type.(type) {
	case *internal.ClientMessage_Tap_:
		log.Info().
			Str("senderId", senderId).
			Int32("Amount", x.Tap.Amount).
			Msg("Received message")
	case *internal.ClientMessage_ConstructBuilding_:
		h.handleConstructBuilding(ctx, senderId, x.ConstructBuilding)
	case *internal.ClientMessage_Train_:
		h.handleTrain(ctx, senderId, x.Train)
	default:
		log.Info().Str("senderId", senderId).Msg("Received message")
	}
}

func (h *handler) handleConstructBuilding(ctx context.Context, senderId string, m *internal.ClientMessage_ConstructBuilding) {
	log.Info().
		Str("senderId", senderId).
		Str("Building", m.Building).
		Str("LotId", m.LotId).
		Msg("Received message")
	err := h.rdb.JsonSet(
		ctx,
		fmt.Sprintf("user:%s:gamestate", senderId),
		fmt.Sprintf(".lots[\"%s\"]", m.LotId),
		fmt.Sprintf("{ \"building\": \"%s\" }", m.Building)).Err()
	if err != nil {
		log.Error().Err(err).Msg("Failed to handle construct building message")
	}
	h.sendFullStateUpdate(ctx, senderId)
}

func (h *handler) handleTrain(ctx context.Context, senderId string, m *internal.ClientMessage_Train) {
	log.Info().
		Str("senderId", senderId).
		Interface("Education", m.Education).
		Int32("Amount", m.Amount).
		Msg("Received train message")

	gameStateKey := fmt.Sprintf("user:%s:gamestate", senderId)

	txf := func(tx *redis.Tx) error {
		n, err := internal.RedisJsonGet(tx, ctx, gameStateKey, ".population.unemployed").Int64()
		if err != nil && err != redis.Nil {
			return err
		}
		if n < int64(m.Amount) {
			return errors.New("Too few unemployed")
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			_, err := internal.RedisJsonNumIncrBy(
				pipe,
				ctx,
				gameStateKey,
				".population.unemployed",
				int64(-m.Amount)).Result()
			if err != nil {
				log.Error().Err(err).Msg("Failed to decrease unemployed")
				return err
			}

			training := internal.Training{
				CompleteAt: time.Now().UnixNano() + trainTime[m.Education] * 1e9,
				Education: m.Education,
				Amount: m.Amount,
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

	h.sendFullStateUpdate(ctx, senderId)
}

func (h *handler) sendFullStateUpdate(ctx context.Context, senderId string) {
	s, err := h.rdb.JsonGet(ctx, fmt.Sprintf("user:%s:gamestate", senderId), ".").Result()
	if err != nil {
		log.Error().Err(err).Msg("Failed to send full state update")
		return
	}

	gs := internal.GameState{}
	err = protojson.Unmarshal([]byte(s), &gs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send full state update")
		return
	}

	msg := gs.ToStateChangeMessage()
	err = h.send(ctx, senderId, msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send full state update")
		return
	}
}

func (h *handler) send(ctx context.Context, senderId string, m *internal.ServerMessage) error {
	b, err := protojson.Marshal(m)
	if err != nil {
		return err
	}

	h.rdb.RPush(ctx, "wsout", &internal.OutgoingMessage{
		ReceiverId: senderId,
		Body:       string(b),
	})

	return nil
}
