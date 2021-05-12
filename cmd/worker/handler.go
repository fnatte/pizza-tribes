package main

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

type handler struct {
	rdb internal.RedisClient
	world *internal.WorldService
}

func (h *handler) Handle(ctx context.Context, senderId string, m *internal.ClientMessage) {
	var err error

	switch x := m.Type.(type) {
	case *internal.ClientMessage_Tap_:
		err = h.handleTap(ctx, senderId, x.Tap)
	case *internal.ClientMessage_ConstructBuilding_:
		h.handleConstructBuilding(ctx, senderId, x.ConstructBuilding)
	case *internal.ClientMessage_Train_:
		h.handleTrain(ctx, senderId, x.Train)
	case *internal.ClientMessage_Steal_:
		err = h.handleSteal(ctx, senderId, x.Steal)
	case *internal.ClientMessage_ReadReport_:
		err = h.handleReadReport(ctx, senderId, x.ReadReport)
	case *internal.ClientMessage_UpgradeBuilding_:
		err = h.handleUpgrade(ctx, senderId, x.UpgradeBuilding)
	default:
		log.Info().Str("senderId", senderId).Msg("Received message")
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to handle message")
	}
}

func (h *handler) fetchAndUpdateTimestamp(ctx context.Context, userId string) (int64, error) {
	b, err := h.rdb.JsonGet(ctx, fmt.Sprintf("user:%s:gamestate", userId), ".").Result()
	if err != nil {
		return 0, err
	}

	gs := internal.GameState{}
	gs.LoadProtoJson([]byte(b))
	if err != nil {
		return 0, err
	}

	return internal.SetNextUpdate(h.rdb, ctx, userId, &gs)
}

func (h *handler) sendFullStateUpdate(ctx context.Context, senderId string) {
	s, err := h.rdb.JsonGet(ctx, fmt.Sprintf("user:%s:gamestate", senderId), ".").Result()
	if err != nil {
		log.Error().Err(err).Msg("Failed to send full state update")
		return
	}

	gs := internal.GameState{}
	gs.LoadProtoJson([]byte(s))
	if err != nil {
		log.Error().Err(err).Msg("Failed to send full state update")
		return
	}

	msg := gs.ToStateChangeMessage()
	err = h.send(ctx, senderId, msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send state update")
		return
	}

	msg = internal.CalculateStats(&gs).ToServerMessage()
	err = h.send(ctx, senderId, msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send stats message")
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
