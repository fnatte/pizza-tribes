package main

import (
	"net/http"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func GameDataHandler(w http.ResponseWriter, r *http.Request) {
	b, err := protojson.MarshalOptions{
		UseEnumNumbers: true,
	}.Marshal(&game.FullGameData)
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("Failed to marhsla full game data")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
