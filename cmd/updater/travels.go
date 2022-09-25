package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"text/template"
	"time"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/message"
)

var messagePrinter = message.NewPrinter(message.MatchLanguage("en"))

const thiefReportTemplateText = `
{{if gt .SuccessfulThieves 0}}
Our heist with {{ .Thieves }} thieves on *{{ .TargetUsername }}'s* town was successful.
{{if gt .CaughtThieves 0}}
**{{ .CaughtThieves }} thieves** were **caught**, but {{ .SuccessfulThieves }} thieves got away with **{{ .Loot | mprintf "%d" }} coins**.
{{- else}}
No thieves were caught, and they got away with **{{ .Loot | mprintf "%d" }} coins**.
{{- end}}
{{- else}}
Our heist on *{{ .TargetUsername }}* was a failure. All **{{ .Thieves }}** thieves got **caught**.
{{- end}}
{{if gt .SleepingGuards 0}}
{{ .SleepingGuards }} guards were sleeping during the heist and have probably been fired.
{{- end}}
`
const targetReportTemplateText = `
{{if gt .SuccessfulThieves 0}}
{{if gt .CaughtThieves 0}}
**{{.CaughtThieves}}** thieves were **caught** trying to steal from our town, but {{ .SuccessfulThieves }} thieves got away with **{{ .Loot | mprintf "%d" }} coins**!
{{- else}}
It looks like someone stole **{{ .Loot | mprintf "%d" }} coins** from us.
{{- end}}
{{- else}}
**{{ .CaughtThieves }}** thieves were **caught** trying to steal from our town.
{{- end}}
{{if gt .SleepingGuards 0}}
{{ .SleepingGuards }} guards were fired for sleeping during their shift.
{{- end}}
`

var thiefReportTemplate *template.Template
var targetReportTemplate *template.Template

type reportTemplateData struct {
	TargetUsername    string
	Loot              int64
	Thieves           int32
	SuccessfulThieves int32
	CaughtThieves     int32
	SleepingGuards    int32
}

type pipeFn func(redis.Pipeliner) error

func init() {
	tmplFuncMap := template.FuncMap{
		"mprintf": messagePrinter.Sprintf,
	}

	thiefReportTemplate = template.Must(template.New("root").
		Funcs(tmplFuncMap).
		Parse(thiefReportTemplateText))

	targetReportTemplate = template.Must(template.New("root").
		Funcs(tmplFuncMap).
		Parse(targetReportTemplateText))
}

func getThiefEvadeBonus(gs *models.GameState) float64 {
	bonus := 0.0

	if gs.HasDiscovery(models.ResearchDiscovery_TIP_TOE) {
		bonus += 0.25
	}

	if gs.HasDiscovery(models.ResearchDiscovery_SHADOW_EXPERT) {
		bonus += 0.50
	}

	return bonus
}

func getThiefCapacityBonus(gs *models.GameState) float64 {
	bonus := 0.0

	if gs.HasDiscovery(models.ResearchDiscovery_BIG_POCKETS) {
		bonus += 0.25
	}

	if gs.HasDiscovery(models.ResearchDiscovery_THIEVES_FAVORITE_BAG) {
		bonus += 0.50
	}

	return bonus
}

func getGuardEfficiencyBonus(gs *models.GameState) float64 {
	bonus := 0.0

	if gs.HasDiscovery(models.ResearchDiscovery_NIGHTS_WATCH) {
		bonus += 0.15
	}

	if gs.HasDiscovery(models.ResearchDiscovery_TRIP_WIRE) {
		bonus += 0.20
	}

	if gs.HasDiscovery(models.ResearchDiscovery_CARDIO) {
		bonus += 0.25
	}

	if gs.HasDiscovery(models.ResearchDiscovery_LASER_ALARM) {
		bonus += 0.30
	}

	return bonus
}

func getGuardAwarenessBonus(gs *models.GameState) float64 {
	bonus := 0.0

	if gs.HasDiscovery(models.ResearchDiscovery_COFFEE) {
		bonus += 0.25
	}

	return bonus
}

func completeSteal(ctx context.Context, userId string, gs *models.GameState, tx *gamestate.GameTx, r redis.RedisClient, world *game.WorldService, travel *models.Travel, travelIndex int) error {
	gsTarget := &models.GameState{}
	x := travel.DestinationX
	y := travel.DestinationY

	// Validate target town
	worldEntry, err := world.GetEntryXY(ctx, int(x), int(y))
	if err != nil {
		return fmt.Errorf("could not find world entry: %w", err)
	}
	town := worldEntry.GetTown()
	if town == nil {
		return fmt.Errorf("no town at %d, %d", x, y)
	}
	if town.UserId == userId {
		return errors.New("can't steal from own town")
	}

	// Get game state of target
	gsKeyTarget := fmt.Sprintf("user:%s:gamestate", town.UserId)
	s, err := redis.RedisJsonGet(r, ctx, gsKeyTarget, ".").Result()
	if err != nil {
		return fmt.Errorf("failed to complete steal: %w", err)
	}
	if err = protojson.Unmarshal([]byte(s), gsTarget); err != nil {
		return fmt.Errorf("failed to complete steal: %w", err)
	}

	// Get username of target
	targetUsername, err := r.HGet(ctx, fmt.Sprintf("user:%s", town.UserId), "username").Result()
	if err != nil {
		return fmt.Errorf("failed to complete steal: %w", err)
	}

	// Calculate outcome
	targetEducations := game.CountTownPopulationEducations(gsTarget)
	outcome := game.CalculateHeist(game.Heist{
		Guards:      targetEducations[models.Education_GUARD],
		Thieves:     travel.Thieves,
		TargetCoins: gsTarget.Resources.Coins,
		ThiefEvadeBonus: getThiefEvadeBonus(gs),
		ThiefCapacityBonus: getThiefCapacityBonus(gs),
		GuardEfficiencyBonus: getThiefEvadeBonus(gsTarget),
		GuardAwarenessBonus: getGuardAwarenessBonus(gsTarget),
	}, rand.NewSource(uint64(time.Now().UnixNano())))

	// Make caught thiefs uneducated
	for n := int32(0); n < outcome.CaughtThieves; n++ {
		for k, m := range tx.Users[userId].Gs.Mice {
			if m.IsEducated && m.Education == models.Education_THIEF {
				tx.Users[userId].SetMouseUneducated(k)
				break
			}
		}
	}

	// Prepare return travel - but not if all thieves got caught
	if outcome.SuccessfulThieves > 0 {
		arrivalAt := game.CalculateArrivalTime(
			travel.DestinationX, travel.DestinationY,
			gs.TownX, gs.TownY,
			game.GetThiefSpeed(gs),
		)

		returnTravel := models.Travel{
			ArrivalAt:    arrivalAt,
			DestinationX: travel.DestinationX,
			DestinationY: travel.DestinationY,
			Returning:    true,
			Thieves:      outcome.SuccessfulThieves,
			Coins:        outcome.Loot,
		}
		if err != nil {
			return fmt.Errorf("failed to marshal travel: %w", err)
		}

		// Update patch with return travel
		tx.Users[userId].SetTravelQueue(append(gs.TravelQueue, &returnTravel))
	}

	// Build reports
	tmplData := reportTemplateData{
		TargetUsername:    targetUsername,
		Loot:              outcome.Loot,
		Thieves:           travel.Thieves,
		SuccessfulThieves: outcome.SuccessfulThieves,
		CaughtThieves:     outcome.CaughtThieves,
		SleepingGuards:    outcome.SleepingGuards,
	}
	buf := new(bytes.Buffer)
	if err = thiefReportTemplate.Execute(buf, &tmplData); err != nil {
		return fmt.Errorf("failed to get thief report contents: %w", err)
	}
	thiefReport := &models.Report{
		Id:        xid.New().String(),
		CreatedAt: time.Now().UnixNano(),
		Title:     "Thief report",
		Content:   buf.String(),
		Unread:    true,
	}
	buf = new(bytes.Buffer)
	if err = targetReportTemplate.Execute(buf, &tmplData); err != nil {
		return fmt.Errorf("failed to get target report contents: %w", err)
	}
	targetReportTitle := "We have been robbed!"
	if outcome.SuccessfulThieves == 0 {
		targetReportTitle = "We caught thieves!"
	}
	targetReport := &models.Report{
		Id:        xid.New().String(),
		CreatedAt: time.Now().UnixNano(),
		Title:     targetReportTitle,
		Content:   buf.String(),
		Unread:    true,
	}

	// Prepare patch to target user (whoms coins was stoled)
	tx.InitUser(town.UserId, gsTarget)
	tx.Users[town.UserId].IncrCoins(-int32(outcome.Loot))

	for n := int32(0); n < outcome.SleepingGuards; n++ {
		for k, m := range tx.Users[town.UserId].Gs.Mice {
			if m.IsEducated && m.Education == models.Education_GUARD {
				tx.Users[town.UserId].SetMouseUneducated(k)
				break
			}
		}
	}

	// Append reports to patch
	tx.Users[userId].AppendReport(thiefReport)
	tx.Users[town.UserId].AppendReport(targetReport)

	return nil
}

func completeStealReturn(userId string, gs *models.GameState, tx *gamestate.GameTx, world *game.WorldService, travel *models.Travel, travelIndex int) error {
	u := tx.Users[userId]
	u.IncrCoins(int32(travel.Coins))

	log.Info().
		Str("userId", userId).
		Int64("loot", travel.Coins).
		Msg("Steal return completed")

	return nil
}

func completeTravels(ctx context.Context, userId string, gs *models.GameState, tx *gamestate.GameTx, r redis.RedisClient, world *game.WorldService) error {
	completedTravels := game.GetCompletedTravels(gs)
	if len(completedTravels) == 0 {
		return nil
	}

	// Update patch
	u := tx.Users[userId]
	u.SetTravelQueue(gs.TravelQueue[len(completedTravels):])

	// Complete travels
	for travelIndex, travel := range completedTravels {
		if travel.Returning {
			if travel.Thieves > 0 {
				err := completeStealReturn(userId, gs, tx, world, travel, travelIndex)
				if err != nil {
					return err
				}
			}
		} else {
			if travel.Thieves > 0 {
				err := completeSteal(ctx, userId, gs, tx, r, world, travel, travelIndex)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
