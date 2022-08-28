package game

import (
	"context"
	"errors"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/game/persist"
)

type UserService struct {
	userRepo    persist.UserRepository
}

func NewUserService(userRepo persist.UserRepository) *UserService {
	return &UserService{
		userRepo:    userRepo,
	}
}

func (u *UserService) GetUserByUsername(ctx context.Context, username string) (*persist.User, error) {
	user, err := u.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, persist.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (u *UserService) CreateUser(ctx context.Context, username string, password string) (*persist.User, error) {
	userId, err := u.userRepo.CreateUser(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user, err := u.userRepo.GetUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
