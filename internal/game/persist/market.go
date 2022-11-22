package persist

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
)

type marketRepo struct {
	rc redis.RedisClient
}

func NewMarketRepository(rc redis.RedisClient) *marketRepo {
	return &marketRepo{
		rc: rc,
	}
}

func (r *marketRepo) GetGlobalDemandScore(ctx context.Context) (float64, error) {
	sum, err := r.rc.Get(ctx, "demands:sum").Float64()
	if err != nil && err != redis.Nil {
		return 0, fmt.Errorf("failed to get global demand: %w", err)
	}

	return sum, nil
}

func (r *marketRepo) SetUserDemandScore(ctx context.Context, userId string, demand float64) error {
	err := r.rc.ZAdd(ctx, "demands", &redis.Z{
		Score:  demand,
		Member: userId,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to set demand: %w", err)
	}

	_, err = r.updateGlobalDemandScore(ctx)
	if err != nil {
		return fmt.Errorf("failed to update global demand: %w", err)
	}

	return nil
}

func (r *marketRepo) setGlobalDemandScore(ctx context.Context, sum float64) error {
	return r.rc.Set(ctx, "demands:sum", sum, 0).Err()
}

func (r *marketRepo) updateGlobalDemandScore(ctx context.Context) (float64, error) {
	packed, err := r.rc.ZRangeWithScores(ctx, "demands", 0, -1).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get demands: %w", err)
	}
	if len(packed) == 0 {
		return 0, nil
	}

	sum := 0.0
	for _, item := range packed {
		sum += item.Score
	}

	err = r.setGlobalDemandScore(ctx, sum)
	if err != nil {
		err = fmt.Errorf("failed to set global demand sum: %w", err)
		return 0, err
	}

	return sum, nil
}

func (r *marketRepo) GetDemandRankByUserId(ctx context.Context, userId string) (int64, error) {
	rank, err := r.rc.ZRevRank(ctx, "demands", userId).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get user demand rank: %w", err)
	}

	return rank + 1, nil
}

func (r *marketRepo) GetDemandLeaderboard(ctx context.Context, skip int) (*models.DemandLeaderboard, error) {
	res, err := r.rc.ZRevRangeWithScores(ctx, "demands", int64(skip), int64(skip)+20).Result()
	if err != nil {
		return nil, err
	}

	sum, err := r.GetGlobalDemandScore(ctx)
	if err != nil {
		return nil, err
	}

	board := &models.DemandLeaderboard{
		Skip: int32(skip),
		Rows: make([]*models.DemandLeaderboard_Row, len(res)),
	}

	for i, row := range res {
		userId := row.Member.(string)
		userKey := fmt.Sprintf("user:%s", userId)
		username, err := r.rc.HGet(ctx, userKey, "username").Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get username: %w", err)
		}

		board.Rows[i] = &models.DemandLeaderboard_Row{
			UserId:   userId,
			Username: username,
			Demand:   row.Score,
			MarketShare: row.Score / sum,
		}
	}

	return board, nil
}
