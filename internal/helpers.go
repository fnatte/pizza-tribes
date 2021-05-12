package internal

import (
	"google.golang.org/protobuf/encoding/protojson"
)

func NewInt64(i int64) *int64    { return &i }
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

func CountMaxEmployed(gs *GameState) (counts map[int32]int32) {
	counts = map[int32]int32{}
	for _, lot := range gs.Lots {
		info := FullGameData.Buildings[int32(lot.Building)]
		if info != nil && info.LevelInfos[lot.Level].Employer != nil {
			counts[int32(lot.Building)] = counts[int32(lot.Building)] +
				info.LevelInfos[lot.Level].Employer.MaxWorkforce
		}
	}
	return counts
}

func CountMaxPopulation(gs *GameState) (count int32) {
	for _, lot := range gs.Lots {
		info := FullGameData.Buildings[int32(lot.Building)]
		if info != nil && info.LevelInfos[lot.Level].Residence != nil {
			count = count +
				info.LevelInfos[lot.Level].Residence.Beds
		}
	}
	return
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
	for _, t := range travelQueue {
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
