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
}

type GameTx struct {
	Users map[string]GameTx_User
}

func (u *GameTx_User) SetCoins(val int32) {
	u.Gs.Resources.Coins = val
	u.GsPatch.Replace("/resources/coins", val)
}

func (u *GameTx_User) SetPizzas(val int32) {
	u.Gs.Resources.Coins = val
	u.GsPatch.Replace("/resources/pizzas", val)
}

func (u *GameTx_User) SetTappedAt(lotId string, val int64) {
	u.Gs.Lots[lotId].TappedAt = val
	u.GsPatch.Replace(fmt.Sprintf("/lots/%s/tappedAt", lotId), val)
}

func (u *GameTx_User) SetTaps(lotId string, val int32) {
	u.Gs.Lots[lotId].Taps = val
	u.GsPatch.Replace(fmt.Sprintf("/lots/%s/taps", lotId), val)
}

func (u *GameTx_User) SetStreak(lotId string, val int32) {
	u.Gs.Lots[lotId].Streak = val
	u.GsPatch.Replace(fmt.Sprintf("/lots/%s/streak", lotId), val)
}

func NewGameTx(userId string, gs *models.GameState) *GameTx {
	return &GameTx{
		Users: map[string]GameTx_User{
			userId: {
				Reports: []*models.Report{},
				Gs:      gs,
				GsPatch: &persist.Patch{
					Ops: []persist.Operation{},
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
