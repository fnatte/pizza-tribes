package main

import (
	"context"
	"os"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func envOrDefault(key string, defaultVal string) string{
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func ensureWorld(ctx context.Context, r internal.RedisClient) error {
	world := internal.NewWorldService(r)
	if err := world.Initilize(ctx); err != nil {
		return err
	}

	return nil
}

func main() {
	log.Info().Msg("Starting migrator")

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0,  // use default DB
	})
	r := internal.NewRedisClient(rdb)

	err := ensureWorld(ctx, r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure world")
	}

	log.Info().Msg("Migrator done")
}
