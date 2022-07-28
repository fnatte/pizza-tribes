package db

import (
	"os"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
)

func envOrDefault(key string, defaultVal string) string{
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}


func NewRedisClient() internal.RedisClient {
	return internal.NewRedisClient(redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	}))
}

