package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"text/template"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/rand"
	"golang.org/x/text/message"
	"gonum.org/v1/gonum/stat/distuv"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var messagePrinter = message.NewPrinter(message.MatchLanguage("en"))

const thiefReportTemplateText = `
{{if gt .SuccessfulThieves 0}}
Our heist with {{ .Thieves }} thieves on {{ .TargetUsername }}'s town was successful.
{{if gt .CaughtThieves 0}}
{{ .CaughtThieves }} thieves were caught, but {{ .SuccessfulThieves }} thieves got away with {{ .Loot | mprintf "%d" }} coins.
{{- else}}
No thieves were caught, and they got away with {{ .Loot | mprintf "%d" }} coins.
{{- end}}
{{- else}}
Our heist on {{ .TargetUsername }} was a failure. All {{ .Thieves }} thieves got caught.
{{- end}}
`
const targetReportTemplateText = `
{{if gt .SuccessfulThieves 0}}
{{if gt .CaughtThieves 0}}
{{.CaughtThieves}} thieves were caught trying to steal from our town, but {{ .SuccessfulThieves }} thieves got away with {{ .Loot | mprintf "%d" }} of our coins!
{{- else}}
It looks like someone stole {{ .Loot | mprintf "%d" }} coins from us.
{{- end}}
{{- else}}
{{ .CaughtThieves }} thieves were caught trying to steal from our town.
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

func completeSteal(ctx updateContext, tx *redis.Tx, world *internal.WorldService, travel *models.Travel, travelIndex int) (pipeFn, error) {
	gsTarget := &models.GameState{}
	x := travel.DestinationX
	y := travel.DestinationY
	gsKeyThief := fmt.Sprintf("user:%s:gamestate", ctx.userId)

	// Validate target town
	worldEntry, err := world.GetEntryXY(ctx, int(x), int(y))
	if err != nil {
		return nil, fmt.Errorf("could not find world entry: %w", err)
	}
	town := worldEntry.GetTown()
	if town == nil {
		return nil, fmt.Errorf("no town at %d, %d", x, y)
	}
	if town.UserId == ctx.userId {
		return nil, errors.New("can't steal from own town")
	}

	// Get game state of target
	gsKeyTarget := fmt.Sprintf("user:%s:gamestate", town.UserId)
	s, err := internal.RedisJsonGet(tx, ctx, gsKeyTarget, ".").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to complete steal: %w", err)
	}
	if err = protojson.Unmarshal([]byte(s), gsTarget); err != nil {
		return nil, fmt.Errorf("failed to complete steal: %w", err)
	}

	// Get username of target
	targetUsername, err := tx.HGet(ctx, fmt.Sprintf("user:%s", town.UserId), "username").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to complete steal: %w", err)
	}

	// Calculate outcome
	guards := float64(gsTarget.Population.Guards)
	thieves := float64(travel.Thieves)
	dist := distuv.Binomial{
		N:   thieves,
		P:   thieves / (thieves + guards/2),
		Src: rand.NewSource(uint64(time.Now().UnixNano())),
	}
	successfulThieves := int32(dist.Rand())
	caughtThieves := travel.Thieves - successfulThieves
	maxLoot := successfulThieves * internal.ThiefCapacity
	loot := int64(internal.MinInt32(maxLoot, gsTarget.Resources.Coins))

	// Prepare return travel - but not if all thieves got caught
	var returnTravelBytes []byte
	if successfulThieves > 0 {
		arrivalAt := internal.CalculateArrivalTime(
			travel.DestinationX, travel.DestinationY,
			ctx.gs.TownX, ctx.gs.TownY,
			internal.ThiefSpeed,
		)

		returnTravel := models.Travel{
			ArrivalAt:    arrivalAt,
			DestinationX: travel.DestinationX,
			DestinationY: travel.DestinationY,
			Returning:    true,
			Thieves:      successfulThieves,
			Coins:        loot,
		}
		returnTravelBytes, err = protojson.Marshal(&returnTravel)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal travel: %w", err)
		}

		// Update patch with return travel
		ctx.gsPatch.TravelQueue = append(ctx.gsPatch.TravelQueue, &returnTravel)
	}

	// Build reports
	tmplData := reportTemplateData{
		TargetUsername:    targetUsername,
		Loot:              loot,
		Thieves:           travel.Thieves,
		SuccessfulThieves: successfulThieves,
		CaughtThieves:     caughtThieves,
	}
	buf := new(bytes.Buffer)
	if err = thiefReportTemplate.Execute(buf, &tmplData); err != nil {
		return nil, fmt.Errorf("failed to get thief report contents: %w", err)
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
		return nil, fmt.Errorf("failed to get target report contents: %w", err)
	}
	targetReportTitle := "We have been robbed!"
	if successfulThieves == 0 {
		targetReportTitle = "We caught thieves!"
	}
	targetReport := &models.Report{
		Id:        xid.New().String(),
		CreatedAt: time.Now().UnixNano(),
		Title:     targetReportTitle,
		Content:   buf.String(),
		Unread:    true,
	}
	*ctx.sendReports = true

	return func(pipe redis.Pipeliner) error {
		// TODO: notify (send game state patch) to target user
		// Decrease coins in target town
		_, err := internal.RedisJsonNumIncrBy(
			pipe, ctx, gsKeyTarget, ".resources.coins", -loot).Result()
		if err != nil {
			return fmt.Errorf("failed to decrease coins from target town: %w", err)
		}

		if returnTravelBytes != nil {
			// Append return travel
			if err = internal.RedisJsonArrAppend(pipe, ctx, gsKeyThief,
				".travelQueue", returnTravelBytes).Err(); err != nil {
				return fmt.Errorf("failed to append new travel: %w", err)
			}
		}

		internal.SaveReport(ctx, pipe, ctx.userId, thiefReport)
		internal.SaveReport(ctx, pipe, town.UserId, targetReport)

		log.Info().Str("userId", ctx.userId).Int64("loot", loot).Msg("Steal completed")

		return nil
	}, nil
}

func completeStealReturn(ctx updateContext, tx *redis.Tx, world *internal.WorldService, travel *models.Travel, travelIndex int) (pipeFn, error) {
	gsKey := fmt.Sprintf("user:%s:gamestate", ctx.userId)

	// Update patch with coins
	if ctx.gsPatch.Resources == nil {
		ctx.gsPatch.Resources = &models.GameStatePatch_ResourcesPatch{}
	}
	if ctx.gsPatch.Resources.Coins == nil {
		ctx.gsPatch.Resources.Coins = &wrapperspb.Int32Value{}
	}
	ctx.gsPatch.Resources.Coins.Value = ctx.gsPatch.Resources.Coins.Value + int32(travel.Coins)

	// Update patch with thieves
	if ctx.gsPatch.Population == nil {
		ctx.gsPatch.Population = &models.GameStatePatch_PopulationPatch{}
	}
	if ctx.gsPatch.Population.Thieves == nil {
		ctx.gsPatch.Population.Thieves = &wrapperspb.Int32Value{}
	}
	ctx.gsPatch.Population.Thieves.Value = ctx.gsPatch.Population.Thieves.Value + travel.Thieves

	return func(pipe redis.Pipeliner) error {
		var err error
		// Increase coins with the loot
		if err = internal.RedisJsonNumIncrBy(
			pipe, ctx, gsKey, ".resources.coins", travel.Coins).Err(); err != nil {
			return fmt.Errorf("failed to increase coins with loot: %w", err)
		}

		// Add back thieves to town population
		if err = internal.RedisJsonNumIncrBy(
			pipe, ctx, gsKey, ".population.thieves",
			int64(travel.Thieves)).Err(); err != nil {
			return fmt.Errorf("failed to increase coins with loot: %w", err)
		}

		log.Info().
			Str("userId", ctx.userId).
			Int64("loot", travel.Coins).
			Msg("Steal return completed")

		return nil
	}, nil
}

func completeTravels(ctx updateContext, tx *redis.Tx, world *internal.WorldService) (pipeFn, error) {
	completedTravels := internal.GetCompletedTravels(ctx.gs)
	if len(completedTravels) == 0 {
		return nil, nil
	}

	// Update patch
	ctx.gsPatch.TravelQueue = ctx.gs.TravelQueue[len(completedTravels):]
	ctx.gsPatch.TravelQueuePatched = true

	pipeFns := []pipeFn{}

	gsKey := fmt.Sprintf("user:%s:gamestate", ctx.userId)

	// Complete travels
	for travelIndex, travel := range completedTravels {
		if travel.Returning {
			if travel.Thieves > 0 {
				fn, err := completeStealReturn(ctx, tx, world, travel, travelIndex)
				if err != nil {
					return nil, err
				}
				pipeFns = append(pipeFns, fn)
			}
		} else {
			if travel.Thieves > 0 {
				fn, err := completeSteal(ctx, tx, world, travel, travelIndex)
				if err != nil {
					return nil, err
				}
				pipeFns = append(pipeFns, fn)
			}
		}
	}

	return func(pipe redis.Pipeliner) error {
		// Remove completed travels from queue
		err := internal.RedisJsonArrTrim(
			pipe, ctx, gsKey,
			".travelQueue",
			len(completedTravels),
			math.MaxInt32,
		).Err()
		if err != nil {
			return fmt.Errorf("failed to remove completed travels: %w", err)
		}

		for i := range pipeFns {
			if err = pipeFns[i](pipe); err != nil {
				return err
			}
		}

		return nil
	}, nil
}
