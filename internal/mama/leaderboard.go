package mama

import (
	"context"
	"database/sql"
)

type LeaderboardRow struct {
	UserId string
	Username string
	Coins  int64
}

type LeaderboardResult struct {
	GameId string
	Skip int
	Limit int
	Rows []*LeaderboardRow
}

func GetLeaderboard(ctx context.Context, db *sql.DB, gameId string, skip int, limit int) (*LeaderboardResult, error) {
	query := `SELECT l.user_id, u.username, l.coins FROM leaderboard l
		LEFT JOIN user u ON u.id = l.user_id
		WHERE l.game_id = ?`
	rows, err := db.QueryContext(ctx, query, gameId)
	if err != nil {
		return nil, err
	}

	res := &LeaderboardResult{
		GameId: gameId,
		Skip: skip,
		Limit: limit,
	}

	for rows.Next() {
		row := &LeaderboardRow{}
		err := rows.Scan(&row.UserId, &row.Username, &row.Coins)
		if err != nil {
			return nil, err
		}

		res.Rows = append(res.Rows, row)
	}

	return res, nil
}

func GetLeaderboardRank(ctx context.Context, db *sql.DB, gameId string, userId string) (int, error) {
	query := `SELECT rank FROM (
		SELECT row_number() OVER (ORDER BY l.coins DESC) AS rank, l.user_id FROM leaderboard l
		WHERE l.game_id = ?
	)
	WHERE user_id = ?`

	var rank int
	err := db.QueryRowContext(ctx, query, gameId, userId).Scan(&rank)
	if err != nil {
		return 0, err
	}

	return rank, nil
}

