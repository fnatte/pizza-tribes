package internal

import (
	"context"
	"time"

	. "github.com/fnatte/pizza-tribes/internal/models"
	"github.com/go-redis/redis/v8"
)

func GetNextUpdateTimestamp(gs *GameState) int64 {
	t := time.Now().Add(10 * time.Second).UnixNano()
	for i := range gs.ConstructionQueue {
		t = Min(t, gs.ConstructionQueue[i].CompleteAt)
	}
	for i := range gs.TrainingQueue {
		t = Min(t, gs.TrainingQueue[i].CompleteAt)
	}
	for i := range gs.TravelQueue {
		t = Min(t, gs.TravelQueue[i].ArrivalAt)
	}

	// Make the update time at least 100ms in the future to avoid
	// update loops in case of failures.
	t = Max(t, time.Now().Add(100*time.Millisecond).UnixNano())

	return t
}

func SetNextUpdate(r redis.Cmdable, ctx context.Context, userId string, gs *GameState) (int64, error) {
	return SetNextUpdateTo(r, ctx, userId, GetNextUpdateTimestamp(gs))
}

func SetNextUpdateTo(r redis.Cmdable, ctx context.Context, userId string, value int64) (int64, error) {
	return r.ZAdd(ctx, "user_updates", &redis.Z{
		Score:  float64(value),
		Member: userId,
	}).Result()
}
