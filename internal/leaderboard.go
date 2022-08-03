package internal

import (
	"context"
	"fmt"

	. "github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/redis"
)

type LeaderboardService struct {
	r redis.RedisClient
}

func NewLeaderboardService(r redis.RedisClient) *LeaderboardService {
	return &LeaderboardService{r: r}
}

func (s *LeaderboardService) GetRankByUserId(ctx context.Context, userId string) (int64, error) {
	rank, err := s.r.ZRank(ctx, "leaderboard", userId).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get user leaderboard rank: %w", err)
	}

	return rank, nil
}

func (s *LeaderboardService) Get(ctx context.Context, skip int) (*Leaderboard, error) {
	res, err := s.r.ZRevRangeWithScores(ctx, "leaderboard", int64(skip), 20).Result()
	if err != nil {
		return nil, err
	}

	board := &Leaderboard{
		Skip: int32(skip),
		Rows: make([]*Leaderboard_Row, len(res)),
	}

	for i, row := range res {
		userId := row.Member.(string)
		userKey := fmt.Sprintf("user:%s", userId)
		username, err := s.r.HGet(ctx, userKey, "username").Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get username: %w", err)
		}

		board.Rows[i] = &Leaderboard_Row{
			UserId:   userId,
			Coins:    int64(row.Score),
			Username: username,
		}
	}

	return board, nil
}

func (s *LeaderboardService) UpdateUser(ctx context.Context, userId string, coins int64) error {
	s.r.ZAdd(ctx, "leaderboard", &redis.Z{
		Score:  float64(coins),
		Member: userId,
	})

	return nil
}
