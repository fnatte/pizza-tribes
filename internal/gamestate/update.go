package gamestate

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
)

type GameStateUpdater func(gs *models.GameState, tx *GameTx) error

func PerformUpdate(ctx context.Context, gsRepo persist.GameStateRepository, reportRepo persist.ReportsRepository, userId string, updater GameStateUpdater) (*GameTx, error) {
	mutex := gsRepo.NewMutex(userId)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to obtain lock: %w", err)
	}

	var tx *GameTx

	gs, err := gsRepo.Get(ctx, userId)
	if err == nil {
		tx = NewGameTx(userId, gs)
		err = updater(gs, tx)
		if err == nil {
			err = persistTx(ctx, userId, tx, gsRepo, reportRepo)
		}
	}

	if ok, err2 := mutex.Unlock(); !ok || err2 != nil {
		return nil, fmt.Errorf("failed to unlock: %w", err2)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to perform update: %w", err)
	}

	return tx, nil
}

func persistTx(ctx context.Context, userId string, tx *GameTx, gsRepo persist.GameStateRepository, reportRepo persist.ReportsRepository) error {
	for _, txu := range tx.Users {
		err := gsRepo.Patch(ctx, userId, txu.Gs, txu.PatchMask)
		if err != nil {
			return err
		}

		if txu.ReportsInvalidated {
			for _, report := range txu.Reports {
				if err = reportRepo.Save(ctx, userId, report); err != nil {
					return fmt.Errorf("failed to save report: %w", err)
				}
			}

		}
	}

	return nil
}
