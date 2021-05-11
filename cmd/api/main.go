package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/fnatte/pizza-tribes/cmd/api/ws"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting Api")
	ctx := context.Background()

	port, err := strconv.Atoi(envOrDefault("PORT", "8080"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse port")
		return
	}
	origin := envOrDefault("ORIGIN", "http://localhost:8080")

	// Setup redis client
	rc := internal.NewRedisClient(redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	}))

	// Initialize services and controllers
	auth := NewAuthService(rc)
	world := internal.NewWorldService(rc)
	leaderboard := internal.NewLeaderboardService(rc)
	wsHub := ws.NewHub()
	handler := wsHandler{rc: rc, world: world}
	wsEndpoint := ws.NewEndpoint(auth.Authorize, wsHub, &handler, origin)
	poller := poller{rdb: rc, hub: wsHub}
	ts := &TimeseriesService{r: rc, auth: auth}
	worldController := &WorldController{auth: auth, world: world}
	userController := &UserController{auth: auth, r: rc}
	leaderboardController := &LeaderboardController{
		auth:        auth,
		leaderboard: leaderboard}

	r := mux.NewRouter()
	r.Handle("/ws", wsEndpoint)
	r.HandleFunc("/gamedata", GameDataHandler)
	registerSubrouter(r, "/auth", auth.Handler())
	registerSubrouter(r, "/timeseries", ts.Handler())
	registerSubrouter(r, "/world", worldController.Handler())
	registerSubrouter(r, "/user", userController.Handler())
	registerSubrouter(r, "/leaderboard", leaderboardController.Handler())

	// Start web socket loop
	go wsHub.Run()
	go poller.run(ctx)

	// Start HTTP server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}

func envOrDefault(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}


func registerSubrouter(r *mux.Router, prefix string, handler http.Handler) {
	r.PathPrefix(prefix).Handler(http.StripPrefix(prefix, handler))
}
