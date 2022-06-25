package gamestate

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
)

type GameStateUpdater func(gs *models.GameState, tx *GameTx) error

func PerformUpdate(ctx context.Context, repo persist.GameStateRepository, userId string, updater GameStateUpdater) (*GameTx, error) {
	mutex := repo.NewMutex(userId)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to obtain lock: %w", err)
	}

	var err2 error
	var gs *models.GameState
	var tx *GameTx

	gs, err2 = repo.Get(ctx, userId)

	if err2 == nil {
		tx = NewGameTx(userId, gs)
		err2 = updater(gs, tx)
		if err2 == nil {
			for _, txu := range tx.Users {
				err2 = repo.Patch(ctx, userId, txu.GsPatch)
				if err2 != nil {
					break
				}
			}
		}
	}

	if ok, err := mutex.Unlock(); !ok || err != nil {
		return nil, fmt.Errorf("failed to unlock: %w", err)
	}

	if err2 != nil {
		return nil, fmt.Errorf("failed to perform update: %w", err2)
	}

	return tx, nil
}
