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
	sqldb *sql.DB
	userDeleter services.UserDeleter
}

func NewTestController(r redis.RedisClient, sqldb *sql.DB, userDeleter services.UserDeleter) *testController {
	return &testController{
		rc: r,
		sqldb: sqldb,
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
	gameUserRepo := persist.NewGameUserRepository(c.rc)
	gsRepo := persist.NewGameStateRepository(c.rc)
	world := game.NewWorldService(c.rc)
	users := game.NewUserService(userRepo)
	leaderboard := game.NewLeaderboardService(c.rc)
	gameCtrl := game.NewGameCtrl(gsRepo, gameUserRepo, world, leaderboard)

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

	err = gameCtrl.JoinGame(ctx, usr.Id, usr.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
