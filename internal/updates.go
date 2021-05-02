package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func UpdateTimestamp(r RedisClient, ctx context.Context, userId string, gs *GameState) (int64, error) {
	//return RedisZAddLt(r, ctx, "user_updates", &redis.Z{
	return r.ZAdd(ctx, "user_updates", &redis.Z{
		Score: float64(GetNextUpdateTimestamp(gs)),
		Member: userId,
	}).Result()
}
