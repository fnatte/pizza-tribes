package internal

import (
	"context"
	"fmt"
	"time"
)

func ensureTimeserie(ctx context.Context, r RedisClient, key string, retention int64) error {
	err := r.TsInfo(ctx, key).Err()
	if err == nil {
		return nil
	}

	err = r.TsCreate(ctx, key, retention).Err()
	if err != nil {
		return err
	}

	return nil
}

func EnsureTimeseries(ctx context.Context, r RedisClient, userId string) error {
	retention := (7 * 24 * time.Hour).Milliseconds()
	err := ensureTimeserie(ctx, r, fmt.Sprintf("user:%s:ts_coins", userId), retention)
	if err != nil {
		return err
	}

	err = ensureTimeserie(ctx, r, fmt.Sprintf("user:%s:ts_pizzas", userId), retention)
	if err != nil {
		return err
	}

	return nil
}

func AddMetricCoins(ctx context.Context, r RedisClient, userId string, time, value int64) error {
	k := fmt.Sprintf("user:%s:ts_coins", userId)
	return r.TsAdd(ctx, k, time, value).Err()
}

func AddMetricPizzas(ctx context.Context, r RedisClient, userId string, time, value int64) error {
	k := fmt.Sprintf("user:%s:ts_pizzas", userId)
	return r.TsAdd(ctx, k, time, value).Err()
}

func FetchPizzasTimeseries(ctx context.Context, r RedisClient, userId string) ([]*TimeseriesDataPoint, error) {
	from := time.Now().Unix() * 1000 - (24 * time.Hour).Milliseconds()
	to := time.Now().Unix() * 1000
	k := fmt.Sprintf("user:%s:ts_pizzas", userId)
	timeBucket := (60 * time.Minute).Milliseconds()

	return r.TsRangeAggr(ctx, k, from, to, "avg", timeBucket)
}

func FetchCoinsTimeseries(ctx context.Context, r RedisClient, userId string) ([]*TimeseriesDataPoint, error) {
	from := time.Now().Unix() * 1000 - (24 * time.Hour).Milliseconds()
	to := time.Now().Unix() * 1000
	k := fmt.Sprintf("user:%s:ts_coins", userId)
	timeBucket := (60 * time.Minute).Milliseconds()

	return r.TsRangeAggr(ctx, k, from, to, "avg", timeBucket)
}
