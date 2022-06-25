package persist

import (
	"context"

	"github.com/fnatte/pizza-tribes/internal/models"
)

type GameStateRepository interface {
	NewMutex(userId string) Mutex
	Get(ctx context.Context, userId string) (*models.GameState, error)
	Patch(ctx context.Context, userId string, patch *Patch) error
}

type Mutex interface {
	Lock() error
	LockContext(context.Context) error
	Unlock() (bool, error)
}
