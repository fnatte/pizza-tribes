package main

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting Api")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	auth := NewAuthService(rdb)

	r := mux.NewRouter()
	r.PathPrefix("/auth").Handler(http.StripPrefix("/auth", auth.Router()))
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		err := auth.Authorize(r)
		if err != nil {
			log.Error().Err(err).Msg("Auth error")
			http.Error(w, "Auth error", http.StatusUnauthorized)
			return
		}

		userId, ok := r.Context().Value("userId").(string)
		if !ok {
			log.Error().Msg("Failed to get userId")
			http.Error(w, "Ops. An error occured.", http.StatusInternalServerError)
			return
		}

		log.Info().Str("userId", userId).Msg("Hello from test")
		w.Write([]byte(userId))
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
