package main

import (
	"context"
	"net/http"

	"github.com/fnatte/pizza-mouse/cmd/api/ws"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type wsHandler struct{}

func (h *wsHandler) HandleMessage(ctx context.Context, m []byte, c *ws.Client) {
	log.Info().Str("userId", c.UserId()).Msg("Received message")
}
func (h *wsHandler) HandleInit(ctx context.Context, c *ws.Client) error {
	log.Info().Str("userId", c.UserId()).Msg("Client connected")
	go (func() {
		c.Send([]byte("{ \"resources\": { \"coins\": 1, \"pizzas\": 1 } }"))
	})()
	return nil
}

func main() {
	log.Info().Msg("Starting Api")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	handler := wsHandler{}
	auth := NewAuthService(rdb)
	wsHub := ws.NewHub()
	wsEndpoint := ws.NewEndpoint(auth.Authorize, wsHub, &handler)

	r := mux.NewRouter()
	r.PathPrefix("/auth").Handler(http.StripPrefix("/auth", auth.Router()))
	r.Handle("/ws", wsEndpoint)

	go wsHub.Run()

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
