package internal

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
)

type UserService struct {
	userRepo    persist.UserRepository
	gsRepo      persist.GameStateRepository
	world       *WorldService
	leaderboard *LeaderboardService
}

func NewUserService(userRepo persist.UserRepository, gsRepo persist.GameStateRepository, world *WorldService, leaderboard *LeaderboardService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		gsRepo:      gsRepo,
		world:       world,
		leaderboard: leaderboard,
	}
}

func (u *UserService) CreateUser(ctx context.Context, username string, password string) (*persist.UserDbo, error) {
	userId, err := u.userRepo.CreateUser(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user, err := u.userRepo.GetUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	gs := models.NewGameState()
	u.gsRepo.Save(ctx, userId, gs)
	if err != nil {
		return nil, fmt.Errorf("failed to save game state: %w", err)
	}

	x, y, err := u.world.AcquireTown(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire town: %w", err)
	}
	gs.TownX = int32(x)
	gs.TownY = int32(y)
	u.gsRepo.Patch(ctx, userId, gs, &models.PatchMask{
		Paths: []string{"townX", "townY"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to patch acquired town: %w", err)
	}

	coins := int64(gs.Resources.Coins)
	if err = u.leaderboard.UpdateUser(ctx, userId, coins); err != nil {
		return nil, fmt.Errorf("failed to update leaderboard: %w", err)
	}

	return user, nil
}
