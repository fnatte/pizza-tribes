package serve

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fnatte/pizza-tribes/cmd/admin/services"
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	persistsql "github.com/fnatte/pizza-tribes/internal/game/persist/sql"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	"github.com/fnatte/pizza-tribes/internal/game/ws"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type user struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type batchDeleteUserRequest struct {
	UserIds   []string `json:"userIds"`
	Usernames []string `json:"usernames"`
}

type batchCreateUserRequest struct {
	Users []createUserRequest `json:"users"`
}

type batchDeleteUserResponseItem struct {
	Id         string `json:"id"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
}

type batchDeleteUserResponse struct {
	Users []*batchDeleteUserResponseItem `json:"users"`
}

type batchCreateUserResponseItem struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
}

type batchCreateUserResponse struct {
	Users []*batchCreateUserResponseItem `json:"users"`
}

type incrCoinsRequest struct {
	Amount int32 `json:"amount"`
}

type userController struct {
	rc      redis.RedisClient
	sqldb   *sql.DB
	updater gamestate.Updater
	userDeleter services.UserDeleter
}

func contains(slice []string, item string) bool {
	for _, str := range slice {
		if str == item {
			return true
		}
	}

	return false
}

func NewUserController(r redis.RedisClient, sqldb *sql.DB, updater gamestate.Updater, userDeleter services.UserDeleter) *userController {
	return &userController{
		rc:      r,
		updater: updater,
		sqldb:   sqldb,
		userDeleter: userDeleter,
	}
}

func (c *userController) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		w.WriteHeader(400)
		return
	}

	err := c.userDeleter.DeleteUser(ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(204)
}

func (c *userController) HandleBatchDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := batchDeleteUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userRepo := persistsql.NewUserRepo(c.sqldb)

	resp := batchDeleteUserResponse{
		Users: []*batchDeleteUserResponseItem{},
	}

	userIds := req.UserIds

	for _, username := range req.Usernames {
		u, err := userRepo.GetUserByUsername(ctx, username)
		if err != nil {
			respItem := batchDeleteUserResponseItem{}
			resp.Users = append(resp.Users, &respItem)
			respItem.Id = username
			respItem.StatusText = err.Error()

			if errors.Is(err, persist.ErrUserNotFound) {
				respItem.Status = http.StatusNotFound
			} else {
				respItem.Status = http.StatusInternalServerError
			}

			continue
		}

		userIds = append(userIds, u.Id)
	}

	for _, userId := range userIds {
		respItem := batchDeleteUserResponseItem{}
		resp.Users = append(resp.Users, &respItem)

		err := c.userDeleter.DeleteUser(ctx, userId)
		if err != nil {
			respItem.Status = http.StatusInternalServerError
			respItem.StatusText = err.Error()
			continue
		}

		respItem.Id = userId
		respItem.Status = http.StatusOK
		respItem.StatusText = "OK"
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write(b)
}

func (c *userController) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := createUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userRepo := persistsql.NewUserRepo(c.sqldb)
	users := game.NewUserService(userRepo)
	gameUserRepo := persist.NewGameUserRepository(c.rc)
	gsRepo := persist.NewGameStateRepository(c.rc)
	world := game.NewWorldService(c.rc)
	leaderboard := game.NewLeaderboardService(c.rc)
	gameCtrl := game.NewGameCtrl(gsRepo, gameUserRepo, world, leaderboard, c.rc)

	u, err := users.CreateUser(ctx, req.Username, req.Password)
	if err != nil {
		err = fmt.Errorf("failed to create user: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = gameCtrl.JoinGame(ctx, u.Id, u.Username, []string{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(&user{
		Id:       u.Id,
		Username: u.Username,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func (c *userController) HandleBatchCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := batchCreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userRepo := persistsql.NewUserRepo(c.sqldb)
	users := game.NewUserService(userRepo)
	gameUserRepo := persist.NewGameUserRepository(c.rc)
	gsRepo := persist.NewGameStateRepository(c.rc)
	world := game.NewWorldService(c.rc)
	leaderboard := game.NewLeaderboardService(c.rc)
	gameCtrl := game.NewGameCtrl(gsRepo, gameUserRepo, world, leaderboard, c.rc)

	resp := batchCreateUserResponse{
		Users: []*batchCreateUserResponseItem{},
	}

	for _, user := range req.Users {
		respItem := batchCreateUserResponseItem{}
		resp.Users = append(resp.Users, &respItem)

		u, err := users.CreateUser(ctx, user.Username, user.Password)
		if err != nil {
			respItem.Username = user.Username
			respItem.Status = http.StatusInternalServerError
			respItem.StatusText = fmt.Errorf("failed to create user: %w", err).Error()
			continue
		}

		err = gameCtrl.JoinGame(ctx, u.Id, u.Username, []string{})
		if err != nil {
			respItem.Username = user.Username
			respItem.Status = http.StatusInternalServerError
			respItem.StatusText = fmt.Errorf("failed to create user: %w", err).Error()
			continue
		}

		respItem.Id = u.Id
		respItem.Username = u.Username
		respItem.Status = 200
		respItem.StatusText = "OK"
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusMultiStatus)
	w.Write(b)
}

func (c *userController) HandleShowUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		w.WriteHeader(400)
		return
	}

	userRepo := persistsql.NewUserRepo(c.sqldb)
	u, err := userRepo.GetUser(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if u == nil {
		http.Error(w, "User not found", 404)
		return
	}

	b, err := json.Marshal(&user{
		Id:       u.Id,
		Username: u.Username,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func (c *userController) HandleCompleteUserQueues(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		w.WriteHeader(400)
		return
	}

	userRepo := persist.NewGameUserRepository(c.rc)

	u, err := userRepo.GetUser(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if u == nil {
		http.Error(w, "User not found", 404)
		return
	}

	c.updater.PerformUpdate(r.Context(), userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		// Complete construction queue
		if len(gs.ConstructionQueue) > 0 {
			q := []*models.Construction{}
			for _, c := range gs.ConstructionQueue {
				c2 := proto.Clone(c).(*models.Construction)
				c2.CompleteAt = time.Now().Add(-10 * time.Second).UnixNano()
				q = append(q, c2)
			}
			tx.Users[userId].SetConstructionQueue(q)
		}

		// Complete travel queue
		if len(gs.TravelQueue) > 0 {
			q := []*models.Travel{}
			for _, c := range gs.TravelQueue {
				c2 := proto.Clone(c).(*models.Travel)
				c2.ArrivalAt = time.Now().Add(-10 * time.Second).UnixNano()
				q = append(q, c2)
			}
			tx.Users[userId].SetTravelQueue(q)
		}

		// Complete research queue
		if len(gs.ResearchQueue) > 0 {
			q := []*models.OngoingResearch{}
			for _, c := range gs.ResearchQueue {
				c2 := proto.Clone(c).(*models.OngoingResearch)
				c2.CompleteAt = time.Now().Add(-10 * time.Second).UnixNano()
				q = append(q, c2)
			}
			tx.Users[userId].SetResearchQueue(q)
		}

		// Complete train queue
		if len(gs.TrainingQueue) > 0 {
			q := []*models.Training{}
			for _, c := range gs.TrainingQueue {
				c2 := proto.Clone(c).(*models.Training)
				c2.CompleteAt = time.Now().Add(-10 * time.Second).UnixNano()
				q = append(q, c2)
			}
			tx.Users[userId].SetTrainingQueue(q)
		}

		return nil
	})

	game.SetNextUpdateTo(c.rc, r.Context(), userId, time.Now().UnixNano())

	w.WriteHeader(http.StatusNoContent)
}

func (c *userController) HandleIncrCoins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		w.WriteHeader(400)
		return
	}

	req := incrCoinsRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := c.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		log.Info().Int32("amount", req.Amount).Msg("perform update")
		tx.Users[userId].IncrCoins(req.Amount)
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ws.Send(ctx, c.rc, userId, tx.Users[userId].ToServerMessage())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send state change")
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *userController) HandleGetGameState(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		w.WriteHeader(400)
		return
	}

	gsRepo := persist.NewGameStateRepository(c.rc)
	gs, err := gsRepo.Get(ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(gs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func (c *userController) HandlePatchGameState(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Admin: HandlePatchGameState")

	ctx := r.Context()
	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		w.WriteHeader(400)
		return
	}

	req := models.GameStatePatch{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = protojson.Unmarshal(b, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gsRepo := persist.NewGameStateRepository(c.rc)
	err = gsRepo.Patch(ctx, userId, req.GameState, req.PatchMask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update leaderboard
	if contains(req.PatchMask.Paths, "resources") || contains(req.PatchMask.Paths, "resources.coins") {
		coins := int64(req.GameState.Resources.Coins)
		leaderboard := game.NewLeaderboardService(c.rc)
		if err = leaderboard.UpdateUser(ctx, userId, coins); err != nil {
			log.Error().Err(err).Msg("Failed to update leaderboard")
		}
	}

	err = ws.Send(ctx, c.rc, userId, req.ToServerMessage())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send state change")
	}

	game.SetNextUpdateTo(c.rc, r.Context(), userId, time.Now().UnixNano())

	w.WriteHeader(http.StatusNoContent)
}

func (c *userController) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	var userIds []string
	userRepo := persistsql.NewUserRepo(c.sqldb)
	if username != "" {
		u, err := userRepo.GetUserByUsername(r.Context(), username)
		if err != nil {
			if errors.Is(err, persist.ErrUserNotFound) {
				userIds = []string{}
			} else {
				http.Error(w, err.Error(), 500)
				return
			}
		}
		if u != nil {
			userIds = []string{u.Id}
		}
	} else {
		var err error
		users, err := userRepo.GetAllUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		for _, u := range users {
			userIds = append(userIds, u.Id)
		}
	}

	b, err := json.Marshal(userIds)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}

func (c *userController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/users/{userId}/completeQueues", c.HandleCompleteUserQueues).Methods("POST")
	r.HandleFunc("/users/{userId}/incrCoins", c.HandleIncrCoins).Methods("POST")
	r.HandleFunc("/users/{userId}/gameState", c.HandlePatchGameState).Methods("PATCH")
	r.HandleFunc("/users/{userId}/gameState", c.HandleGetGameState).Methods("GET")

	r.HandleFunc("/users/batch", c.HandleBatchCreateUser).Methods("POST")
	r.HandleFunc("/users/batch", c.HandleBatchDeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{userId}", c.HandleDeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{userId}", c.HandleShowUser).Methods("GET")
	r.HandleFunc("/users", c.HandleCreateUser).Methods("POST")
	r.HandleFunc("/users", c.HandleListUsers).Methods("GET")

	return r
}
