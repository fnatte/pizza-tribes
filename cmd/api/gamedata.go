package main

import (
	"encoding/json"
	"net/http"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/rs/zerolog/log"
)

func GameDataHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(&internal.FullGameData)
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("Failed to marhsla full game data")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
