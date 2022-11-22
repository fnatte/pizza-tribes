package main

import (
	"net/http"
	"strconv"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type DemandLeaderboardController struct {
	marketRepo persist.MarketRepository
	auth        *game.AuthService
}

func (c *DemandLeaderboardController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/me/rank", func(w http.ResponseWriter, r *http.Request) {
		err := c.auth.Authorize(r)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize")
			w.WriteHeader(403)
			return
		}

		userId, ok := game.GetUserIdFromContext(r.Context())
		if !ok {
			log.Error().Err(err).Msg("Failed to get user id")
			return
		}

		rank, err := c.marketRepo.GetDemandRankByUserId(r.Context(), userId)
		if err != nil {
			w.WriteHeader(500)
			log.Error().Err(err).Msg("Failed to get user demand rank")
			return
		}

		rankStr := strconv.FormatInt(rank, 10)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(rankStr))
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := c.auth.Authorize(r)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize")
			w.WriteHeader(403)
			return
		}

		skip := 0
		paramSkip := r.URL.Query().Get("skip")
		if paramSkip != "" {
			if skip, err = strconv.Atoi(paramSkip); err != nil {
				w.WriteHeader(400)
				log.Error().Err(err).Msg("Could not parse skip")
				return
			}
		}

		board, err := c.marketRepo.GetDemandLeaderboard(r.Context(), skip)
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
	})

	return r
}
