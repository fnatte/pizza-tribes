package gamestate

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
)

type GameTx_User struct {
	Gs      *models.GameState
	GsPatch *persist.Patch
	Reports []*models.Report

	NextUpdateInvalidated bool
	StatsInvalidated bool
}

type GameTx struct {
	Users map[string]GameTx_User
}

func (u *GameTx_User) SetCoins(val int32) {
	if u.Gs.Resources.Coins != val {
		u.Gs.Resources.Coins = val
		u.GsPatch.Replace("/resources/coins", val)
	}
}

func (u *GameTx_User) SetPizzas(val int32) {
	if u.Gs.Resources.Coins != val {
		u.Gs.Resources.Coins = val
		u.GsPatch.Replace("/resources/pizzas", val)
	}
}

func (u *GameTx_User) SetTappedAt(lotId string, val int64) {
	if u.Gs.Lots[lotId].TappedAt != val {
		u.Gs.Lots[lotId].TappedAt = val
		u.GsPatch.Replace(fmt.Sprintf("/lots/%s/tappedAt", lotId), val)
	}
}

func (u *GameTx_User) SetTaps(lotId string, val int32) {
	if u.Gs.Lots[lotId].Taps != val {
		u.Gs.Lots[lotId].Taps = val
		u.GsPatch.Replace(fmt.Sprintf("/lots/%s/taps", lotId), val)
	}
}

func (u *GameTx_User) SetStreak(lotId string, val int32) {
	if u.Gs.Lots[lotId].Streak != val {
		u.Gs.Lots[lotId].Streak = val
		u.GsPatch.Replace(fmt.Sprintf("/lots/%s/streak", lotId), val)
	}
}

func (u *GameTx_User) SetUneducated(val int32) {
	if u.Gs.Population.Uneducated != val {
		u.Gs.Population.Uneducated = val
		u.GsPatch.Replace("/population/uneducated", val)
	}
}

func (u *GameTx_User) SetMouseIsBeingEducated(mouseId string, val bool) {
	if u.Gs.Mice[mouseId].IsBeingEducated != val {
		u.Gs.Mice[mouseId].IsBeingEducated = val
		u.GsPatch.Replace(fmt.Sprintf("/mice/%s/isBeingEducated", mouseId), val)
	}
}

func (u *GameTx_User) AppendTrainingQueue(val *models.Training) {
	u.Gs.TrainingQueue = append(u.Gs.TrainingQueue, val)
	u.GsPatch.Add("/trainingQueue/-", val)
	u.NextUpdateInvalidated = true
}

func NewGameTx(userId string, gs *models.GameState) *GameTx {
	return &GameTx{
		Users: map[string]GameTx_User{
			userId: {
				Reports: []*models.Report{},
				Gs:      gs,
				GsPatch: &persist.Patch{
					Ops: []*persist.Operation{},
				},
			},
		},
	}
}

type updateContext struct {
	context.Context
	userId string
	tx     *GameTx
}

func NewUpdateContext(ctx context.Context, userId string, gs *models.GameState) *updateContext {
	return &updateContext{
		Context: ctx,
		userId:  userId,
		tx:      NewGameTx(userId, gs),
	}
}
