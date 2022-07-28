package persist

import (
	"context"

	"github.com/fnatte/pizza-tribes/internal/models"
)

type GameStateRepository interface {
	NewMutex(userId string) Mutex
	Get(ctx context.Context, userId string) (*models.GameState, error)
	Patch(ctx context.Context, userId string, gs *models.GameState, patch *models.ServerMessage_PatchMask) error
	Save(ctx context.Context, userId string, gs *models.GameState) error
}

type ReportsRepository interface {
	Save(ctx context.Context, userId string, report *models.Report) error
	Get(ctx context.Context, userId string) ([]*models.Report, error)
	MarkRead(ctx context.Context, userId string, reportId string) error
}

type UserRepository interface {
	SetUserLatestActivity(ctx context.Context, userId string, value int64) error
	GetUserLatestActivity(ctx context.Context, userId string) (int64, error)
	GetAllUsers(ctx context.Context) ([]string, error)
}

type Mutex interface {
	Lock() error
	LockContext(context.Context) error
	Unlock() (bool, error)
}
