package main

import (
	"context"
	"time"

	"github.com/fnatte/pizza-mouse/internal"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting worker")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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

		log.Info().Str("senderId", msg.SenderId).Msg("Received message")

		rdb.RPush(ctx, "wsout", &internal.OutgoingMessage{
			ReceiverId: msg.SenderId,
			Body: "{ \"resources\": { \"coins\": 2, \"pizzas\": 2 } }",
		})
	}

}
