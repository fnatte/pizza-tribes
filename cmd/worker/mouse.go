package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
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
		s, err := redis.RedisJsonGet(h.rdb, ctx, gsKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
			return err
		}

		if _, ok := gs.Mice[m.MouseId]; !ok {
			return errors.New("invalid mouse id")
		}

		mousePath := fmt.Sprintf(".mice[\"%s\"]", m.MouseId)
		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			err := redis.RedisJsonSet(pipe, ctx, gsKey, fmt.Sprintf("%s.isEducated", mousePath), "false").Err()
			if err != nil {
				return fmt.Errorf("failed to set isEducated: %w", err)
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

var validMouseName = regexp.MustCompile(`^[a-zA-Z]+(\s[a-zA-Z]+)?(\s[a-zA-Z]+)?$`)

func validateName(name string) error {
	if len(name) < 2 {
		return errors.New("name too short")
	}

	if len(name) > 30 {
		return errors.New("name too long")
	}

	if !validMouseName.MatchString(name) {
		return errors.New("invalid name")
	}

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
		s, err := redis.RedisJsonGet(h.rdb, ctx, gsKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
			return err
		}

		name := strings.TrimSpace(m.Name)
		if err = validateName(name); err != nil {
			return err
		}

		var ok bool
		if _, ok = gs.Mice[m.MouseId]; !ok {
			return errors.New("invalid mouse id")
		}

		mousePath := fmt.Sprintf(".mice[\"%s\"]", m.MouseId)

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			path := fmt.Sprintf("%s.name", mousePath)
			value := fmt.Sprintf("\"%s\"", name)
			err := redis.RedisJsonSet(pipe, ctx, gsKey, path, value).Err()
			if err != nil {
				return fmt.Errorf("failed to set name: %w", err)
			}

			// Change name quest
			if q, ok := gs.Quests["4"]; ok && !q.Completed {
				path := ".quests[\"4\"].completed"
				value := "true"
				err := redis.RedisJsonSet(pipe, ctx, gsKey, path, value).Err()
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

func (h *handler) handleSaveMouseAppearance(ctx context.Context, userId string, m *models.ClientMessage_SaveMouseAppearance) error {
	if !game.IsValidMouseAppearance(m.Appearance) {
		return errors.New("invalid appearance")
	}

	tx, err := h.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		if _, ok := gs.Mice[m.MouseId]; !ok {
			return errors.New("invalid mouse id")
		}

		tx.Users[userId].SetMouseAppearance(m.MouseId, m.Appearance)

		// Complete the set mouse appearance quest
		if q, ok := gs.Quests["14"]; ok && !q.Completed {
			tx.Users[userId].SetQuestCompleted("14")
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}

	err = h.sendGameTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send game tx: %w", err)
	}

	return nil
}

func (h *handler) handleSetAmbassadorMouse(ctx context.Context, userId string, m *models.ClientMessage_SetAmbassadorMouse) error {
	tx, err := h.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		if _, ok := gs.Mice[m.MouseId]; !ok {
			return errors.New("invalid mouse id")
		}

		tx.Users[userId].SetAmbassadorMouse(m.MouseId)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}

	err = h.sendGameTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send game tx: %w", err)
	}

	return nil
}
