package main

import (
	"net/http"
	"strconv"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type WorldController struct {
	r     redis.RedisClient
	world *game.WorldService
	auth  *game.AuthService
}

func (c *WorldController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		err := c.auth.Authorize(req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize")
			w.WriteHeader(403)
			return
		}

		state, err := c.world.GetState(req.Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed to get state")
			w.WriteHeader(500)
			return
		}

		b, err := protojson.Marshal(state)
		if err != nil {
			w.WriteHeader(500)
			log.Error().Err(err).Msg("Failed to entries response")
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(b)
	})

	r.HandleFunc("/entries", func(w http.ResponseWriter, req *http.Request) {
		err := c.auth.Authorize(req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize")
			w.WriteHeader(403)
			return
		}

		paramX := req.URL.Query().Get("x")
		paramY := req.URL.Query().Get("y")
		paramR := req.URL.Query().Get("r")

		var x, y, radius int
		if x, err = strconv.Atoi(paramX); err != nil {
			w.WriteHeader(400)
			log.Error().Err(err).Msg("Param x and y is required")
			return
		}
		if y, err = strconv.Atoi(paramY); err != nil {
			w.WriteHeader(400)
			log.Error().Err(err).Msg("Param x and y is required")
			return
		}
		if paramR != "" {
			if radius, err = strconv.Atoi(paramR); err != nil {
				w.WriteHeader(400)
				log.Error().Err(err).Msg("Failed to parse radius")
				return
			}
		} else {
			radius = 10
		}

		entries, err := c.world.GetEntries(req.Context(), x, y, radius)
		if err != nil {
			w.WriteHeader(500)
			log.Error().Err(err).Msg("Failed to get zone")
			return
		}

		resp := models.EntriesResponse{
			Entries: entries,
		}

		b, err := protojson.Marshal(&resp)
		if err != nil {
			w.WriteHeader(500)
			log.Error().Err(err).Msg("Failed to entries response")
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(b)
	})

	return r
}
