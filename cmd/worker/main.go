package main

import (
	"context"
	"os"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func envOrDefault(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func main() {
	log.Info().Msg("Starting worker")

	debug := envOrDefault("DEBUG", "0") == "1"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	rc := internal.NewRedisClient(rdb)

	world := internal.NewWorldService(rc)

	h := &handler{rdb: rc, world: world}

	ctx := context.Background()

	for {
		res, err := rdb.BLPop(ctx, 30*time.Second, "wsin").Result()
		if err != nil {
			if err != redis.Nil {
				log.Error().Err(err).Msg("Error when dequeuing message")
			}
			continue
		}

		if len(res) < 2 {
			log.Error().Err(err).Msg("This should never happend. BLPop should always return a slice with two values")
			continue
		}

		msg := &internal.IncomingMessage{}
		msg.UnmarshalBinary([]byte(res[1]))

		m := &models.ClientMessage{}
		err = protojson.Unmarshal([]byte(msg.Body), m)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse incoming message")
			continue
		}

		h.Handle(ctx, msg.SenderId, m)
	}

}

