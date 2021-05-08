package main

import (
	"net/http"
	"strconv"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

type WorldController struct {
	r internal.RedisClient
	world *internal.WorldService
	auth *AuthService
}

func (c *WorldController) Router() http.Handler {
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

		var zone *internal.WorldZone

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
