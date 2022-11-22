package persist

import (
	"context"
	"errors"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal/game/models"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUsernameTaken = errors.New("username is taken")
var ErrInvalidUsername = errors.New("username is invalid")


type User struct {
	Id             string
	Username       string
	HashedPassword string
}

type GameUser struct {
	Uid      string
	Username string
}

type GameStateRepository interface {
	NewMutex(userId string) Mutex
	Get(ctx context.Context, userId string) (*models.GameState, error)
	Patch(ctx context.Context, userId string, gs *models.GameState, patch *models.PatchMask) error
	Save(ctx context.Context, userId string, gs *models.GameState) error
}

type ReportsRepository interface {
	Save(ctx context.Context, userId string, report *models.Report) error
	Get(ctx context.Context, userId string) ([]*models.Report, error)
	MarkRead(ctx context.Context, userId string, reportId string) error
}

type GameUserRepository interface {
	SetUserLatestActivity(ctx context.Context, userId string, value int64) error
	GetUserLatestActivity(ctx context.Context, userId string) (int64, error)
	GetAllUsers(ctx context.Context) ([]string, error)
	GetUserCount(ctx context.Context) (int64, error)
	CreateUser(ctx context.Context, userId, username string) (error)
	GetUser(ctx context.Context, userId string) (*GameUser, error)
	DeleteUser(ctx context.Context, userId string) (error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, username string, password string) (string, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	GetUser(ctx context.Context, userId string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserItems(ctx context.Context, userId string) ([]string, error)
	DeleteUser(ctx context.Context, userId string) (error)
}

type NotifyRepository interface {
	SendPushNotification(ctx context.Context, msg *messaging.Message) (int64, error)
	SchedulePushNotification(ctx context.Context, msg *messaging.Message, sendAt time.Time) (int64, error)
}

type MarketRepository interface {
	GetGlobalDemandScore(ctx context.Context) (float64, error)
	SetUserDemandScore(ctx context.Context, userId string, demand float64) error
	GetDemandRankByUserId(ctx context.Context, userId string) (int64, error)
	GetDemandLeaderboard(ctx context.Context, skip int) (*models.DemandLeaderboard, error)
}

type WorldRepository interface {
	GetState(ctx context.Context) (*models.WorldState, error)
}

type Mutex interface {
	Lock() error
	LockContext(context.Context) error
	Unlock() (bool, error)
}
