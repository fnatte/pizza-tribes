package main

import (
	"encoding/json"
	"net/http"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/rs/zerolog/log"
)

type JoinController struct {
	gameCtrl *game.GameCtrl
}

type joinGameRequest struct {
	Username string `json:"username"`
	UserId string `json:"user_id"`
	Items []string `json:"items"`
}

type gameResponse struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
	Joined bool   `json:"joined"`
}

func NewJoinController(gameCtrl *game.GameCtrl) *JoinController {
	return &JoinController{
		gameCtrl: gameCtrl,
	}
}

func (c *JoinController) JoinGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := joinGameRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.UserId == "" || req.Username == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = c.gameCtrl.JoinGame(ctx, req.UserId, req.Username, req.Items)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get join game")
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

