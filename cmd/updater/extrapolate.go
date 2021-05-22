package main

import (
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

func extrapolate(ctx updateContext, tx *redis.Tx) (error) {
	changes := calculateExtrapolateChanges(ctx.gs)

	// Update patch
	ctx.IncrCoins(changes.coins)
	ctx.IncrPizzas(changes.pizzas)
	ctx.patch.gsPatch.Timestamp = &wrapperspb.Int64Value{Value: changes.timestamp}

	return nil
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
		coins:     pizzasSold,
		pizzas:    pizzasProduced - pizzasSold,
		timestamp: now.Unix(),
	}
}
