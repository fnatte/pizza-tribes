package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type UserController struct {
	r internal.RedisClient
	auth *AuthService
}

func (c *UserController) Router() http.Handler {
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

		b, err := json.Marshal(struct{ Username string `json:"username"` }{
			Username: username,
		})
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

