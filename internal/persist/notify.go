package persist

import (
	"context"
	"encoding/json"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
)

type notifyRepo struct {
	rc internal.RedisClient
}

func NewNotifyRepository(rc internal.RedisClient) *notifyRepo {
	return &notifyRepo{
		rc: rc,
	}
}

func (r *notifyRepo) SchedulePushNotification(ctx context.Context, msg *messaging.Message, sendAt time.Time) (int64, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return 0, err
	}

	return r.rc.ZAdd(ctx, "push_notifications", &redis.Z{
		Score:  float64(sendAt.Unix()),
		Member: b,
	}).Result()
}

func (r *notifyRepo) SendPushNotification(ctx context.Context, msg *messaging.Message) (int64, error) {
	return r.SchedulePushNotification(ctx, msg, time.Now())
}
