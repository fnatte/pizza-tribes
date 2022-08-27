package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/game"
	sqlrepo "github.com/fnatte/pizza-tribes/internal/game/persist/sql"
	_ "github.com/mattn/go-sqlite3"
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
	log.Info().Msg("Starting Central")

	setupLogging()

	port, err := strconv.Atoi(internal.EnvOrDefault("PORT", "8083"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse port")
		return
	}

	origins := strings.Split(internal.EnvOrDefault("ORIGIN", fmt.Sprintf("http://localhost:%d", port)), " ")

	sqliteDSN := internal.EnvOrDefault("SQLITE_DSN", "./pizzatribes.db")
	sqldb, err := sql.Open("sqlite3", sqliteDSN)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	userRepo := sqlrepo.NewUserRepo(sqldb)
	auth := game.NewAuthService()
	users := game.NewUserService(userRepo)

	gc := NewGamesController(sqldb, auth, userRepo)
	ac := NewAuthController(auth, users)

	r := mux.NewRouter()
	registerSubrouter(r, "/auth", ac.Handler())
	registerSubrouter(r, "/games", gc.Handler())

	h := handlers.CORS(
		handlers.AllowedOrigins(origins),
		handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{
			"Accept", "Accept-Language", "Authorization", "Content-Language", "Content-Type", "Origin",
		}),
	)(r)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), h)
}

func registerSubrouter(r *mux.Router, prefix string, handler http.Handler) {
	r.PathPrefix(prefix).Handler(http.StripPrefix(prefix, handler))
}
