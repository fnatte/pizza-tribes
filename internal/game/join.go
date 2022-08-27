package game

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
)

type GameCtrl struct {
	gsRepo persist.GameStateRepository
	gameUserRepo persist.GameUserRepository
	world *WorldService
	leaderboard *LeaderboardService
}

func NewGameCtrl(
	gsRepo persist.GameStateRepository,
	gameUserRepo persist.GameUserRepository,
	world *WorldService,
	leaderboard *LeaderboardService,
) *GameCtrl {
	return &GameCtrl{
		gsRepo: gsRepo,
		gameUserRepo: gameUserRepo,
		world: world,
		leaderboard: leaderboard,
	}
}

func (g *GameCtrl) JoinGame(ctx context.Context, userId string, username string) error {
	if err := g.gameUserRepo.CreateUser(ctx, userId, username); err != nil {
		return fmt.Errorf("failed to create game user: %w", err)
	}

	gs := models.NewGameState()
	err := g.gsRepo.Save(ctx, userId, gs)
	if err != nil {
		return fmt.Errorf("failed to save game state: %w", err)
	}

	x, y, err := g.world.AcquireTown(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to acquire town: %w", err)
	}
	gs.TownX = int32(x)
	gs.TownY = int32(y)
	g.gsRepo.Patch(ctx, userId, gs, &models.PatchMask{
		Paths: []string{"townX", "townY"},
	})
	if err != nil {
		return fmt.Errorf("failed to patch acquired town: %w", err)
	}

	coins := int64(gs.Resources.Coins)
	if err = g.leaderboard.UpdateUser(ctx, userId, coins); err != nil {
		return fmt.Errorf("failed to update leaderboard: %w", err)
	}

	return nil
}
