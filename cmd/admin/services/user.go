package services

import (
	"context"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
)

type userDeleter struct {
	gameUserRepo persist.GameUserRepository
	gsRepo persist.GameStateRepository
	world *game.WorldService
	userRepo persist.UserRepository
}

type UserDeleter interface {
	DeleteUser(ctx context.Context, userId string) error
}

func NewUserDeleter(gameUserRepo persist.GameUserRepository, gsRepo persist.GameStateRepository, world *game.WorldService, userRepo persist.UserRepository) *userDeleter {
	return &userDeleter{
		gameUserRepo: gameUserRepo,
		gsRepo: gsRepo,
		world: world,
		userRepo: userRepo,
	}
}

func (ud *userDeleter) DeleteUser(ctx context.Context, userId string) error {
	err := ud.userRepo.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}

	gs, err := ud.gsRepo.Get(ctx, userId)
	if err != nil {
		return err
	}

	err = ud.gameUserRepo.DeleteUser(ctx, userId)
	if err != nil {
		return err
	}

	if gs != nil {
		err = ud.world.RemoveEntry(ctx, int(gs.TownX), int(gs.TownY))
		if err != nil {
			return err
		}
	}

	return nil
}
