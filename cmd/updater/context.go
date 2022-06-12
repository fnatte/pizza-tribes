package main

import (
	"context"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/rs/xid"
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

func (u *updateContext) setEducatedMouse(education models.Education) {
	for id, m := range u.gs.Mice {
		if m.IsBeingEducated {
			if u.patch.gsPatch.Mice == nil {
				u.patch.gsPatch.Mice = map[string]*models.GameStatePatch_MousePatch{}
			}
			if u.patch.gsPatch.Mice[id] == nil {
				u.patch.gsPatch.Mice[id] = &models.GameStatePatch_MousePatch{}
			}
			mp := u.patch.gsPatch.Mice[id]

			m.IsBeingEducated = false
			mp.IsBeingEducated = wrapperspb.Bool(false)
			m.IsEducated = true
			mp.IsEducated = wrapperspb.Bool(true)
			m.Education = education
			mp.Education = &models.GameStatePatch_EducationPatch{
				Value: education,
			}

			break
		}
	}
}

func (u *updateContext) appendNewMouse() {
	id := xid.New().String()
	newMouse := &models.Mouse{
		Name:       GetNewMouseName(u.gs.Mice),
		IsEducated: false,
		Education:  0,
	}
	u.gs.Mice[id] = newMouse
	if u.patch.gsPatch.Mice == nil {
		u.patch.gsPatch.Mice = map[string]*models.GameStatePatch_MousePatch{}
	}
	u.patch.gsPatch.Mice[id] = newMouse.ToPatch(true)
}

func (u *updateContext) removeMouse(isEducated bool, education models.Education) {
	for id, m := range u.gs.Mice {
		if (!m.IsEducated && !isEducated) || m.Education == education {
			u.gs.Mice[id] = nil
			if u.patch.gsPatch.Mice == nil {
				u.patch.gsPatch.Mice = map[string]*models.GameStatePatch_MousePatch{}
			}
			u.patch.gsPatch.Mice[id] = nil
		}
	}
}

func (u *updateContext) IncrUneducated(amount int32) {
	if amount == 0 {
		return
	}

	if u.patch.gsPatch.Population.Uneducated == nil {
		u.patch.gsPatch.Population.Uneducated = &wrapperspb.Int32Value{
			Value: u.gs.Population.Uneducated,
		}
	}

	u.patch.gsPatch.Population.Uneducated.Value = u.patch.gsPatch.Population.Uneducated.Value + amount
	u.gs.Population.Uneducated = u.patch.gsPatch.Population.Uneducated.Value

	if amount > 0 {
		for n := 0; n < int(amount); n++ {
			u.appendNewMouse()
		}
	} else {
		for n := 0; n < int(-amount); n++ {
			u.removeMouse(false, 0)
		}
	}

}

func (u *updateContext) IncrChefs(amount int32) {
	if u.patch.gsPatch.Population.Chefs == nil {
		u.patch.gsPatch.Population.Chefs = &wrapperspb.Int32Value{
			Value: u.gs.Population.Chefs,
		}
	}

	u.patch.gsPatch.Population.Chefs.Value = u.patch.gsPatch.Population.Chefs.Value + amount
	u.gs.Population.Chefs = u.patch.gsPatch.Population.Chefs.Value

	if amount > 0 {
		for n := 0; n < int(amount); n++ {
			u.setEducatedMouse(models.Education_CHEF)
		}
	} else {
		for n := 0; n < int(-amount); n++ {
			u.removeMouse(true, models.Education_CHEF)
		}
	}
}

func (u *updateContext) IncrSalesmice(amount int32) {
	if u.patch.gsPatch.Population.Salesmice == nil {
		u.patch.gsPatch.Population.Salesmice = &wrapperspb.Int32Value{
			Value: u.gs.Population.Salesmice,
		}
	}

	u.patch.gsPatch.Population.Salesmice.Value = u.patch.gsPatch.Population.Salesmice.Value + amount
	u.gs.Population.Salesmice = u.patch.gsPatch.Population.Salesmice.Value

	if amount > 0 {
		for n := 0; n < int(amount); n++ {
			u.setEducatedMouse(models.Education_SALESMOUSE)
		}
	} else {
		for n := 0; n < int(-amount); n++ {
			u.removeMouse(true, models.Education_SALESMOUSE)
		}
	}
}

func (u *updateContext) IncrGuards(amount int32) {
	if u.patch.gsPatch.Population.Guards == nil {
		u.patch.gsPatch.Population.Guards = &wrapperspb.Int32Value{
			Value: u.gs.Population.Guards,
		}
	}

	u.patch.gsPatch.Population.Guards.Value = u.patch.gsPatch.Population.Guards.Value + amount
	u.gs.Population.Guards = u.patch.gsPatch.Population.Guards.Value

	if amount > 0 {
		for n := 0; n < int(amount); n++ {
			u.setEducatedMouse(models.Education_GUARD)
		}
	} else {
		for n := 0; n < int(-amount); n++ {
			u.removeMouse(true, models.Education_GUARD)
		}
	}
}

func (u *updateContext) IncrThieves(amount int32) {
	if u.patch.gsPatch.Population.Thieves == nil {
		u.patch.gsPatch.Population.Thieves = &wrapperspb.Int32Value{
			Value: u.gs.Population.Thieves,
		}
	}

	u.patch.gsPatch.Population.Thieves.Value = u.patch.gsPatch.Population.Thieves.Value + amount
	u.gs.Population.Thieves = u.patch.gsPatch.Population.Thieves.Value

	if amount > 0 {
		for n := 0; n < int(amount); n++ {
			u.setEducatedMouse(models.Education_THIEF)
		}
	} else {
		for n := 0; n < int(-amount); n++ {
			u.removeMouse(true, models.Education_THIEF)
		}
	}
}

func (u *updateContext) IncrPublicists(amount int32) {
	if u.patch.gsPatch.Population.Publicists == nil {
		u.patch.gsPatch.Population.Publicists = &wrapperspb.Int32Value{
			Value: u.gs.Population.Publicists,
		}
	}

	u.patch.gsPatch.Population.Publicists.Value = u.patch.gsPatch.Population.Publicists.Value + amount
	u.gs.Population.Publicists = u.patch.gsPatch.Population.Publicists.Value

	if amount > 0 {
		for n := 0; n < int(amount); n++ {
			u.setEducatedMouse(models.Education_PUBLICIST)
		}
	} else {
		for n := 0; n < int(-amount); n++ {
			u.removeMouse(true, models.Education_PUBLICIST)
		}
	}
}
