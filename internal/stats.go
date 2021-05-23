package internal

import (
	. "github.com/fnatte/pizza-tribes/internal/models"
)

const CHEF_PIZZAS_PER_SECOND = 0.2
const SALESMICE_SELLS_PER_SECOND = 0.5
const DEMAND_BASE = 0.2
const DEMAND_RUSH_HOUR_BONUS = 0.55

func calculateTasteScore(gs *GameState) float64 {
	score := 1.0

	if gs.HasDiscovery(ResearchDiscovery_DURUM_WHEAT) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscovery_DOUBLE_ZERO_FLOUR) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscovery_SAN_MARZANO_TOMATOES) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscovery_OCIMUM_BASILICUM) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscovery_EXTRA_VIRGIN) {
		score = score + 0.05
	}
	if gs.HasDiscovery(ResearchDiscovery_MASONRY_OVEN) {
		score = score + 0.1
	}

	return score
}

func calculatePopularity(gs *GameState) float64 {
	popularityBonus := 1.0

	if gs.HasDiscovery(ResearchDiscovery_WEBSITE) {
		popularityBonus = popularityBonus + 0.1
	}
	if gs.HasDiscovery(ResearchDiscovery_MOBILE_APP) {
		popularityBonus = popularityBonus + 0.1
	}

	tasteScore := calculateTasteScore(gs)

	return float64(gs.Population.Publicists) * 5.0 * popularityBonus * tasteScore
}

func calculateSalesBonus(gs *GameState) float64 {
	bonus := 1.0

	if gs.HasDiscovery(ResearchDiscovery_DIGITAL_ORDERING_SYSTEM) {
		bonus = bonus + 0.2
	}

	return bonus
}

func calculateBakeBonus(gs *GameState) float64 {
	bonus := 1.0

	if gs.HasDiscovery(ResearchDiscovery_GAS_OVEN) {
		bonus = bonus + 0.1
	}
	if gs.HasDiscovery(ResearchDiscovery_HYBRID_OVEN) {
		bonus = bonus + 0.1
	}

	return bonus
}

func CalculateStats(gs *GameState) *Stats {
	// No changes if there are no population
	if gs.Population == nil {
		return &Stats{}
	}

	popularity := calculatePopularity(gs)
	demandOffpeak := DEMAND_BASE * popularity
	demandRushHour := (DEMAND_BASE + DEMAND_RUSH_HOUR_BONUS) * popularity

	maxEmployed := CountMaxEmployed(gs)

	employedChefs := MinInt32(gs.Population.Chefs, maxEmployed[int32(Building_KITCHEN)])
	pizzasProducedPerSecond := float64(employedChefs) *
		CHEF_PIZZAS_PER_SECOND *
		calculateBakeBonus(gs)

	employedSalesmice := MinInt32(gs.Population.Salesmice, maxEmployed[int32(Building_SHOP)])
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
