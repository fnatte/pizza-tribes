package serve

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/gorilla/mux"
)

type setupTestResponse struct {
	AccessToken string `json:"accessToken"`
	User        *user  `json:"user"`
}

type setupTestRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type testController struct {
	rc redis.RedisClient
}

func NewTestController(r redis.RedisClient) *testController {
	return &testController{
		rc: r,
	}
}

func (c *testController) HandleSetupTest(w http.ResponseWriter, r *http.Request) {
	req := setupTestRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	userRepo := persist.NewUserRepository(c.rc)
	gsRepo := persist.NewGameStateRepository(c.rc)
	world := internal.NewWorldService(c.rc)

	// Delete previous user
	userId, err := userRepo.FindUser(ctx, req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userId != "" {
		gs, err := gsRepo.Get(ctx, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = userRepo.DeleteUser(ctx, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if gs != nil {
			err = world.RemoveEntry(ctx, int(gs.TownX), int(gs.TownY))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Create new user
	id, err := userRepo.CreateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gs := models.NewGameState()
	gsRepo.Save(ctx, id, gs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := userRepo.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a := internal.AuthService{}
	token, err := a.CreateToken(id, time.Now().Add(2*time.Minute))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(&setupTestResponse{
		AccessToken: token,
		User: &user{
			Id:       u.Id,
			Username: u.Username,
		}})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func (c *testController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/test/setup", c.HandleSetupTest).Methods("POST")

	return r
}
