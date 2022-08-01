package main

import (
	"fmt"
	"net/http"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type UserController struct {
	r      redis.RedisClient
	auth   *AuthService
	gsRepo persist.GameStateRepository
}

func (c *UserController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/{userId}", func(w http.ResponseWriter, r *http.Request) {
		err := c.auth.Authorize(r)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize")
			w.WriteHeader(403)
			return
		}

		params := mux.Vars(r)
		userId := params["userId"]
		if userId == "" {
			w.WriteHeader(400)
			return
		}

		username, err := c.r.HGet(r.Context(), fmt.Sprintf("user:%s", userId), "username").Result()
		if err != nil {
			w.WriteHeader(404)
			return
		}

		gs, err := c.gsRepo.Get(r.Context(), userId)
		if err != nil {
			log.Error().Msg("failed to get game state")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		var ambassador *models.Mouse
		if gs.AmbassadorMouseId != "" {
			var ok bool
			ambassador, ok = gs.Mice[gs.AmbassadorMouseId]
			if !ok {
				log.Error().Msg("failed to look up ambassador mouse")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		res := &models.ApiUserResponse{
			Username: username,
		}
		if ambassador != nil {
			res.Ambassador = &models.ApiUserResponse_Ambassador{
				Appearance: ambassador.Appearance,
			}
		}

		b, err := protojson.Marshal(res)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal username struct")
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(b)
	})

	return r
}
