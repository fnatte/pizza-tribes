package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func (h *handler) handleReschoolMouse(ctx context.Context, senderId string, m *models.ClientMessage_ReschoolMouse) error {
	log.Info().
		Str("senderId", senderId).
		Str("MouseId", m.MouseId).
		Msg("Received reschool mouse message")

	gsKey := fmt.Sprintf("user:%s:gamestate", senderId)

	var gs models.GameState

	txf := func() error {
		// Get current game state
		s, err := internal.RedisJsonGet(h.rdb, ctx, gsKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
			return err
		}

		var mouse *models.Mouse
		var ok bool
		if mouse, ok = gs.Mice[m.MouseId]; !ok {
			return errors.New("invalid mouse id")
		}

		mousePath := fmt.Sprintf(".mice[\"%s\"]", m.MouseId)

		uneducated := gs.Population.Uneducated + 1

		var educatedCount int
		var educationCountPath string
		switch mouse.Education {
		case models.Education_CHEF:
			educatedCount = int(gs.Population.Chefs) - 1
			educationCountPath = ".population.chefs"
		case models.Education_SALESMOUSE:
			educatedCount = int(gs.Population.Salesmice) - 1
			educationCountPath = ".population.salesmice"
		case models.Education_GUARD:
			educatedCount = int(gs.Population.Guards) - 1
			educationCountPath = ".population.guards"
		case models.Education_THIEF:
			educatedCount = int(gs.Population.Thieves) - 1
			educationCountPath = ".population.thieves"
		case models.Education_PUBLICIST:
			educatedCount = int(gs.Population.Publicists) - 1
			educationCountPath = ".population.publicists"
		}

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			err := internal.RedisJsonSet(pipe, ctx, gsKey, fmt.Sprintf("%s.isEducated", mousePath), "false").Err()
			if err != nil {
				return fmt.Errorf("failed to set isEducated: %w", err)
			}

			err = internal.RedisJsonSet(pipe, ctx, gsKey, ".population.uneducated", uneducated).Err()
			if err != nil {
				return fmt.Errorf("failed to increase uneducated: %w", err)
			}

			err = internal.RedisJsonSet(pipe, ctx, gsKey, educationCountPath, educatedCount).Err()
			if err != nil {
				return fmt.Errorf("failed to decrease the educated value: %w", err)
			}

			return nil
		})
		return err
	}

	mutex := h.rdb.NewMutex("lock:" + gsKey)
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to obtain lock: %w", err)
	}
	err2 := txf()
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}
	if err2 != nil {
		return fmt.Errorf("failed to handle reschool: %w", err2)
	}

	h.fetchAndUpdateTimestamp(ctx, senderId)
	h.sendFullStateUpdate(ctx, senderId)

	return nil
}

func (h *handler) handleRenameMouse(ctx context.Context, senderId string, m *models.ClientMessage_RenameMouse) error {
	log.Info().
		Str("senderId", senderId).
		Str("MouseId", m.MouseId).
		Str("Name", m.Name).
		Msg("Received rename mouse message")

	gsKey := fmt.Sprintf("user:%s:gamestate", senderId)

	var gs models.GameState

	txf := func() error {
		// Get current game state
		s, err := internal.RedisJsonGet(h.rdb, ctx, gsKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
			return err
		}

		name := strings.TrimSpace(m.Name)

		if len(name) < 2 {
			return errors.New("name too short")
		}

		if len(name) > 30 {
			return errors.New("name too long")
		}

		var ok bool
		if _, ok = gs.Mice[m.MouseId]; !ok {
			return errors.New("invalid mouse id")
		}

		mousePath := fmt.Sprintf(".mice[\"%s\"]", m.MouseId)

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			path := fmt.Sprintf("%s.name", mousePath)
			value := fmt.Sprintf("\"%s\"", name)
			err := internal.RedisJsonSet(pipe, ctx, gsKey, path, value).Err()
			if err != nil {
				return fmt.Errorf("failed to set name: %w", err)
			}

			// Change name quest
			if q, ok := gs.Quests["4"]; ok && !q.Completed {
				path := ".quests[\"4\"].completed"
				value := "true"
				err := internal.RedisJsonSet(pipe, ctx, gsKey, path, value).Err()
				if err != nil {
					return fmt.Errorf("failed to set claimed reward to true: %w", err)
				}

			}

			return nil
		})
		return err
	}

	mutex := h.rdb.NewMutex("lock:" + gsKey)
	if err := mutex.Lock(); err != nil {
		return fmt.Errorf("failed to obtain lock: %w", err)
	}
	err2 := txf()
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return fmt.Errorf("failed to unlock: %w", err)
	}
	if err2 != nil {
		return fmt.Errorf("failed to handle rename: %w", err2)
	}

	h.fetchAndUpdateTimestamp(ctx, senderId)
	h.sendFullStateUpdate(ctx, senderId)

	return nil
}
