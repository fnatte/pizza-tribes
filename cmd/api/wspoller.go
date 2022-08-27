package main

import (
	"context"
	"time"

	"github.com/fnatte/pizza-tribes/cmd/api/ws"
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type poller struct {
	rdb redis.UniversalClient
	hub *ws.Hub
}

// Pumps websocket messages from redis to the websocket hub
func (p *poller) run(ctx context.Context) {
	for {
		res, err := p.rdb.BLPop(ctx, 30*time.Second, "wsout").Result()
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

		msg := &game.OutgoingMessage{}
		msg.UnmarshalBinary([]byte(res[1]))

		p.hub.SendTo(msg.ReceiverId, []byte(msg.Body))
	}
}

