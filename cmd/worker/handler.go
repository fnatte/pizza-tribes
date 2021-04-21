package main

import (
	"context"

	"github.com/fnatte/pizza-mouse/internal"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type handler struct {
	rdb *redis.Client
}

func (h *handler) Handle(ctx context.Context, senderId string, m *internal.ClientMessage) {
	switch x := m.Type.(type) {
	case *internal.ClientMessage_Tap_:
		log.Info().
			Str("senderId", senderId).
			Int32("Amount", x.Tap.Amount).
			Msg("Received message")
	default:
		log.Info().Str("senderId", senderId).Msg("Received message")
	}

	h.send(ctx, senderId)
}

func (h *handler) send(ctx context.Context, senderId string) {
	h.rdb.RPush(ctx, "wsout", &internal.OutgoingMessage{
		ReceiverId: senderId,
		Body:       "{ \"resources\": { \"coins\": 2, \"pizzas\": 2 } }",
	})
}
