package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fnatte/pizza-tribes/cmd/api/ws"
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting Api")

	debug := envOrDefault("DEBUG", "0") == "1"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	ctx := context.Background()

	port, err := strconv.Atoi(envOrDefault("PORT", "8080"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse port")
		return
	}
	origins := strings.Split(envOrDefault("ORIGIN", "http://localhost:8080"), " ")

	// Setup redis client
	rc := redis.NewRedisClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	// Initialize services and controllers
	auth := game.NewAuthService()
	world := game.NewWorldService(rc)
	leaderboard := game.NewLeaderboardService(rc)
	gameUserRepo := persist.NewGameUserRepository(rc)
	gsRepo := persist.NewGameStateRepository(rc)
	marketRepo := persist.NewMarketRepository(rc)
	wsHub := ws.NewHub()
	handler := wsHandler{rc: rc, world: world, gsRepo: gsRepo, marketRepo: marketRepo, userRepo: gameUserRepo}
	wsEndpoint := ws.NewEndpoint(auth.Authorize, wsHub, &handler, origins)
	poller := poller{rdb: rc, hub: wsHub}
	ts := &TimeseriesService{r: rc, auth: auth}

	worldController := &WorldController{auth: auth, world: world}
	userController := &UserController{auth: auth, r: rc, gsRepo: gsRepo}
	leaderboardController := &LeaderboardController{
		auth:        auth,
		leaderboard: leaderboard}
	pushNotificationsController := &PushNotificationsController{
		auth: auth,
		r:    rc}
	demandLeaderboardController := &DemandLeaderboardController{
		auth:        auth,
		marketRepo: marketRepo}

	speed, err := world.GetSpeed(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get game speed")
		os.Exit(1)
		return
	}

	game.AlterGameDataForSpeed(speed)

	r := mux.NewRouter()
	r.Handle("/ws", wsEndpoint)
	r.HandleFunc("/gamedata", GameDataHandler)
	registerSubrouter(r, "/timeseries", ts.Handler())
	registerSubrouter(r, "/world", worldController.Handler())
	registerSubrouter(r, "/user", userController.Handler())
	registerSubrouter(r, "/leaderboard", leaderboardController.Handler())
	registerSubrouter(r, "/demand_leaderboard", demandLeaderboardController.Handler())
	registerSubrouter(r, "/push_notifications", pushNotificationsController.Handler())

	// Start web socket loop
	go wsHub.Run()
	go poller.run(ctx)

	// Start HTTP server
	h := handlers.CORS(
		handlers.AllowedOrigins(origins),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Authorization", "Content-Language", "Content-Type", "Origin"}),
	)(r)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), h)
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
