package internal

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetNextUpdateTimestamp (gs *GameState) int64 {
	t := time.Now().Add(10 * time.Second).UnixNano()
	for i := range(gs.ConstructionQueue) {
		t = Min(t, gs.ConstructionQueue[i].CompleteAt)
	}
	for i := range(gs.TrainingQueue) {
		t = Min(t, gs.TrainingQueue[i].CompleteAt)
	}
	for i := range(gs.TravelQueue) {
		t = Min(t, gs.TravelQueue[i].ArrivalAt)
	}
	return t
}


func UpdateTimestamp(r RedisClient, ctx context.Context, userId string, gs *GameState) (int64, error) {
	return r.ZAdd(ctx, "user_updates", &redis.Z{
		Score: float64(GetNextUpdateTimestamp(gs)),
		Member: userId,
	}).Result()
}
