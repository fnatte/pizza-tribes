package serve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
	"github.com/fnatte/pizza-tribes/internal/protojson"
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
	rc      internal.RedisClient
	updater gamestate.Updater
}

func NewUserController(r internal.RedisClient, updater gamestate.Updater) *userController {
	return &userController{
		rc:      r,
		updater: updater,
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

	userRepo := persist.NewUserRepository(c.rc)
	gsRepo := persist.NewGameStateRepository(c.rc)
	world := internal.NewWorldService(c.rc)

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
		err = fmt.Errorf("failed to create user: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := userRepo.GetUser(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to get user: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gsRepo := persist.NewGameStateRepository(c.rc)
	gs := models.NewGameState()
	gsRepo.Save(ctx, id, gs)
	if err != nil {
		err = fmt.Errorf("failed to save game state: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	world := internal.NewWorldService(c.rc)
	x, y, err := world.AcquireTown(ctx, id)
	if err != nil {
		err = fmt.Errorf("failed to acquire town: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gs.TownX = int32(x)
	gs.TownY = int32(y)
	gsRepo.Patch(ctx, id, gs, &models.PatchMask{
		Paths: []string{"townX", "townY"},
	})
	if err != nil {
		err = fmt.Errorf("failed to patch acquired town: %w", err)
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

	err = ws.Send(ctx, c.rc, userId, req.ToServerMessage())
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
	r.HandleFunc("/users/{userId}/gameState", c.HandlePatchGameState).Methods("PATCH")
	r.HandleFunc("/users/{userId}/gameState", c.HandleGetGameState).Methods("GET")

	r.HandleFunc("/users/{userId}", c.HandleDeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{userId}", c.HandleShowUser).Methods("GET")
	r.HandleFunc("/users", c.HandleCreateUser).Methods("POST")
	r.HandleFunc("/users", c.HandleListUsers).Methods("GET")

	return r
}
