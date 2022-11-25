package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

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

type leaderboardRowResponse struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Coins    int64  `json:"coins"`
}

type leaderboardResponse struct {
	GameId string                    `json:"gameId"`
	Skip   int                       `json:"skip"`
	Limit  int                       `json:"limit"`
	Rows   []*leaderboardRowResponse `json:"rows"`
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

func (c *GamesController) ListPreviousGames(w http.ResponseWriter, req *http.Request) {
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
		return
	}

	games, err := mama.GetArchivedGamesWithUser(ctx, c.sqldb, userId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get games")
		w.WriteHeader(500)
		return
	}

	resp := []gameResponse{}
	for _, game := range games {
		resp = append(resp, gameResponse{
			Id:     game.Id,
			Title:  game.Title,
			Status: game.Status,
			Joined: true,
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
		return
	}

	games, err := mama.GetActiveGames(ctx, c.sqldb)
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

	items, err := c.userRepo.GetUserItems(ctx, userId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user items")
		w.WriteHeader(500)
		return
	}

	gcClient := gamelet.NewGameletClient(game.Host)
	err = gcClient.JoinGame(userId, user.Username, items)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post join game to gamelet")
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (c *GamesController) ShowGame(w http.ResponseWriter, req *http.Request) {
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
		return
	}

	game, err := mama.GetGame(ctx, c.sqldb, gameId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get game")
		w.WriteHeader(500)
		return
	}

	joined, err := mama.HasJoinedGame(ctx, c.sqldb, userId, gameId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get joined status")
		w.WriteHeader(500)
		return
	}

	resp := &gameResponse{
		Id:     game.Id,
		Title:  game.Title,
		Status: game.Status,
		Joined: joined,
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

func (c *GamesController) ShowGameLeaderboard(w http.ResponseWriter, req *http.Request) {
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

	skip := 0
	paramSkip := req.URL.Query().Get("skip")
	if paramSkip != "" {
		if skip, err = strconv.Atoi(paramSkip); err != nil {
			w.WriteHeader(400)
			log.Warn().Err(err).Msg("Could not parse skip")
			return
		}
	}

	limit := 20
	paramLimit := req.URL.Query().Get("limit")
	if paramLimit != "" {
		if limit, err = strconv.Atoi(paramLimit); err != nil {
			w.WriteHeader(400)
			log.Warn().Err(err).Msg("Could not parse limit")
			return
		}
		if limit < 1 {
			limit = 1
		} else if limit > 100 {
			limit = 100
		}
	}

	ctx := req.Context()

	l, err := mama.GetLeaderboard(ctx, c.sqldb, gameId, skip, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get leaderboard")
		w.WriteHeader(500)
		return
	}

	resp := &leaderboardResponse{
		GameId: l.GameId,
		Skip:   l.Skip,
		Limit:  l.Limit,
		Rows:   make([]*leaderboardRowResponse, 0, len(l.Rows)),
	}
	for _, row := range l.Rows {
		resp.Rows = append(resp.Rows, &leaderboardRowResponse{
			UserId:   row.UserId,
			Coins:    row.Coins,
			Username: row.Username,
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

func (c *GamesController) ShowGameLeaderboardMeRank(w http.ResponseWriter, req *http.Request) {
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
		return
	}

	rank, err := mama.GetLeaderboardRank(ctx, c.sqldb, gameId, userId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get rank")
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	io.WriteString(w, strconv.Itoa(rank))
}

func (c *GamesController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", c.ListGames)
	r.HandleFunc("/previous", c.ListPreviousGames)
	r.HandleFunc("/{gameId}", c.ShowGame).Methods("GET")
	r.HandleFunc("/{gameId}/join", c.JoinGame).Methods("POST")
	r.HandleFunc("/{gameId}/leaderboard", c.ShowGameLeaderboard).Methods("GET")
	r.HandleFunc("/{gameId}/leaderboard/me/rank", c.ShowGameLeaderboardMeRank).Methods("GET")

	return r
}
