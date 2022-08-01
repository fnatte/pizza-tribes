package main

import (
	"context"
	"os"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/rs/zerolog/log"
)

func envOrDefault(key string, defaultVal string) string{
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func ensureWorld(ctx context.Context, r redis.RedisClient) error {
	world := internal.NewWorldService(r)
	if err := world.Initialize(ctx); err != nil {
		return err
	}

	return nil
}

func main() {
	log.Info().Msg("Starting migrator")

	ctx := context.Background()

	r := redis.NewRedisClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0,  // use default DB
	})

	err := ensureWorld(ctx, r)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure world")
	}

	log.Info().Msg("Migrator done")
}
