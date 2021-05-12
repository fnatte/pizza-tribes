package internal

import (
	"google.golang.org/protobuf/encoding/protojson"
)

func NewInt64(i int64) *int64 { return &i }
func NewString(s string) *string { return &s }

var protojsonu = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

func CountBuildings(gs *GameState) (counts map[int32]int32) {
	counts = map[int32]int32{}
	for _, lot := range gs.Lots {
		counts[int32(lot.Building)] = counts[int32(lot.Building)] + 1
	}
	return counts
}

func CountBuildingsUnderConstruction(gs *GameState) (counts map[int32]int32) {
	counts = map[int32]int32{}
	for _, c := range gs.ConstructionQueue {
		counts[int32(c.Building)] = counts[int32(c.Building)] + 1
	}
	return counts
}

func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MinInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func CountMaxEmployed(buildingCount map[int32]int32) (counts map[int32]int32) {
	counts = map[int32]int32{}
	buildings := FullGameData.Buildings
	for k := range Building_name {
		employer := buildings[k].Employer
		if employer != nil {
			maxWorkforce := employer.MaxWorkforce
			counts[k] = buildingCount[k] * maxWorkforce
		}
	}
	return counts
}

func CountTownPopulation(population *GameState_Population) int32 {
	if population == nil {
		return 0
	}

	return (population.Uneducated +
		population.Chefs +
		population.Salesmice +
		population.Guards +
		population.Thieves)
}

func CountTravellingPopulation(travelQueue []*Travel) int32 {
	var count int32 = 0
	for _, t := range(travelQueue) {
		count = count + t.Thieves
	}

	return count
}

func CountAllPopulation(gs *GameState) int32 {
	if gs == nil {
		return 0
	}

	return CountTownPopulation(gs.Population) +
		CountTravellingPopulation(gs.TravelQueue)

}
