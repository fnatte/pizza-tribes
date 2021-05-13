package main

import (
	"net/http"
	"strconv"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type WorldController struct {
	r     internal.RedisClient
	world *internal.WorldService
	auth  *AuthService
}

func (c *WorldController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/zone", func(w http.ResponseWriter, r *http.Request) {
		err := c.auth.Authorize(r)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize")
			w.WriteHeader(403)
			return
		}

		paramX := r.URL.Query().Get("x")
		paramY := r.URL.Query().Get("y")
		paramIdx := r.URL.Query().Get("idx")

		var zone *models.WorldZone

		if paramIdx != "" {
			var idx int
			if idx, err = strconv.Atoi(paramIdx); err != nil {
				w.WriteHeader(400)
				log.Error().Err(err).Msg("Could not parse idx")
				return
			}

			zone, err = c.world.GetZoneIdx(r.Context(), idx)
			if err != nil {
				w.WriteHeader(500)
				log.Error().Err(err).Msg("Failed to get zone")
				return
			}
		} else {
			var x, y int
			if x, err = strconv.Atoi(paramX); err != nil {
				w.WriteHeader(400)
				log.Error().Err(err).Msg("Param x and y, or idx is required")
				return
			}
			if y, err = strconv.Atoi(paramY); err != nil {
				w.WriteHeader(400)
				log.Error().Err(err).Msg("Param x and y, or idx is required")
				return
			}

			zone, err = c.world.GetZoneXY(r.Context(), x, y)
			if err != nil {
				w.WriteHeader(500)
				log.Error().Err(err).Msg("Failed to get zone")
				return
			}
		}

		b, err := protojson.Marshal(zone)
		if err != nil {
			w.WriteHeader(500)
			log.Error().Err(err).Msg("Failed to marshal zone")
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(b)
	})

	return r
}
