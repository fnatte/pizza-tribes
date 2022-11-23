package main

import (
	"time"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/mtime"
	"github.com/rs/zerolog/log"
)

type extrapolateChanges struct {
	timestamp int64
	coins     int32
	pizzas    int32
}

func extrapolate(userId string, gs *models.GameState, tx *gamestate.GameTx, worldState *models.WorldState, globalDemandScore float64, userCount int64, speed float64) error {
	changes := calculateExtrapolateChanges(gs, worldState, globalDemandScore, userCount, speed)

	u := tx.Users[userId]
	u.SetCoins(gs.Resources.Coins + changes.coins)
	u.SetPizzas(gs.Resources.Pizzas + changes.pizzas)
	u.SetTimestamp(changes.timestamp)

	return nil
}

func calculateExtrapolateChanges(gs *models.GameState, worldState *models.WorldState, globalDemandScore float64, userCount int64, speed float64) extrapolateChanges {
	// No changes if there are no population
	if game.CountTownPopulation(gs) == 0 {
		return extrapolateChanges{}
	}

	now := time.Now()

	then := gs.Timestamp
	if then <= 0 {
		then = now.Unix()
	}

	dt := float64(now.Unix() - then) * speed

	rush, offpeak := mtime.GetRush(then, now.Unix())

	stats := game.CalculateStats(gs, globalDemandScore, worldState, userCount)

	demand := int32(stats.DemandOffpeak*float64(offpeak) +
		stats.DemandRushHour*float64(rush))

	pizzasProduced := int32(stats.PizzasProducedPerSecond * dt)
	pizzasAvailable := gs.Resources.Pizzas + pizzasProduced

	maxSellsByMice := int32(stats.MaxSellsByMicePerSecond * dt)
	pizzasSold := game.MinInt32(demand,
		game.MinInt32(maxSellsByMice, pizzasAvailable))

	pizzaPrice := gs.GetValidPizzaPrice()

	log.Debug().
		Int32("pizzasAvailable", pizzasAvailable).
		Int32("pizzasProduced", pizzasProduced).
		Int32("maxSellsByMice", maxSellsByMice).
		Int32("pizzasSold", pizzasSold).
		Int32("demand", demand).
		Int32("price", pizzaPrice).
		Interface("stats", &stats).
		Msg("Game state update")

	return extrapolateChanges{
		coins:     pizzasSold * pizzaPrice,
		pizzas:    pizzasProduced - pizzasSold,
		timestamp: now.Unix(),
	}
}
