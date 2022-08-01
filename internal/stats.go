package internal

import (
	. "github.com/fnatte/pizza-tribes/internal/models"
)

const CHEF_PIZZAS_PER_SECOND = 0.2
const SALESMICE_SELLS_PER_SECOND = 0.5
const DEMAND_BASE = 0.2
const DEMAND_RUSH_HOUR_BONUS = 0.55

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

func CalculateStats(gs *GameState, globalDemandScore float64, worldState *WorldState) *Stats {
	// No changes if there are no population
	if CountTownPopulation(gs) == 0 {
		return &Stats{}
	}

	e := CountTownPopulationEducations(gs)

	marketShare := CalculateDemandScore(gs) / globalDemandScore
	marketDemandBase := marketShare * CalculateGlobalDemand(worldState)
	marketDemandOffpeak := marketDemandBase
	marketDemandRushHour := marketDemandBase * 2

	demandOffpeak := DEMAND_BASE + marketDemandOffpeak
	demandRushHour := (DEMAND_BASE + DEMAND_RUSH_HOUR_BONUS) + marketDemandRushHour

	maxEmployed := CountMaxEmployed(gs)

	employedChefs := MinInt32(e[Education_CHEF], maxEmployed[int32(Building_KITCHEN)])
	pizzasProducedPerSecond := float64(employedChefs) *
		CHEF_PIZZAS_PER_SECOND *
		calculateBakeBonus(gs)

	employedSalesmice := MinInt32(e[Education_SALESMOUSE], maxEmployed[int32(Building_SHOP)])
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
		MarketShare:             marketShare,
	}
}
