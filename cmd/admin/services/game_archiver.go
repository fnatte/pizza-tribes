package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fnatte/pizza-tribes/internal/gamelet"
)

type GameArchiver struct {
	db  *sql.DB
	glc *gamelet.GameletClient
}

type row struct {
	gameId string
	userId string
	coins  int64
}

func NewGameArchiver(db *sql.DB, glc *gamelet.GameletClient) *GameArchiver {
	return &GameArchiver{db: db, glc: glc}
}

func (ga *GameArchiver) getRows(gameId string) ([]row, error) {
	rows := []row{}

	for skip := 0; true; skip++ {
		lb, err := ga.glc.GetLeaderboard(0)
		if err != nil {
			return nil, fmt.Errorf("failed to get leaderboard: %w", err)
		}

		for _, r := range lb.Rows {
			rows = append(rows, row{
				gameId: gameId,
				userId: r.UserId,
				coins:  r.Coins,
			})
		}

		// If we got less than 20 rows we are at the last frame
		if len(lb.Rows) < 20 {
			break
		}
	}

	return rows, nil
}

func (ga *GameArchiver) ArchiveGame(ctx context.Context, gameId string) error {
	rows, err := ga.getRows(gameId)
    if err != nil {
		return fmt.Errorf("archive game: %w", err)
    }

	args := make([]interface{}, 0, len(rows) * 3)
	for _, row := range rows {
		args = append(args, row.gameId, row.userId, row.coins)
	}
	vals := strings.Repeat("(?, ?, ?), ", len(rows) - 1) + "(?, ?, ?)"

	tx, err := ga.db.BeginTx(ctx, nil)
    if err != nil {
		return fmt.Errorf("archive game: %w", err)
    }
    defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE game SET status = ? WHERE id = ?", "archived", gameId)
	if err != nil {
		return fmt.Errorf("archive game: %w", err)
	}

	_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO leaderboard (game_id, user_id, coins) VALUES %s", vals), args...)
	if err != nil {
		return fmt.Errorf("archive game: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("archive game: %w", err)
	}

	return nil
}

