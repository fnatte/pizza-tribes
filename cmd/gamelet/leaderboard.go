package main

import (
	"net/http"
	"strconv"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/rs/zerolog/log"
)

type LeaderboardController struct {
	leaderboard *game.LeaderboardService
}

func NewLeaderboardController(leaderboard *game.LeaderboardService) *LeaderboardController {
	return &LeaderboardController{leaderboard: leaderboard}
}

func (c *LeaderboardController) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	var err error
	skip := 0
	paramSkip := r.URL.Query().Get("skip")
	if skip, err = strconv.Atoi(paramSkip); err != nil {
		w.WriteHeader(400)
		log.Error().Err(err).Msg("Could not parse skip")
		return
	}

	board, err := c.leaderboard.Get(r.Context(), skip)
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("Failed to get leaderboard")
		return
	}

	b, err := protojson.Marshal(board)
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("Failed to marshal leaderboard")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}
