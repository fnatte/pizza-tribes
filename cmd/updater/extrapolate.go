package main

import (
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/mtime"
	"github.com/rs/zerolog/log"
)

type extrapolateChanges struct {
	timestamp int64
	coins     int32
	pizzas    int32
}

func extrapolate(userId string, gs *models.GameState, tx *gamestate.GameTx) error {
	changes := calculateExtrapolateChanges(gs)

	u := tx.Users[userId]
	u.SetCoins(gs.Resources.Coins + changes.coins)
	u.SetPizzas(gs.Resources.Pizzas + changes.pizzas)
	u.SetTimestamp(changes.timestamp)

	return nil
}

func calculateExtrapolateChanges(gs *models.GameState) extrapolateChanges {
	// No changes if there are no population
	if internal.CountTownPopulation(gs) == 0 {
		return extrapolateChanges{}
	}

	now := time.Now()

	then := gs.Timestamp
	if then <= 0 {
		then = now.Unix()
	}

	dt := float64(now.Unix() - then)

	rush, offpeak := mtime.GetRush(then, now.Unix())

	stats := internal.CalculateStats(gs)

	demand := int32(stats.DemandOffpeak*float64(offpeak) +
		stats.DemandRushHour*float64(rush))

	pizzasProduced := int32(stats.PizzasProducedPerSecond * dt)
	pizzasAvailable := gs.Resources.Pizzas + pizzasProduced

	maxSellsByMice := int32(stats.MaxSellsByMicePerSecond * dt)
	pizzasSold := internal.MinInt32(demand,
		internal.MinInt32(maxSellsByMice, pizzasAvailable))

	log.Debug().
		Int32("pizzasAvailable", pizzasAvailable).
		Int32("pizzasProduced", pizzasProduced).
		Int32("maxSellsByMice", maxSellsByMice).
		Int32("pizzasSold", pizzasSold).
		Int32("demand", demand).
		Interface("stats", &stats).
		Msg("Game state update")

	return extrapolateChanges{
		coins:     pizzasSold,
		pizzas:    pizzasProduced - pizzasSold,
		timestamp: now.Unix(),
	}
}
