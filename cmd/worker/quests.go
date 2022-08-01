package main

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/rs/zerolog/log"
)

func (h *handler) handleOpenQuest(ctx context.Context, senderId string, m *models.ClientMessage_OpenQuest) error {
	log.Info().
		Str("senderId", senderId).
		Str("QuestId", m.QuestId).
		Msg("Received open quest message")

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

		if _, ok := gs.Quests[m.QuestId]; !ok {
			return nil
		}

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			path := fmt.Sprintf(".quests[\"%s\"].opened", m.QuestId)
			value := "true"
			err := redis.RedisJsonSet(pipe, ctx, gsKey, path, value).Err()
			if err != nil {
				return fmt.Errorf("failed to set quest to opened: %w", err)
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
		return fmt.Errorf("failed to handle open quest: %w", err2)
	}

	h.fetchAndUpdateTimestamp(ctx, senderId)
	h.sendFullStateUpdate(ctx, senderId)

	return nil
}

func (h *handler) handleClaimQuestReward(ctx context.Context, senderId string, m *models.ClientMessage_ClaimQuestReward) error {
	log.Info().
		Str("senderId", senderId).
		Str("QuestId", m.QuestId).
		Msg("Received claim quest reward message")

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

		if q, ok := gs.Quests[m.QuestId]; !ok || q.ClaimedReward {
			return nil
		}

		var q *models.Quest
		for i := range internal.FullGameData.Quests {
			if internal.FullGameData.Quests[i].Id == m.QuestId {
				q = internal.FullGameData.Quests[i]
			}
		}

		if q == nil {
			log.Error().Msg("Could not find quest")
			return nil
		}

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			path := fmt.Sprintf(".quests[\"%s\"].claimedReward", m.QuestId)
			value := "true"
			err := redis.RedisJsonSet(pipe, ctx, gsKey, path, value).Err()
			if err != nil {
				return fmt.Errorf("failed to set claimed reward to true: %w", err)
			}

			if q.Reward.Coins > 0 {
				err := redis.RedisJsonNumIncrBy(
					pipe, ctx, gsKey,
					".resources.coins",
					int64(q.Reward.Coins)).Err()
				if err != nil {
					log.Error().Err(err).Msg("Failed to increase coins")
					return err
				}
			}

			if q.Reward.Pizzas > 0 {
				err := redis.RedisJsonNumIncrBy(
					pipe, ctx, gsKey,
					".resources.pizzas",
					int64(q.Reward.Pizzas)).Err()
				if err != nil {
					log.Error().Err(err).Msg("Failed to increase coins")
					return err
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
		return fmt.Errorf("failed to handle claim quest reward: %w", err2)
	}

	h.fetchAndUpdateTimestamp(ctx, senderId)
	h.sendFullStateUpdate(ctx, senderId)

	return nil
}

func (h *handler) handleCompleteVisitHelpPageQuest(ctx context.Context, senderId string, m *models.ClientMessage_CompleteVisitHelpPageQuest) error {
	log.Info().
		Str("senderId", senderId).
		Msg("Received complete visit help page quest message")

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

		if q, ok := gs.Quests["6"]; !ok || q.Completed {
			return nil
		}

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if q, ok := gs.Quests["6"]; ok && !q.Completed {
				path := ".quests[\"6\"].completed"
				value := "true"
				err := redis.RedisJsonSet(pipe, ctx, gsKey, path, value).Err()
				if err != nil {
					return fmt.Errorf("failed to set visit help page quest completed to true: %w", err)
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
		return fmt.Errorf("failed to handle complete visit help page quest: %w", err2)
	}

	h.fetchAndUpdateTimestamp(ctx, senderId)
	h.sendFullStateUpdate(ctx, senderId)

	return nil
}
