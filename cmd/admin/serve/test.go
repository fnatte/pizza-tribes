package serve

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fnatte/pizza-tribes/cmd/admin/services"
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	sqlrepo "github.com/fnatte/pizza-tribes/internal/game/persist/sql"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	"github.com/fnatte/pizza-tribes/internal/gamelet"
	"github.com/fnatte/pizza-tribes/internal/mama"
	"github.com/gorilla/mux"
)

type setupTestResponse struct {
	AccessToken string `json:"accessToken"`
	User        *user  `json:"user"`
	GameId      string `json:"gameId"`
}

type setupTestRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type testController struct {
	rc          redis.RedisClient
	sqldb       *sql.DB
	userDeleter services.UserDeleter
}

func NewTestController(r redis.RedisClient, sqldb *sql.DB, userDeleter services.UserDeleter) *testController {
	return &testController{
		rc:          r,
		sqldb:       sqldb,
		userDeleter: userDeleter,
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

	userRepo := sqlrepo.NewUserRepo(c.sqldb)
	users := game.NewUserService(userRepo)

	// Delete previous user
	u, err := userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil && err != persist.ErrUserNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	if u != nil {
		c.userDeleter.DeleteUser(ctx, u.Id)
	}

	// Create new user
	usr, err := users.CreateUser(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	games, err := mama.GetActiveGames(ctx, c.sqldb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(games) == 0 {
		http.Error(w, "There are no games to join", http.StatusInternalServerError)
		return
	}

	err = mama.JoinGame(ctx, c.sqldb, usr.Id, games[0].Id)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	gcClient := gamelet.NewGameletClient(games[0].Host)
	err = gcClient.JoinGame(usr.Id, usr.Username, []string{})
	if err != nil {
		w.WriteHeader(500)
		return
	}

	a := game.AuthService{}
	token, err := a.CreateToken(usr.Id, time.Now().Add(2*time.Minute))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(&setupTestResponse{
		AccessToken: token,
		User: &user{
			Id:       usr.Id,
			Username: usr.Username,
		},
		GameId: games[0].Id,
	})
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
