package main

import (
	"context"

	"github.com/fnatte/pizza-tribes/internal/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type patch struct {
	gsPatch     *models.GameStatePatch
	sendStats   bool
	sendReports bool
}

type updateContext struct {
	context.Context
	userId  string
	gs      *models.GameState
	patch   *patch
	patches map[string]*patch
	reports map[string][]*models.Report
}

func (u *updateContext) AppendReport(userId string, report *models.Report) {
	u.reports[userId] = append(u.reports[userId], report)
	u.patches[userId].sendReports = true
}

func (u *updateContext) IncrPizzas(amount int32) {
	if u.patch.gsPatch.Resources.Pizzas == nil {
		u.patch.gsPatch.Resources.Pizzas = &wrapperspb.Int32Value{
			Value: u.gs.Resources.Pizzas,
		}
	}

	u.patch.gsPatch.Resources.Pizzas.Value = u.patch.gsPatch.Resources.Pizzas.Value + amount
	u.gs.Resources.Pizzas = u.patch.gsPatch.Resources.Pizzas.Value
}

func (u *updateContext) IncrCoins(amount int32) {
	if u.patch.gsPatch.Resources.Coins == nil {
		u.patch.gsPatch.Resources.Coins = &wrapperspb.Int32Value{
			Value: u.gs.Resources.Coins,
		}
	}

	u.patch.gsPatch.Resources.Coins.Value = u.patch.gsPatch.Resources.Coins.Value + amount
	u.gs.Resources.Coins = u.patch.gsPatch.Resources.Coins.Value
}

func (u *updateContext) IncrUneducated(amount int32) {
	if u.patch.gsPatch.Population.Uneducated == nil {
		u.patch.gsPatch.Population.Uneducated = &wrapperspb.Int32Value{
			Value: u.gs.Population.Uneducated,
		}
	}

	u.patch.gsPatch.Population.Uneducated.Value = u.patch.gsPatch.Population.Uneducated.Value + amount
	u.gs.Population.Uneducated = u.patch.gsPatch.Population.Uneducated.Value
}

func (u *updateContext) IncrChefs(amount int32) {
	if u.patch.gsPatch.Population.Chefs == nil {
		u.patch.gsPatch.Population.Chefs = &wrapperspb.Int32Value{
			Value: u.gs.Population.Chefs,
		}
	}

	u.patch.gsPatch.Population.Chefs.Value = u.patch.gsPatch.Population.Chefs.Value + amount
	u.gs.Population.Chefs = u.patch.gsPatch.Population.Chefs.Value
}

func (u *updateContext) IncrSalesmice(amount int32) {
	if u.patch.gsPatch.Population.Salesmice == nil {
		u.patch.gsPatch.Population.Salesmice = &wrapperspb.Int32Value{
			Value: u.gs.Population.Salesmice,
		}
	}

	u.patch.gsPatch.Population.Salesmice.Value = u.patch.gsPatch.Population.Salesmice.Value + amount
	u.gs.Population.Salesmice = u.patch.gsPatch.Population.Salesmice.Value
}

func (u *updateContext) IncrGuards(amount int32) {
	if u.patch.gsPatch.Population.Guards == nil {
		u.patch.gsPatch.Population.Guards = &wrapperspb.Int32Value{
			Value: u.gs.Population.Guards,
		}
	}

	u.patch.gsPatch.Population.Guards.Value = u.patch.gsPatch.Population.Guards.Value + amount
	u.gs.Population.Guards = u.patch.gsPatch.Population.Guards.Value
}

func (u *updateContext) IncrThieves(amount int32) {
	if u.patch.gsPatch.Population.Thieves == nil {
		u.patch.gsPatch.Population.Thieves = &wrapperspb.Int32Value{
			Value: u.gs.Population.Thieves,
		}
	}

	u.patch.gsPatch.Population.Thieves.Value = u.patch.gsPatch.Population.Thieves.Value + amount
	u.gs.Population.Thieves = u.patch.gsPatch.Population.Thieves.Value
}

