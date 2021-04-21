package main

import (
	"context"
	"net/http"
	"time"

	"github.com/fnatte/pizza-mouse/cmd/api/ws"
	"github.com/fnatte/pizza-mouse/internal"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type wsHandler struct {
	rdb *redis.Client
}

func (h *wsHandler) HandleMessage(ctx context.Context, m []byte, c *ws.Client) {
	log.Debug().Str("userId", c.UserId()).Msg("Received message")
	err := h.rdb.RPush(ctx, "wsin", &internal.IncomingMessage{
		SenderId: c.UserId(),
		Body:     string(m),
	}).Err()
	if err != nil {
		log.Error().Err(err).Msg("Error when pushing incoming message to redis")
	}
}
func (h *wsHandler) HandleInit(ctx context.Context, c *ws.Client) error {
	log.Info().Str("userId", c.UserId()).Msg("Client connected")
	// go (func() {
	// 	c.Send([]byte("{ \"resources\": { \"coins\": 1, \"pizzas\": 1 } }"))
	// })()
	return nil
}

type poller struct {
	rdb *redis.Client
	hub *ws.Hub
}

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

		msg := &internal.OutgoingMessage{}
		msg.UnmarshalBinary([]byte(res[1]))

		log.Info().Str("receiverId", msg.ReceiverId).Msg("OutgoingMessage")

		p.hub.SendTo(msg.ReceiverId, []byte(msg.Body))
	}

}

func main() {
	/*
		msg := internal.ClientMessage{
			Id: "test-123",
			Type: &internal.ClientMessage_Tap_{
				Tap: &internal.ClientMessage_Tap{
					Amount: 32,
				},
			},
		}

		val, err := protojson.Marshal(&msg)
		log.Info().Msg(string(val))
	*/

	log.Info().Msg("Starting Api")

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	handler := wsHandler{rdb: rdb}
	auth := NewAuthService(rdb)
	wsHub := ws.NewHub()
	wsEndpoint := ws.NewEndpoint(auth.Authorize, wsHub, &handler)
	poller := poller{rdb: rdb, hub: wsHub}

	r := mux.NewRouter()
	r.PathPrefix("/auth").Handler(http.StripPrefix("/auth", auth.Router()))
	r.Handle("/ws", wsEndpoint)

	go wsHub.Run()
	go poller.run(ctx)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
