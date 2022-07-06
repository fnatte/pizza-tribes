package internal

import (
	"context"
	"strconv"
	"time"

	. "github.com/fnatte/pizza-tribes/internal/models/gamestate"
	"github.com/go-redis/redis/v8"
)

func strToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func GetNextUpdateTimestamp(gs *GameState) int64 {
	t := time.Now().Add(10 * time.Second).UnixNano()
	for i := range gs.ConstructionQueue {
		t = Min(t, strToInt64(gs.ConstructionQueue[i].CompleteAt))
	}
	for i := range gs.TrainingQueue {
		t = Min(t, strToInt64(gs.TrainingQueue[i].CompleteAt))
	}
	for i := range gs.TravelQueue {
		t = Min(t, strToInt64(gs.TravelQueue[i].ArrivalAt))
	}

	// Make the update time at least 100ms in the future to avoid
	// update loops in case of failures.
	t = Max(t, time.Now().Add(100 * time.Millisecond).UnixNano())

	return t
}

func SetNextUpdate(r redis.Cmdable, ctx context.Context, userId string, gs *GameState) (int64, error) {
	return r.ZAdd(ctx, "user_updates", &redis.Z{
		Score:  float64(GetNextUpdateTimestamp(gs)),
		Member: userId,
	}).Result()
}
