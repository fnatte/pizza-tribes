package internal

import (
	"context"
	"encoding/json"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/go-redis/redis/v8"
)

func SchedulePushNotification(ctx context.Context, r RedisClient, msg *messaging.Message, sendAt time.Time) (int64, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return 0, err
	}

	return r.ZAdd(ctx, "push_notifications", &redis.Z{
		Score:  float64(sendAt.Unix()),
		Member: b,
	}).Result()
}
