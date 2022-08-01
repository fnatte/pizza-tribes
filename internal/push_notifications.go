package internal

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fnatte/pizza-tribes/internal/redis"

	"firebase.google.com/go/messaging"
)

func SchedulePushNotification(ctx context.Context, r redis.RedisClient, msg *messaging.Message, sendAt time.Time) (int64, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return 0, err
	}

	return r.ZAdd(ctx, "push_notifications", &redis.Z{
		Score:  float64(sendAt.Unix()),
		Member: b,
	}).Result()
}
