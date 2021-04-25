package main

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-mouse/internal"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

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
	default:
		log.Info().Str("senderId", senderId).Msg("Received message")
	}

	/*
	err := h.send(ctx, senderId, &internal.ServerMessage{
		Id: xid.New().String(),
		Payload: &internal.ServerMessage_StateChange{
			StateChange: &internal.GameStatePatch{
				Resources: &internal.GameStatePatch_ResourcesPatch{
					Coins:  internal.NewInt64(10),
					Pizzas: internal.NewInt64(2),
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send messsage back to client")
	}
	*/
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
