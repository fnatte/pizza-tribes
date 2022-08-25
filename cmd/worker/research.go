package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
)

func (h *handler) handleStartResearch(ctx context.Context, userId string, m *models.ClientMessage_StartResearch) error {
	tx, err := h.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		if gs.GeniusFlashes <= 0 {
			return fmt.Errorf("not enough genius flashes")
		}

		var r *models.ResearchInfo
		var ok bool
		if r, ok = internal.FullGameData.Research[int32(m.Discovery)]; !ok {
			return fmt.Errorf("discovery not found")
		}

		if gs.HasDiscovery(m.Discovery) {
			return fmt.Errorf("research already been discovered")
		}

		for _, x := range gs.ResearchQueue {
			if x.Discovery == m.Discovery {
				return fmt.Errorf("research already being researched")
			}
		}

		for _, d := range r.Requirements {
			if !gs.HasDiscovery(d) {
				return fmt.Errorf("research not available, all requirements were not fulfilled")
			}
		}

		// Calculate when this research will be completed.
		// If there's already already something being researched, this research will
		// be started at the end of previous one. If there's nothing in queue, it can
		// be started immediately (time.Now()).
		timeOffset := time.Now().UnixNano()
		if n := len(gs.ResearchQueue); n > 0 {
			timeOffset = gs.ResearchQueue[n-1].CompleteAt
		}
		completeAt := timeOffset + int64(r.ResearchTime)*1e9

		tx.Users[userId].AppendResearchQueue(&models.OngoingResearch{
			CompleteAt: completeAt,
			Discovery:  m.Discovery,
		})
		tx.Users[userId].IncrGeniusFlashes(-1)
		

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to handle start research: %w", err)
	}

	err = h.sendGameTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send game tx: %w", err)
	}

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

