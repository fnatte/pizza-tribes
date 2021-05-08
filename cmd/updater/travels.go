package main

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type pipeFn func(redis.Pipeliner) error

func completeSteal(ctx context.Context, tx *redis.Tx, world *internal.WorldService, gs *internal.GameState, gsPatch *internal.GameStatePatch, userId string, travel *internal.Travel, travelIndex int) (pipeFn, error) {
	gsTarget := &internal.GameState{}
	x := travel.DestinationX
	y := travel.DestinationY
	gsKeyThief := fmt.Sprintf("user:%s:gamestate", userId)

	// Validate target town
	worldEntry, err := world.GetEntryXY(ctx, int(x), int(y))
	if err != nil {
		return nil, fmt.Errorf("could not find world entry: %w", err)
	}
	town := worldEntry.GetTown()
	if town == nil {
		return nil, fmt.Errorf("no town at %d, %d", x, y)
	}
	if town.UserId == userId {
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

	// Calculate loot
	maxLoot := travel.Thieves * 3_000
	loot := int64(internal.MinInt32(maxLoot, gsTarget.Resources.Coins))

	// Calculate arrival time
	arrivalAt := internal.CalculateArrivalTime(
		travel.DestinationX, travel.DestinationY,
		gs.TownX, gs.TownY,
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
	gsPatch.TravelQueue = append(gsPatch.TravelQueue, &returnTravel)

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

		log.Info().Str("userId", userId).Int64("loot", loot).Msg("Steal completed")

		return nil
	}, nil
}

func completeStealReturn(ctx context.Context, tx *redis.Tx, world *internal.WorldService, gs *internal.GameState, gsPatch *internal.GameStatePatch, userId string, travel *internal.Travel, travelIndex int) (pipeFn, error) {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)

	// Update patch with coins
	if gsPatch.Resources == nil {
		gsPatch.Resources = &internal.GameStatePatch_ResourcesPatch{}
	}
	if gsPatch.Resources.Coins == nil {
		gsPatch.Resources.Coins = &wrapperspb.Int32Value{}
	}
	gsPatch.Resources.Coins.Value = gsPatch.Resources.Coins.Value + int32(travel.Coins)

	// Update patch with thieves
	if gsPatch.Population == nil {
		gsPatch.Population = &internal.GameStatePatch_PopulationPatch{}
	}
	if gsPatch.Population.Thieves == nil {
		gsPatch.Population.Thieves = &wrapperspb.Int32Value{}
	}
	gsPatch.Population.Thieves.Value = gsPatch.Population.Thieves.Value + travel.Thieves

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
			Str("userId", userId).
			Int64("loot", travel.Coins).
			Msg("Steal return completed")

		return nil
	}, nil
}

func completeTravels(ctx context.Context, tx *redis.Tx, world *internal.WorldService, userId string, gs *internal.GameState, gsPatch *internal.GameStatePatch) (pipeFn, error) {
	completedTravels := gs.GetCompletedTravels()
	if len(completedTravels) == 0 {
		return nil, nil
	}

	// Update patch
	gsPatch.TravelQueue = gs.TravelQueue[len(completedTravels):]
	gsPatch.TravelQueuePatched = true

	pipeFns := []pipeFn{}

	gsKey := fmt.Sprintf("user:%s:gamestate", userId)

	// Complete travels
	for travelIndex, travel := range completedTravels {
		if travel.Returning {
			if travel.Thieves > 0 {
				fn, err := completeStealReturn(ctx, tx, world, gs, gsPatch, userId, travel, travelIndex)
				if err != nil {
					return nil, err
				}
				pipeFns = append(pipeFns, fn)
			}
		} else {
			if travel.Thieves > 0 {
				fn, err := completeSteal(ctx, tx, world, gs, gsPatch, userId, travel, travelIndex)
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
