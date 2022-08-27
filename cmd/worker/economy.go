package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
)

func (h *handler) handleSetPizzaPrice(ctx context.Context, userId string, m *models.ClientMessage_SetPizzaPrice) error {

	if !game.IsValidPizzaPrice(m.PizzaPrice) {
		return errors.New("invalid pizza price")
	}

	tx, err := h.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		if !game.HasBuildingMinLevel(gs, models.Building_TOWN_CENTRE, 2) {
			return errors.New("cannot change pizza price without town centre level 2")
		}

		tx.Users[userId].SetPizzaPrice(m.PizzaPrice)

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
