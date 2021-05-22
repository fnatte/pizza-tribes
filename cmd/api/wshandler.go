package main

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/cmd/api/ws"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"github.com/fnatte/pizza-tribes/internal/protojson"
)

type wsHandler struct {
	rc    internal.RedisClient
	world *internal.WorldService
}

func (h *wsHandler) HandleMessage(ctx context.Context, m []byte, c *ws.Client) {
	log.Debug().Str("userId", c.UserId()).Msg("Received message")
	err := h.rc.RPush(ctx, "wsin", &internal.IncomingMessage{
		SenderId: c.UserId(),
		Body:     string(m),
	}).Err()
	if err != nil {
		log.Error().Err(err).Msg("Error when pushing incoming message to redis")
	}
}
func (h *wsHandler) HandleInit(ctx context.Context, c *ws.Client) error {
	// Please excuse me for this looong func :|

	// Get username
	userKey := fmt.Sprintf("user:%s", c.UserId())
	username, err := h.rc.HGet(ctx, userKey, "username").Result()
	if err != nil {
		return fmt.Errorf("failed to get username: %w", err)
	}

	gs := models.GameState{
		Population: &models.GameState_Population{},
		Resources:  &models.GameState_Resources{},
		Lots:       map[string]*models.GameState_Lot{},
	}

	log.Info().Str("userId", c.UserId()).Msg("Client connected")
	gsKey := fmt.Sprintf("user:%s:gamestate", c.UserId())
	s, err := h.rc.JsonGet(ctx, gsKey, ".").Result()
	if err != nil {
		if err != redis.Nil {
			return err
		}
		b, err := protojson.MarshalOptions{
			EmitUnpopulated: true,
		}.Marshal(&gs)
		if err != nil {
			return err
		}
		err = h.rc.JsonSet(ctx, gsKey, ".", string(b)).Err()
		if err != nil {
			return err
		}
		log.Info().Msg("Initilized new game state for user")
	} else {
		if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
			return err
		}
	}

	// Make sure the user has town in world
	if gs.TownX == 0 && gs.TownY == 0 {
		x, y, err := h.world.AcquireTown(ctx, c.UserId())
		if err != nil {
			return fmt.Errorf("failed to acquire town: %w", err)
		}
		gs.TownX = int32(x)
		gs.TownY = int32(y)
		log.Info().Msg("Town acquired")
	}

	// Make sure the user is enqueued for updates
	_, err = internal.SetNextUpdate(h.rc, ctx, c.UserId(), &gs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure user updates")
		return err
	}

	go (func() {
		msg := &models.ServerMessage{
			Id: xid.New().String(),
			Payload: &models.ServerMessage_User_{
				User: &models.ServerMessage_User{
					Username: username,
				},
			},
		}
		b, err := protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send init game state patch")
			return
		}

		c.Send(b)
		log.Info().Msg("Sent user message")
	})()

	go (func() {
		msg := gs.ToStateChangeMessage()
		b, err := protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send init game state patch")
			return
		}

		c.Send(b)

		msg = internal.CalculateStats(&gs).ToServerMessage()
		b, err = protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send init stats")
			return
		}
		c.Send(b)

		log.Info().Msg("Sent init game state and stats")
	})()

	// Send reports
	go (func() {
		r, err := internal.GetReports(ctx, h.rc, c.UserId())
		if err != nil {
			log.Error().Err(err).Msg("Failed to send inital reports")
			return
		}

		msg := &models.ServerMessage{
			Id: xid.New().String(),
			Payload: &models.ServerMessage_Reports_{
				Reports: &models.ServerMessage_Reports{
					Reports: r,
				},
			},
		}
		b, err := protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send inital reports")
			return
		}
		c.Send(b)
		log.Info().Msg("Sent init reports")
	})()

	return nil
}

