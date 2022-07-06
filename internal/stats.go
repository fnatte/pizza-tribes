package internal

import (
	. "github.com/fnatte/pizza-tribes/internal/models/gamestate"
)

const CHEF_PIZZAS_PER_SECOND = 0.2
const SALESMICE_SELLS_PER_SECOND = 0.5
const DEMAND_BASE = 0.2
const DEMAND_RUSH_HOUR_BONUS = 0.55

func calculateTasteScore(gs *GameState) float64 {
	score := 1.0

	if gs.HasDiscovery(ResearchDiscoveryDurumWheat) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscoveryDoubleZeroFlour) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscoverySanMarzanoTomatoes) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscoveryOcimumBasilicum) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscoveryExtraVirgin) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscoveryMasonryOven) {
		score = score + 0.1
	}

	return score
}

func calculatePopularity(gs *GameState) float64 {
	popularityBonus := 1.0

	if gs.HasDiscovery(ResearchDiscoveryWebsite) {
		popularityBonus = popularityBonus + 0.1
	}
	if gs.HasDiscovery(ResearchDiscoveryMobileApp) {
		popularityBonus = popularityBonus + 0.1
	}

	tasteScore := calculateTasteScore(gs)

	return (3 + float64(gs.Population.Publicists)*2) * 5.0 * popularityBonus * tasteScore
}

func calculateSalesBonus(gs *GameState) float64 {
	bonus := 1.0

	if gs.HasDiscovery(ResearchDiscoveryDigitalOrderingSystem) {
		bonus = bonus + 0.2
	}

	return bonus
}

func calculateBakeBonus(gs *GameState) float64 {
	bonus := 1.0

	if gs.HasDiscovery(ResearchDiscoveryGasOven) {
		bonus = bonus + 0.1
	}
	if gs.HasDiscovery(ResearchDiscoveryHybridOven) {
		bonus = bonus + 0.1
	}

	return bonus
}

func CalculateStats(gs *GameState) *Stats {
	popularity := calculatePopularity(gs)
	demandOffpeak := DEMAND_BASE * popularity
	demandRushHour := (DEMAND_BASE + DEMAND_RUSH_HOUR_BONUS) * popularity

	maxEmployed := CountMaxEmployed(gs)

	employedChefs := MinInt32(gs.Population.Chefs, maxEmployed[string(BuildingKitchen)])
	pizzasProducedPerSecond := float64(employedChefs) *
		CHEF_PIZZAS_PER_SECOND *
		calculateBakeBonus(gs)

	employedSalesmice := MinInt32(gs.Population.Salesmice, maxEmployed[string(BuildingShop)])
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
