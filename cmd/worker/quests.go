package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
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

func (h *handler) handleClaimQuestReward(ctx context.Context, userId string, m *models.ClientMessage_ClaimQuestReward) error {
	log.Info().
		Str("userId", userId).
		Str("questId", m.QuestId).
		Str("selectedItem", m.SelectedOneOfItem).
		Msg("Received claim quest reward message")

	var q *models.Quest
	for i := range game.FullGameData.Quests {
		if game.FullGameData.Quests[i].Id == m.QuestId {
			q = game.FullGameData.Quests[i]
			break
		}
	}
	if q == nil {
		return errors.New("could not find quest")
	}

	itemIndex := 0

	if len(q.Reward.OneOfItems) > 0 {
		if m.SelectedOneOfItem == "" {
			return errors.New("must select item")
		}

		validItem := false
		for i, item := range q.Reward.OneOfItems {
			if item == m.SelectedOneOfItem {
				validItem = true
				itemIndex = i
				break
			}
		}
		if !validItem {
			return errors.New("selected item was not one of the available ones")
		}
	}

	tx, err := h.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		if qs, ok := gs.Quests[m.QuestId]; !ok || qs.ClaimedReward {
			return nil
		}

		txu := tx.Users[userId]
		txu.SetQuestClaimedReward(m.QuestId)
		if q.Reward.Coins > 0 {
			txu.IncrCoins(q.Reward.Coins)
		}
		if q.Reward.Pizzas > 0 {
			txu.IncrPizzas(q.Reward.Pizzas)
		}
		if len(q.Reward.OneOfItems) > 0 {
			item := q.Reward.OneOfItems[itemIndex]
			txu.AppendAppearancePart(item)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to perform update on claim quest reward: %w", err)
	}

	err = h.sendGameTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send game tx: %w", err)
	}

	return nil
}

func (h *handler) handleCompleteQuest(ctx context.Context, senderId string, m *models.ClientMessage_CompleteQuest) error {
	log.Info().
		Str("senderId", senderId).
		Str("questId", m.QuestId).
		Msg("Received complete quest message")

	tx, err := h.updater.PerformUpdate(ctx, senderId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		switch m.QuestId {
		case "6", "9", "11":
			if q, ok := gs.Quests[m.QuestId]; ok && !q.Completed {
				tx.Users[senderId].SetQuestCompleted(m.QuestId)
			}
		default:
			return errors.New("invalid quest id")
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed handle complete quest: %w", err)
	}

	err = h.sendGameTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed handle complete quest: %w", err)
	}

	h.fetchAndUpdateTimestamp(ctx, senderId)

	return nil
}
