package internal

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type TimeseriesDataPoint struct {
	Timestamp int64
	Value float64
}

type RedisClient interface {
	redis.UniversalClient
	JsonGet(ctx context.Context, key string, path string) *redis.StringCmd
	JsonSet(ctx context.Context, key string, path string, data interface{}) *redis.StatusCmd
	JsonNumIncrBy(ctx context.Context, key string, path string, value int64) *redis.StringCmd
	ZAddLt(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
	TsCreate(ctx context.Context, key string, retention int64) *redis.StatusCmd
	TsInfo(ctx context.Context, key string) *redis.SliceCmd
	TsAdd(ctx context.Context, key string, timestamp, value int64) *redis.StatusCmd
	TsRange(ctx context.Context, key string, from, to int64) ([]*TimeseriesDataPoint, error)
	TsRangeAggr(ctx context.Context, key string, from, to int64, aggrType string, timeBucket int64) ([]*TimeseriesDataPoint, error)
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

func (c *redisClient) TsCreate(ctx context.Context, key string, retention int64) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx, "TS.CREATE", key, "RETENTION", retention)
	_ = c.Process(ctx, cmd)
	return cmd
}

func (c *redisClient) TsInfo(ctx context.Context, key string) *redis.SliceCmd {
	cmd := redis.NewSliceCmd(ctx, "TS.INFO", key)
	_ = c.Process(ctx, cmd)
	return cmd
}

func (c *redisClient) TsAdd(ctx context.Context, key string, timestamp, value int64) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx, "TS.ADD", key, timestamp, value)
	_ = c.Process(ctx, cmd)
	return cmd
}

func (c *redisClient) TsRange(ctx context.Context, key string, from, to int64) ([]*TimeseriesDataPoint, error) {
	cmd := redis.NewSliceCmd(ctx, "TS.RANGE", key, from, to)
	_ = c.Process(ctx, cmd)

	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return parseDataPoints(res)
}

func (c *redisClient) TsRangeAggr(ctx context.Context, key string, from, to int64, aggrType string, timeBucket int64) ([]*TimeseriesDataPoint, error) {
	cmd := redis.NewSliceCmd(ctx, "TS.RANGE", key, from, to, "AGGREGATION", aggrType, timeBucket)
	_ = c.Process(ctx, cmd)

	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return parseDataPoints(res)
}

func parseDataPoints(res []interface{}) ([]*TimeseriesDataPoint, error) {
	arr := make([]*TimeseriesDataPoint, len(res))
	for i := range(res) {
		raw, ok := res[i].([]interface{})
		if !ok {
			return nil, errors.New("Failed to parse timeseries range result")
		}

		dp, err := parseDataPoint(raw)
		if err != nil {
			return nil, err
		}

		arr[i] = dp
	}

	return arr, nil
}

func parseDataPoint(raw []interface{}) (dp *TimeseriesDataPoint, err error) {
	time, ok := raw[0].(int64)
	if !ok {
		return nil, errors.New("Could not parse data point timestamp")
	}

	val, ok := raw[1].(string)
	if !ok {
		return nil, errors.New("Could not parse data point value")
	}

	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return
	}

	dp = &TimeseriesDataPoint{
		Timestamp: time,
		Value: f,
	}
	return
}

