package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogging() {
	debug := internal.EnvOrDefault("DEBUG", "0") == "1"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func main() {
	log.Info().Msg("Starting Gamelet")

	setupLogging()

	port, err := strconv.Atoi(internal.EnvOrDefault("PORT", "8082"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse port")
		return
	}

	origins := strings.Split(internal.EnvOrDefault("ORIGIN", fmt.Sprintf("http://localhost:%d", port)), " ")

	rc := redis.NewRedisClient(&redis.Options{
		Addr:     internal.EnvOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: internal.EnvOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	world := game.NewWorldService(rc)
	leaderboard := game.NewLeaderboardService(rc)
	gameUserRepo := persist.NewGameUserRepository(rc)
	gsRepo := persist.NewGameStateRepository(rc)
	gameCtrl := game.NewGameCtrl(gsRepo, gameUserRepo, world, leaderboard, rc)

	gc := NewJoinController(gameCtrl)
	lc := NewLeaderboardController(leaderboard)

	r := mux.NewRouter()
	r.HandleFunc("/join", gc.JoinGame).Methods("POST")
	r.HandleFunc("/leaderboard", lc.GetLeaderboard).Methods("GET")

	h := handlers.CORS(
		handlers.AllowedOrigins(origins),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{
			"Accept", "Accept-Language", "Authorization", "Content-Language", "Content-Type", "Origin",
		}),
	)(r)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), h)
}

