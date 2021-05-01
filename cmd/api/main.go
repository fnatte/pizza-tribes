package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fnatte/pizza-tribes/cmd/api/ws"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

type wsHandler struct {
	rc internal.RedisClient
}

func (h *wsHandler) HandleMessage(ctx context.Context, m []byte, c *ws.Client) {
	log.Debug().Str("userId", c.UserId()).Msg("Received message")
	err := h.rc.RPush(ctx, "wsin", &internal.IncomingMessage{
		SenderId: c.UserId(),
		Body:     string(m),
	}).Err()
	if err != nil {
		log.Error().Err(err).Msg("Error when pushing incoming message to redis")
	}
}
func (h *wsHandler) HandleInit(ctx context.Context, c *ws.Client) error {
	gs := internal.GameState{
		Resources: &internal.GameState_Resources{},
		Lots: map[string]*internal.GameState_Lot{},
	}

	log.Info().Str("userId", c.UserId()).Msg("Client connected")
	gsKey := fmt.Sprintf("user:%s:gamestate", c.UserId())
	s, err := h.rc.JsonGet(ctx, gsKey, ".").Result()
	if err != nil {
		if err != redis.Nil {
			return err
		}
		b, err := protojson.Marshal(&gs)
		if err != nil {
			return err
		}
		err = h.rc.JsonSet(ctx, gsKey, ".", string(b)).Err()
		if err != nil {
			return err
		}
		log.Info().Msg("Initilized new game state for user")
	} else {
		gs.LoadProtoJson([]byte(s))
		if err != nil {
			return err
		}
		log.Info().Interface("gameState", &gs).Msg("Got game state")
	}

	go (func() {
		msg := gs.ToStateChangeMessage()
		b, err := protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send init game state patch")
			return
		}

		c.Send(b)
		log.Info().Msg("Sent init game state")
	})()

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
	log.Info().Msg("Starting Api")

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rc := internal.NewRedisClient(rdb)

	handler := wsHandler{rc: rc}
	auth := NewAuthService(rdb)
	wsHub := ws.NewHub()
	wsEndpoint := ws.NewEndpoint(auth.Authorize, wsHub, &handler)
	poller := poller{rdb: rdb, hub: wsHub}

	r := mux.NewRouter()
	r.PathPrefix("/auth").Handler(http.StripPrefix("/auth", auth.Router()))
	r.Handle("/ws", wsEndpoint)
	r.HandleFunc("/gamedata", func (w http.ResponseWriter, r *http.Request) {
		b, err := protojson.MarshalOptions{
			UseEnumNumbers: true,
		}.Marshal(&internal.FullGameData)
		if err != nil {
			w.WriteHeader(500)
			log.Error().Err(err).Msg("Failed to marhsla full game data")
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	})

	go wsHub.Run()
	go poller.run(ctx)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
