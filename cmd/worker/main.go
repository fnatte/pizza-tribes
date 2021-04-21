package main

import (
	"context"
	"time"

	"github.com/fnatte/pizza-mouse/internal"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	log.Info().Msg("Starting worker")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	h := &handler{ rdb: rdb }

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

		m := &internal.ClientMessage{}
		err = protojson.Unmarshal([]byte(msg.Body), m)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse incoming message")
			continue
		}

		h.Handle(ctx, msg.SenderId, m)
	}

}

