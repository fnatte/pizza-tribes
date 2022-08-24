package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/rs/zerolog/log"
)

func findDiscoveredNode(gs *models.GameState, node *models.ResearchNode, d models.ResearchDiscovery) *models.ResearchNode {
	if node.Discovery == d {
		return node
	}

	// If the discovery is researched we can traverse its children
	if gs.HasDiscovery(node.Discovery) {
		for _, subnode := range node.Nodes {
			if n := findDiscoveredNode(gs, subnode, d); n != nil {
				return n
			}
		}
	}

	return nil
}

func (h *handler) handleStartResearch(ctx context.Context, senderId string, m *models.ClientMessage_StartResearch) error {
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

		// Traverse research tracks to find the node
		var node *models.ResearchNode
		for _, track := range internal.FullGameData.ResearchTracks {
			if node = findDiscoveredNode(&gs, track.RootNode, m.Discovery); node != nil {
				break
			}
		}

		if node == nil {
			return fmt.Errorf("All previous research has not been discovered")
		}
		if gs.HasDiscovery(node.Discovery) {
			return fmt.Errorf("This research has already been discovered")
		}

		if gs.Resources.Coins < node.Cost {
			return errors.New("Not enough coins")
		}

		// Calculate when this research will be completed.
		// If there's already already something being researched, this research will
		// be started at the end of previous one. If there's nothing in queue, it can
		// be started immediately (time.Now()).
		timeOffset := time.Now().UnixNano()
		if n := len(gs.ResearchQueue); n > 0 {
			timeOffset = gs.ResearchQueue[n-1].CompleteAt
		}
		completeAt := timeOffset + int64(node.ResearchTime)*1e9

		research := models.OngoingResearch{
			CompleteAt: completeAt,
			Discovery:  m.Discovery,
		}

		researchBytes, err := protojson.Marshal(&research)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal research")
			return err
		}

		_, err = h.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			_, err := redis.RedisJsonNumIncrBy(
				pipe, ctx, gsKey,
				".resources.coins",
				int64(-node.Cost)).Result()
			if err != nil {
				log.Error().Err(err).Msg("Failed to decrease coins")
				return err
			}

			err = redis.RedisJsonArrAppend(
				pipe,
				ctx,
				fmt.Sprintf("user:%s:gamestate", senderId),
				".researchQueue",
				researchBytes,
			).Err()
			if err != nil {
				return err
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
		return fmt.Errorf("failed to handle research: %w", err2)
	}

	internal.SetNextUpdate(h.rdb, ctx, senderId, &gs)

	h.sendFullStateUpdate(ctx, senderId)

	return nil
}


func (h *handler) handleBuyGeniusFlash(ctx context.Context, userId string, m *models.ClientMessage_BuyGeniusFlash) error {
	tx, err := h.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		nextId := len(gs.Discoveries) + int(gs.GeniusFlashes) + 1
		if nextId != int(m.Id) {
			return fmt.Errorf("not current genius id")
		}

		costs := internal.FullGameData.GeniusFlashCosts
		if nextId >= len(costs) {
			return fmt.Errorf("invalid genius id")
		}

		cost := costs[nextId]
		if gs.Resources.Coins < cost.Coins {
			return fmt.Errorf("not enough coins")
		}
		if gs.Resources.Pizzas < cost.Pizzas {
			return fmt.Errorf("not enough pizzas")
		}

		tx.Users[userId].IncrCoins(-cost.Coins)
		tx.Users[userId].IncrPizzas(-cost.Pizzas)
		tx.Users[userId].IncrGeniusFlashes(1)

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to handle genius flash: %w", err)
	}

	err = h.sendGameTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send game tx: %w", err)
	}

	return nil
}

