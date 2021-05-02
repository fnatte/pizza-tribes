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
	ZAddLt(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
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

func RedisJsonSet(c RedisProcesser, ctx context.Context, key string, path string, value interface{}) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx, "JSON.SET", key, path, value)
	_ = c.Process(ctx, cmd)
	return cmd
}

func RedisJsonArrAppend(c RedisProcesser, ctx context.Context, key string, path string, value interface{}) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx, "JSON.ARRAPPEND", key, path, value)
	_ = c.Process(ctx, cmd)
	return cmd
}

func RedisJsonArrTrim(c RedisProcesser, ctx context.Context, key string, path string, start int, end int) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx, "JSON.ARRTRIM", key, path, start, end)
	_ = c.Process(ctx, cmd)
	return cmd
}

func RedisZAddLt(c RedisProcesser, ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	const n = 3
	a := make([]interface{}, n+2*len(members))
	a[0], a[1], a[2] = "zadd", key, "lt"

	for i, m := range members {
		a[n+2*i] = m.Score
		a[n+2*i+1] = m.Member
	}
	cmd := redis.NewIntCmd(ctx, a...)
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

func (c *redisClient) ZAddLt(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return RedisZAddLt(c, ctx, key, members...)
}
