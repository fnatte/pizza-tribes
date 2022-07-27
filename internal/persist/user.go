package persist

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
)

type userRepo struct {
	rdb internal.RedisClient
}

func NewUserRepository(rdb internal.RedisClient) *userRepo {
	return &userRepo{
		rdb: rdb,
	}
}

func (r *userRepo) GetUserLatestActivity(ctx context.Context, userId string) (int64, error) {
	key := fmt.Sprintf("user:%s:latest_activity", userId)

	val, err := r.rdb.Get(ctx, key).Int64()
	if err != nil && err != redis.Nil {
		return 0, fmt.Errorf("failed to get latest user activity: %w", err)
	}

	return val, nil
}

func (r *userRepo) SetUserLatestActivity(ctx context.Context, userId string, val int64) error {
	key := fmt.Sprintf("user:%s:latest_activity", userId)

	err := r.rdb.Set(ctx, key, val, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set latest user activity: %w", err)
	}

	return nil
}

func (r *userRepo) GetAllUsers(ctx context.Context) ([]string, error) {
	return r.rdb.SMembers(ctx, "users").Result()
}
