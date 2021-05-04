package internal

const CHEF_PIZZAS_PER_SECOND = 0.2
const SALESMICE_SELLS_PER_SECOND = 0.5
const DEMAND_BASE = 0.2
const DEMAND_RUSH_HOUR_BONUS = 0.55

func CalculateStats(gs *GameState) *Stats {
	// No changes if there are no population
	if gs.Population == nil {
		return &Stats{}
	}

	popularity := float64(CountPopulation(gs))
	demandOffpeak := DEMAND_BASE * popularity
	demandRushHour := (DEMAND_BASE + DEMAND_RUSH_HOUR_BONUS) * popularity

	buildingCount := CountBuildings(gs)
	maxEmployed := CountMaxEmployed(buildingCount)

	employedChefs := minInt32(gs.Population.Chefs, maxEmployed[int32(Building_KITCHEN)])
	pizzasProducedPerSecond := float64(employedChefs) * CHEF_PIZZAS_PER_SECOND

	employedSalesmice := minInt32(gs.Population.Salesmice, maxEmployed[int32(Building_SHOP)])
	maxSellsByMicePerSecond := float64(employedSalesmice) * SALESMICE_SELLS_PER_SECOND

	return &Stats{
		EmployedChefs:           employedChefs,
		EmployedSalesmice:       employedSalesmice,
		MaxSellsByMicePerSecond: maxSellsByMicePerSecond,
		PizzasProducedPerSecond: pizzasProducedPerSecond,
		DemandOffpeak:           demandOffpeak,
		DemandRushHour:          demandRushHour,
	}
}
