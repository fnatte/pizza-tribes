package internal

import (
	. "github.com/fnatte/pizza-tribes/internal/models2"
)

const CHEF_PIZZAS_PER_SECOND = 0.2
const SALESMICE_SELLS_PER_SECOND = 0.5
const DEMAND_BASE = 0.2
const DEMAND_RUSH_HOUR_BONUS = 0.55

func calculateTasteScore(gs *GameState) float64 {
	score := 1.0

	if gs.HasDiscovery(DurumWheat) {
		score = score + 0.05
	}
	if gs.HasDiscovery(DoubleZeroFlour) {
		score = score + 0.05
	}
	if gs.HasDiscovery(SANMarzanoTomatoes) {
		score = score + 0.05
	}
	if gs.HasDiscovery(OcimumBasilicum) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ExtraVirgin) {
		score = score + 0.05
	}
	if gs.HasDiscovery(MasonryOven) {
		score = score + 0.1
	}

	return score
}

func calculatePopularity(gs *GameState) float64 {
	popularityBonus := 1.0

	if gs.HasDiscovery(Website) {
		popularityBonus = popularityBonus + 0.1
	}
	if gs.HasDiscovery(MobileApp) {
		popularityBonus = popularityBonus + 0.1
	}

	tasteScore := calculateTasteScore(gs)

	return (3 + float64(gs.Population.Publicists)*2) * 5.0 * popularityBonus * tasteScore
}

func calculateSalesBonus(gs *GameState) float64 {
	bonus := 1.0

	if gs.HasDiscovery(DigitalOrderingSystem) {
		bonus = bonus + 0.2
	}

	return bonus
}

func calculateBakeBonus(gs *GameState) float64 {
	bonus := 1.0

	if gs.HasDiscovery(GasOven) {
		bonus = bonus + 0.1
	}
	if gs.HasDiscovery(HybridOven) {
		bonus = bonus + 0.1
	}

	return bonus
}

func CalculateStats(gs *GameState) *Stats {
	popularity := calculatePopularity(gs)
	demandOffpeak := DEMAND_BASE * popularity
	demandRushHour := (DEMAND_BASE + DEMAND_RUSH_HOUR_BONUS) * popularity

	maxEmployed := CountMaxEmployed(gs)

	employedChefs := Min(gs.Population.Chefs, maxEmployed[string(Kitchen)])
	pizzasProducedPerSecond := float64(employedChefs) *
		CHEF_PIZZAS_PER_SECOND *
		calculateBakeBonus(gs)

	employedSalesmice := Min(gs.Population.Salesmice, maxEmployed[string(Shop)])
	maxSellsByMicePerSecond := float64(employedSalesmice) *
		SALESMICE_SELLS_PER_SECOND *
		calculateSalesBonus(gs)

	return &Stats{
		EmployedChefs:           employedChefs,
		EmployedSalesmice:       employedSalesmice,
		MaxSellsByMicePerSecond: maxSellsByMicePerSecond,
		PizzasProducedPerSecond: pizzasProducedPerSecond,
		DemandOffpeak:           demandOffpeak,
		DemandRushHour:          demandRushHour,
	}
}
