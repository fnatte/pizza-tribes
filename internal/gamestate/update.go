package gamestate

import (
	"context"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
)

type GameStateUpdaterFn func(gs *models.GameState, tx *GameTx) error

type updater struct {
	gsRepo      persist.GameStateRepository
	reportsRepo persist.ReportsRepository
	userRepo    persist.UserRepository
	notifyRepo  persist.NotifyRepository
	worldRepo  persist.WorldRepository
}

type Updater interface {
	PerformUpdate(ctx context.Context, userId string, updater GameStateUpdaterFn) (*GameTx, error)
}

func NewUpdater(gsRepo persist.GameStateRepository,
	reportsRepo persist.ReportsRepository,
	userRepo persist.UserRepository,
	notifyRepo persist.NotifyRepository,
	worldRepo persist.WorldRepository) *updater {
	return &updater{
		gsRepo:      gsRepo,
		userRepo:    userRepo,
		reportsRepo: reportsRepo,
		notifyRepo:  notifyRepo,
		worldRepo:   worldRepo,
	}
}

func (u *updater) PerformUpdate(ctx context.Context, userId string, updater GameStateUpdaterFn) (*GameTx, error) {
	mutex := u.gsRepo.NewMutex(userId)
	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to obtain lock: %w", err)
	}

	var tx *GameTx

	gs, err := u.gsRepo.Get(ctx, userId)
	if err == nil {
		var t int64
		t, err = u.userRepo.GetUserLatestActivity(ctx, userId)
		if err == nil {
			tx = NewGameTx(userId, gs, time.Unix(0, t))

			err = updater(gs, tx)
			if err == nil {
				err = u.persistTx(ctx, tx)
			}
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

func (u *updater) persistTx(ctx context.Context, tx *GameTx) error {
	for userId, txu := range tx.Users {
		err := u.gsRepo.Patch(ctx, userId, txu.Gs, txu.PatchMask)
		if err != nil {
			return err
		}

		if txu.ReportsInvalidated {
			for _, report := range txu.Reports {
				if err = u.reportsRepo.Save(ctx, userId, report); err != nil {
					return fmt.Errorf("failed to save report: %w", err)
				}
			}
		}

		if len(txu.Messages) > 0 {
			for _, m := range txu.Messages {
				if _, err = u.notifyRepo.SendPushNotification(ctx, m); err != nil {
					return fmt.Errorf("failed to send push notification: %w", err)
				}
			}
		}
	}

	return nil
}
