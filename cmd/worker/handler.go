package main

import (
	"context"

	"github.com/fnatte/pizza-mouse/internal"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
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

	err := h.send(ctx, senderId, &internal.ServerMessage{
		Id: xid.New().String(),
		Payload: &internal.ServerMessage_StateChange_{
			StateChange: &internal.ServerMessage_StateChange{
				Resources: &internal.ServerMessage_ResourcesPatch{
					Coins: NewInt64(10),
					Pizzas: NewInt64(2),
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send messsage back to client")
	}
}

func (h *handler) send(ctx context.Context, senderId string, m *internal.ServerMessage) error {
	b, err := protojson.Marshal(m)
	if err != nil {
		return err
	}

	h.rdb.RPush(ctx, "wsout", &internal.OutgoingMessage{
		ReceiverId: senderId,
		Body: string(b),
	})

	return nil
}
