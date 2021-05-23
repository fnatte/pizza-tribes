package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
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

	txf := func(tx *redis.Tx) error {
		// Get current game state
		s, err := internal.RedisJsonGet(tx, ctx, gsKey, ".").Result()
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

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			_, err := internal.RedisJsonNumIncrBy(
				pipe, ctx, gsKey,
				".resources.coins",
				int64(-node.Cost)).Result()
			if err != nil {
				log.Error().Err(err).Msg("Failed to decrease coins")
				return err
			}

			err = internal.RedisJsonArrAppend(
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

	if err := h.rdb.Watch(ctx, txf, gsKey); err != nil {
		return err
	}

	internal.SetNextUpdate(h.rdb, ctx, senderId, &gs)

	h.sendFullStateUpdate(ctx, senderId)

	return nil
}
