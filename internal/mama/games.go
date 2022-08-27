package mama

import (
	"context"
	"database/sql"
)

type Game struct {
	Id     string
	Title  string
	Host   string
	Status string
}

type gameCreator struct {
	db *sql.DB
}

func NewGameCreator(db *sql.DB) *gameCreator {
	return &gameCreator{db: db}
}

func (gc *gameCreator) CreateGame(ctx context.Context, title string, host string) error {
	_, err := gc.db.ExecContext(ctx,
		`INSERT INTO game (title, host, status) VALUES (?, ?, "starting")`,
		title, host)
	if err != nil {
		return err
	}

	return nil
}

func GetGame(ctx context.Context, db *sql.DB, gameId string) (*Game, error) {
	row := db.QueryRowContext(ctx, `SELECT id, title, host, status FROM game WHERE id = ?`, gameId)
	game := &Game{}
	err := row.Scan(&game.Id, &game.Title, &game.Host, &game.Status)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func JoinGame(ctx context.Context, db *sql.DB, userId string, gameId string) error {
	_, err := db.ExecContext(ctx, `INSERT INTO user_game (user_id, game_id) VALUES (?, ?)`, userId, gameId)
	if err != nil {
		return err
	}

	return nil
}

func GetAllGames(ctx context.Context, db *sql.DB) ([]*Game, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, title, host, status FROM game`)
	if err != nil {
		return nil, err
	}

	games := []*Game{}

	for rows.Next() {
		game := &Game{}
		err := rows.Scan(&game.Id, &game.Title, &game.Host, &game.Status)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func GetJoinedGames(ctx context.Context, db *sql.DB, userId string) ([]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT game_id FROM user_game WHERE user_id = ?`, userId)
	if err != nil {
		return nil, err
	}

	ids := []string{}

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
