package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"text/template"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const thiefReportTemplateText = `
Our heist with {{ .Thieves }} thieves on {{ .TargetUsername }}'s town was successful.
We got away with {{ .Loot }} coins.
`
const targetReportTemplateText = `
It looks like someone stole {{ .Loot }} coins from us.
`
var thiefReportTemplate *template.Template
var targetReportTemplate *template.Template

type reportTemplateData struct {
	TargetUsername string
	Loot int64
	Thieves int32
}

type pipeFn func(redis.Pipeliner) error

func init() {
	var err error
	thiefReportTemplate, err = template.New("root").Parse(thiefReportTemplateText)
	if err != nil {
		panic(err)
	}

	targetReportTemplate, err = template.New("root").Parse(targetReportTemplateText)
	if err != nil {
		panic(err)
	}
}

func completeSteal(ctx updateContext, tx *redis.Tx, world *internal.WorldService, travel *internal.Travel, travelIndex int) (pipeFn, error) {
	gsTarget := &internal.GameState{}
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
	b, err := internal.RedisJsonGet(tx, ctx, gsKeyTarget, ".").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to complete steal: %w", err)
	}
	if err = gsTarget.LoadProtoJson([]byte(b)); err != nil {
		return nil, fmt.Errorf("failed to complete steal: %w", err)
	}

	// Get username of target
	targetUsername, err := tx.HGet(ctx, fmt.Sprintf("user:%s", town.UserId), "username").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to complete steal: %w", err)
	}

	// Calculate loot
	maxLoot := travel.Thieves * 3_000
	loot := int64(internal.MinInt32(maxLoot, gsTarget.Resources.Coins))

	// Calculate arrival time
	arrivalAt := internal.CalculateArrivalTime(
		travel.DestinationX, travel.DestinationY,
		ctx.gs.TownX, ctx.gs.TownY,
		internal.ThiefSpeed,
	)

	returnTravel := internal.Travel{
		ArrivalAt:    arrivalAt,
		DestinationX: travel.DestinationX,
		DestinationY: travel.DestinationY,
		Returning:    true,
		Thieves:      travel.Thieves,
		Coins:        loot,
	}
	returnTravelBytes, err := protojson.Marshal(&returnTravel)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal travel: %w", err)
	}

	// Update patch with return travel
	ctx.gsPatch.TravelQueue = append(ctx.gsPatch.TravelQueue, &returnTravel)

	// Build reports
	tmplData := reportTemplateData{
		TargetUsername: targetUsername,
		Loot: loot,
		Thieves: travel.Thieves,
	}
	buf := new(bytes.Buffer)
	if err = thiefReportTemplate.Execute(buf, &tmplData); err != nil {
		return nil, fmt.Errorf("failed to get thief report contents: %w", err)
	}
	thiefReport := &internal.Report{
		Id: xid.New().String(),
		CreatedAt: time.Now().UnixNano(),
		Title: "Thief report",
		Content: buf.String(),
		Unread: true,
	}
	buf = new(bytes.Buffer)
	if err = targetReportTemplate.Execute(buf, &tmplData); err != nil {
		return nil, fmt.Errorf("failed to get target report contents: %w", err)
	}
	targetReport := &internal.Report{
		Id: xid.New().String(),
		CreatedAt: time.Now().UnixNano(),
		Title: "We have been robbed!",
		Content: buf.String(),
		Unread: true,
	}
	*ctx.sendReports = true

	return func(pipe redis.Pipeliner) error {
		// TODO: notify (send game state patch) to target user
		// Decrease coins in target town
		_, err := internal.RedisJsonNumIncrBy(
			pipe, ctx, gsKeyTarget, ".resources.coins", loot).Result()
		if err != nil {
			return fmt.Errorf("failed to decrease coins from target town: %w", err)
		}

		// Append return travel
		if err = internal.RedisJsonArrAppend(pipe, ctx, gsKeyThief,
			".travelQueue", returnTravelBytes).Err(); err != nil {
			return fmt.Errorf("failed to append new travel: %w", err)
		}

		internal.SaveReport(ctx, pipe, ctx.userId, thiefReport)
		internal.SaveReport(ctx, pipe, town.UserId, targetReport)

		log.Info().Str("userId", ctx.userId).Int64("loot", loot).Msg("Steal completed")

		return nil
	}, nil
}

func completeStealReturn(ctx updateContext, tx *redis.Tx, world *internal.WorldService, travel *internal.Travel, travelIndex int) (pipeFn, error) {
	gsKey := fmt.Sprintf("user:%s:gamestate", ctx.userId)

	// Update patch with coins
	if ctx.gsPatch.Resources == nil {
		ctx.gsPatch.Resources = &internal.GameStatePatch_ResourcesPatch{}
	}
	if ctx.gsPatch.Resources.Coins == nil {
		ctx.gsPatch.Resources.Coins = &wrapperspb.Int32Value{}
	}
	ctx.gsPatch.Resources.Coins.Value = ctx.gsPatch.Resources.Coins.Value + int32(travel.Coins)

	// Update patch with thieves
	if ctx.gsPatch.Population == nil {
		ctx.gsPatch.Population = &internal.GameStatePatch_PopulationPatch{}
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
	completedTravels := ctx.gs.GetCompletedTravels()
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
