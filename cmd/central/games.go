package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	"github.com/fnatte/pizza-tribes/internal/gamelet"
	"github.com/fnatte/pizza-tribes/internal/mama"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type GamesController struct {
	sqldb    *sql.DB
	auth     *game.AuthService
	userRepo persist.UserRepository
}

type gameResponse struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Joined bool   `json:"joined"`
}

func NewGamesController(sqldb *sql.DB, auth *game.AuthService, userRepo persist.UserRepository) *GamesController {
	return &GamesController{
		sqldb:    sqldb,
		auth:     auth,
		userRepo: userRepo,
	}
}

func (c *GamesController) ListGames(w http.ResponseWriter, req *http.Request) {
	err := c.auth.Authorize(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to authorize")
		w.WriteHeader(403)
		return
	}

	ctx := req.Context()

	userId, ok := ctx.Value("userId").(string)
	if !ok {
		log.Error().Err(err).Msg("Failed to get userId")
		w.WriteHeader(500)
	}

	games, err := mama.GetAllGames(ctx, c.sqldb)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get games")
		w.WriteHeader(500)
		return
	}

	joinedGames, err := mama.GetJoinedGames(ctx, c.sqldb, userId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get joined games")
		w.WriteHeader(500)
		return
	}

	resp := []gameResponse{}
	for _, game := range games {
		joined := false
		for _, id := range joinedGames {
			if id == game.Id {
				joined = true
				break
			}
		}

		resp = append(resp, gameResponse{
			Id:     game.Id,
			Title:  game.Title,
			Status: game.Status,
			Joined: joined,
		})
	}

	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("Failed to marshal games response")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func (c *GamesController) JoinGame(w http.ResponseWriter, req *http.Request) {
	err := c.auth.Authorize(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to authorize")
		w.WriteHeader(403)
		return
	}

	params := mux.Vars(req)
	gameId := params["gameId"]
	if gameId == "" {
		w.WriteHeader(400)
		return
	}

	ctx := req.Context()
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		log.Error().Err(err).Msg("Failed to get userId")
		w.WriteHeader(500)
	}

	game, err := mama.GetGame(ctx, c.sqldb, gameId)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	err = mama.JoinGame(ctx, c.sqldb, userId, gameId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to join game")
		w.WriteHeader(500)
		return
	}

	user, err := c.userRepo.GetUser(ctx, userId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")
		w.WriteHeader(500)
		return
	}

	gcClient := gamelet.NewGameletClient(game.Host)
	err = gcClient.JoinGame(userId, user.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post join game to gamelet")
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (c *GamesController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", c.ListGames)
	r.HandleFunc("/{gameId}/join", c.JoinGame).Methods("POST")

	return r
}
