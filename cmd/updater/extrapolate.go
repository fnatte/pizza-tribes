package main

import (
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/mtime"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type extrapolateChanges struct {
	timestamp int64
	coins     int32
	pizzas    int32
}

func extrapolate(ctx updateContext, tx *redis.Tx) (pipeFn, error) {
	gsKey := fmt.Sprintf("user:%s:gamestate", ctx.userId)
	changes := calculateExtrapolateChanges(ctx.gs)

	// Update patch
	r := ctx.gsPatch.Resources
	if r.Coins == nil {
		r.Coins = &wrapperspb.Int32Value{}
	}
	if r.Pizzas == nil {
		r.Pizzas = &wrapperspb.Int32Value{}
	}
	r.Coins.Value = r.Coins.Value + changes.coins
	r.Pizzas.Value = r.Pizzas.Value + changes.pizzas
	ctx.gsPatch.Timestamp = &wrapperspb.Int64Value{Value: changes.timestamp}

	return func(pipe redis.Pipeliner) error {
		// Write timestamp
		err := internal.RedisJsonSet(
			pipe, ctx, gsKey, ".timestamp", changes.timestamp).Err()
		if err != nil {
			return err
		}

		// Write coins
		err = internal.RedisJsonSet(
			pipe, ctx, gsKey, ".resources.coins", changes.coins).Err()
		if err != nil {
			return err
		}

		// Write pizzas
		err = internal.RedisJsonSet(
			pipe, ctx, gsKey, ".resources.pizzas", changes.pizzas).Err()
		if err != nil {
			return err
		}

		return nil
	}, nil
}

func calculateExtrapolateChanges(gs *models.GameState) extrapolateChanges {
	// No changes if there are no population
	if gs.Population == nil {
		return extrapolateChanges{}
	}

	now := time.Now()
	rush, offpeak := mtime.GetRush(gs.Timestamp, now.Unix())
	dt := float64(now.Unix() - gs.Timestamp)

	stats := internal.CalculateStats(gs)

	demand := int32(stats.DemandOffpeak*float64(offpeak) +
		stats.DemandRushHour*float64(rush))

	pizzasProduced := int32(stats.PizzasProducedPerSecond * dt)
	pizzasAvailable := gs.Resources.Pizzas + pizzasProduced

	maxSellsByMice := int32(stats.MaxSellsByMicePerSecond * dt)
	pizzasSold := internal.MinInt32(demand,
		internal.MinInt32(maxSellsByMice, pizzasAvailable))

	log.Debug().
		Int32("pizzasProduced", pizzasProduced).
		Int32("maxSellsByMice", maxSellsByMice).
		Int32("pizzasSold", pizzasSold).
		Msg("Game state update")

	return extrapolateChanges{
		coins:     gs.Resources.Coins + pizzasSold,
		pizzas:    gs.Resources.Pizzas + pizzasProduced - pizzasSold,
		timestamp: now.Unix(),
	}
}
