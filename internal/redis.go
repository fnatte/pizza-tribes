package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	redis.Cmdable
	JsonGet(ctx context.Context, key string, path string) *redis.StringCmd
	JsonSet(ctx context.Context, key string, path string, data interface{}) *redis.StatusCmd
}

type redisClient struct {
	*redis.Client
}

func NewRedisClient(rdb *redis.Client) RedisClient {
	return &redisClient{rdb}
}

func (c *redisClient) JsonGet(ctx context.Context, key string, path string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx, "JSON.GET", key, path)
	_ = c.Process(ctx, cmd)
	return cmd
}

func (c *redisClient) JsonSet(ctx context.Context, key string, path string, value interface{}) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx, "JSON.SET", key, path, value)
	_ = c.Process(ctx, cmd)
	return cmd
}
