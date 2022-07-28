package serve

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
	"github.com/fnatte/pizza-tribes/internal/ws"
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

type incrCoinsRequest struct {
	Amount int32 `json:"amount"`
}

type userController struct {
	rc internal.RedisClient
}

func NewUserController(r internal.RedisClient) *userController {
	return &userController{
		rc: r,
	}
}

func (c *userController) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		w.WriteHeader(400)
		return
	}

	userRepo := persist.NewUserRepository(c.rc)
	userRepo.DeleteUser(r.Context(), userId)

	w.WriteHeader(204)
}

func (c *userController) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := createUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userRepo := persist.NewUserRepository(c.rc)
	id, err := userRepo.CreateUser(ctx, req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := userRepo.GetUser(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gsRepo := persist.NewGameStateRepository(c.rc)
	gs := models.NewGameState()
	gsRepo.Save(ctx, id, gs)
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

func (c *userController) HandleShowUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		w.WriteHeader(400)
		return
	}

	userRepo := persist.NewUserRepository(c.rc)
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

	userRepo := persist.NewUserRepository(c.rc)
	gsRepo := persist.NewGameStateRepository(c.rc)
	reportsRepo := persist.NewReportsRepository(c.rc)

	u, err := userRepo.GetUser(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if u == nil {
		http.Error(w, "User not found", 404)
		return
	}

	gamestate.PerformUpdate(r.Context(), gsRepo, reportsRepo, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		if len(gs.ConstructionQueue) > 0 {
			q := []*models.Construction{}
			for _, c := range gs.ConstructionQueue {
				c2 := proto.Clone(c).(*models.Construction)
				c2.CompleteAt = time.Now().Add(-10 * time.Second).UnixNano()
				q = append(q, c2)
			}
			tx.Users[userId].SetConstructionQueue(q)
		}
		return nil
	})

	internal.SetNextUpdateTo(c.rc, r.Context(), userId, time.Now().UnixNano())

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

	gsRepo := persist.NewGameStateRepository(c.rc)
	reportsRepo := persist.NewReportsRepository(c.rc)

	tx, err := gamestate.PerformUpdate(ctx, gsRepo, reportsRepo, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
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

func (c *userController) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	var users []string
	userRepo := persist.NewUserRepository(c.rc)
	if username != "" {
		userId, err := userRepo.FindUser(r.Context(), username)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if userId == "" {
			users = []string{}
		} else {
			users = []string{userId}
		}
	} else {
		var err error
		users, err = userRepo.GetAllUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	b, err := json.Marshal(users)
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

	r.HandleFunc("/users/{userId}", c.HandleDeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{userId}", c.HandleShowUser).Methods("GET")
	r.HandleFunc("/users", c.HandleCreateUser).Methods("POST")
	r.HandleFunc("/users", c.HandleListUsers).Methods("GET")

	return r
}
