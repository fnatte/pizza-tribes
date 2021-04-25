package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	redis.UniversalClient
	JsonGet(ctx context.Context, key string, path string) *redis.StringCmd
	JsonSet(ctx context.Context, key string, path string, data interface{}) *redis.StatusCmd
	JsonNumIncrBy(ctx context.Context, key string, path string, value int64) *redis.StringCmd
}

type RedisProcesser interface {
	Process(ctx context.Context, cmd redis.Cmder) error
}

type redisClient struct {
	redis.UniversalClient
}

func NewRedisClient(rdb redis.UniversalClient) RedisClient {
	return &redisClient{rdb}
}

func RedisJsonGet(c RedisProcesser, ctx context.Context, key string, path string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx, "JSON.GET", key, path)
	_ = c.Process(ctx, cmd)
	return cmd
}

func RedisJsonNumIncrBy(c RedisProcesser, ctx context.Context, key string, path string, value int64) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx, "JSON.NUMINCRBY", key, path, value)
	_ = c.Process(ctx, cmd)
	return cmd
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

func (c *redisClient) JsonNumIncrBy(ctx context.Context, key string, path string, value int64) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx, "JSON.NUMINCRBY", key, path, value)
	_ = c.Process(ctx, cmd)
	return cmd
}
