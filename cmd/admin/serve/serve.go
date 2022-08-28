package serve

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fnatte/pizza-tribes/cmd/admin/db"
	"github.com/fnatte/pizza-tribes/cmd/admin/services"
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	persistsql "github.com/fnatte/pizza-tribes/internal/game/persist/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func envOrDefault(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func Serve() {
	log.Info().Msg("Starting admin server")
	origins := strings.Split(envOrDefault("ORIGIN", "http://localhost:8080"), " ")

	port, err := strconv.Atoi(envOrDefault("ADMIN_PORT", "8081"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse port")
		return
	}

	rc := db.NewRedisClient()
	sqldb := db.NewSqlClient()

	gsRepo := persist.NewGameStateRepository(rc)
	reportsRepo := persist.NewReportsRepository(rc)
	notifyRepo := persist.NewNotifyRepository(rc)
	worldRepo := persist.NewWorldRepository(rc)
	gameUserRepo := persist.NewGameUserRepository(rc)
	// marketRepo := persist.NewMarketRepository(rc)
	world := game.NewWorldService(rc)
	userRepo := persistsql.NewUserRepo(sqldb)
	updater := gamestate.NewUpdater(gsRepo, reportsRepo, gameUserRepo, notifyRepo, worldRepo)

	userDeleter := services.NewUserDeleter(gameUserRepo, gsRepo, world, userRepo)

	userController := NewUserController(rc, sqldb, updater, userDeleter)
	testController := NewTestController(rc, sqldb, userDeleter)

	r := mux.NewRouter()
	r.PathPrefix("/users").Handler(userController.Handler())
	r.PathPrefix("/test").Handler(testController.Handler())

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
