package main

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/cmd/api/ws"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

type wsHandler struct {
	rc    redis.RedisClient
	world *internal.WorldService
	gsRepo persist.GameStateRepository
	marketRepo persist.MarketRepository
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

	gs := models.NewGameState()

	log.Info().Str("userId", c.UserId()).Msg("Client connected")

	gsKey := fmt.Sprintf("user:%s:gamestate", c.UserId())
	s, err := h.rc.JsonGet(ctx, gsKey, ".").Result()
	if err != nil {
		if err != redis.Nil {
			return fmt.Errorf("failed to get gamestate: %w", err)
		}

		// Initialize game state for user
		err = h.gsRepo.Save(ctx, c.UserId(), gs)
		if err != nil {
			return fmt.Errorf("failed to initialize gamestate: %w", err)
		}
		log.Info().Msg("Initilized new game state for user")
	} else {
		if err = protojson.Unmarshal([]byte(s), gs); err != nil {
			return fmt.Errorf("failed to unmarshal gamestate: %w", err)
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
		log.Info().
			Int32("x", gs.TownX).
			Int32("y", gs.TownY).
			Msg("Town acquired")
	}

	// Make sure the user is enqueued for updates
	_, err = internal.SetNextUpdate(h.rc, ctx, c.UserId(), gs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure user updates")
		return err
	}

	worldState, err := h.world.GetState(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get world state")
		return err
	}

	globalDemandScore, err := h.marketRepo.GetGlobalDemandScore(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get world state")
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
		log.Debug().Msg("Sent user message")
	})()

	go (func() {
		msg := gs.ToServerMessage()
		b, err := protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send init game state patch")
			return
		}

		c.Send(b)

		msg = internal.CalculateStats(gs, globalDemandScore, worldState).ToServerMessage()
		b, err = protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send init stats")
			return
		}
		c.Send(b)

		log.Debug().Msg("Sent init game state and stats")
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
		log.Debug().Msg("Sent init reports")
	})()

	return nil
}
