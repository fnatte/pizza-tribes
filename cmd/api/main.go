package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fnatte/pizza-tribes/cmd/api/ws"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

type wsHandler struct {
	rc internal.RedisClient
	world *internal.WorldService
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
		Population: &internal.GameState_Population{},
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
		b, err := protojson.MarshalOptions{
			EmitUnpopulated: true,
		}.Marshal(&gs)
		log.Info().Msg(string(b))
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
	}

	// Make sure the user has town in world
	if gs.TownX == 0 && gs.TownY == 0 {
		x, y, err := h.world.AcquireTown(ctx, c.UserId())
		if err != nil {
			return fmt.Errorf("failed to acquire town: %w", err)
		}
		gs.TownX = int32(x)
		gs.TownY = int32(y)
		log.Info().Msg("Town acquired")
	}

	// Make sure the user is enqueued for updates
	_, err = internal.UpdateTimestamp(h.rc, ctx, c.UserId(), &gs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure user updates")
		return err
	}

	// Get username
	userKey := fmt.Sprintf("user:%s", c.UserId())
	username, err := h.rc.HGet(ctx, userKey, "username").Result()
	if err != nil {
		return fmt.Errorf("failed to get username: %w", err)
	}

	go (func() {
		msg := &internal.ServerMessage{
			Id: xid.New().String(),
			Payload: &internal.ServerMessage_User_{
				User: &internal.ServerMessage_User{
					Username: username,
				},
			},
		}
		b, err := protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send init game state patch")
			return
		}

		c.Send(b)
		log.Info().Msg("Sent user message")
	})()

	go (func() {
		msg := gs.ToStateChangeMessage()
		b, err := protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send init game state patch")
			return
		}

		c.Send(b)

		msg = internal.CalculateStats(&gs).ToServerMessage()
		b, err = protojson.Marshal(msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send init stats")
			return
		}
		c.Send(b)

		log.Info().Msg("Sent init game state and stats")
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

		p.hub.SendTo(msg.ReceiverId, []byte(msg.Body))
	}

}

func envOrDefault(key string, defaultVal string) string{
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func main() {
	log.Info().Msg("Starting Api")

	port, err := strconv.Atoi(envOrDefault("PORT", "8080"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse port")
		return
	}

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0,  // use default DB
	})

	rc := internal.NewRedisClient(rdb)

	origin := envOrDefault("ORIGIN", "http://localhost:8080")

	auth := NewAuthService(rdb)
	world := internal.NewWorldService(rc)
	leaderboard := internal.NewLeaderboardService(rc)
	wsHub := ws.NewHub()
	handler := wsHandler{rc: rc, world: world}
	wsEndpoint := ws.NewEndpoint(auth.Authorize, wsHub, &handler, origin)
	poller := poller{rdb: rdb, hub: wsHub}
	ts := &TimeseriesService{ r: rc, auth: auth }
	worldController := &WorldController{ auth: auth, world: world }
	userController := &UserController{ auth: auth, r: rc }
	leaderboardController := &LeaderboardController{ auth: auth, leaderboard: leaderboard }

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
	r.PathPrefix("/timeseries").Handler(http.StripPrefix("/timeseries", ts.Router()))
	r.PathPrefix("/world").Handler(http.StripPrefix("/world", worldController.Router()))
	r.PathPrefix("/user").Handler(http.StripPrefix("/user", userController.Router()))
	r.PathPrefix("/leaderboard").Handler(http.StripPrefix("/leaderboard", leaderboardController.Router()))

	go wsHub.Run()
	go poller.run(ctx)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
