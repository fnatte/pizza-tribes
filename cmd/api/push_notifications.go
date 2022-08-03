package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type PushNotificationsController struct {
	r    redis.RedisClient
	auth *internal.AuthService
}

type registerPushNotificationRequest struct {
	Token string `json:"token"`
}

func (c *PushNotificationsController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		err := c.auth.Authorize(r)
		if err != nil {
			log.Error().Err(err).Msg("Failed to authorize")
			w.WriteHeader(403)
			return
		}

		req := registerPushNotificationRequest{}
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse login request")
			http.Error(w, "Invalid login request", http.StatusBadRequest)
			return
		}

		if req.Token == "" {
			w.WriteHeader(400)
			return
		}

		userId, ok := r.Context().Value("userId").(string)
		if !ok {
			log.Warn().Msg("Failed to get account id")
			w.WriteHeader(500)
			return
		}

		if err := c.r.SAdd(r.Context(), fmt.Sprintf("user:%s:fcm_tokens", userId), req.Token).Err(); err != nil {
			log.Error().Str("userId", userId).Str("token", req.Token).Err(err).Msg("failed to add fcm token")
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
	})

	return r
}
