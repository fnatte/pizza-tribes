package internal

import (
	"strconv"
	"time"

	. "github.com/fnatte/pizza-tribes/internal/models/gamedata"
	. "github.com/fnatte/pizza-tribes/internal/models/gamestate"
)

func NewInt64(i int64) *int64    { return &i }
func NewString(s string) *string { return &s }

func CountBuildings(gs *GameState) (counts map[string]int32) {
	counts = map[string]int32{}
	for _, lot := range gs.Lots {
		counts[string(lot.Building)] = counts[string(lot.Building)] + 1
	}
	return counts
}

func CountBuildingsUnderConstruction(gs *GameState) (counts map[string]int32) {
	counts = map[string]int32{}
	for _, c := range gs.ConstructionQueue {
		counts[string(c.Building)] = counts[string(c.Building)] + 1
	}
	return counts
}

func Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
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

func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func GetBuildingInfo(id Building) *BuildingInfo {
	for _, b := range FullGameData.Buildings {
		if b.ID == string(id) {
			return &b
		}
	}

	return nil
}

func CountMaxEmployed(gs *GameState) (counts map[string]int32) {
	counts = map[string]int32{}
	for _, lot := range gs.Lots {
		info := GetBuildingInfo(lot.Building)
		if info != nil && info.LevelInfos[lot.Level].Employer != nil {
			counts[string(lot.Building)] = counts[string(lot.Building)] +
				info.LevelInfos[lot.Level].Employer.MaxWorkforce
		}
	}
	return counts
}

func CountMaxPopulation(gs *GameState) (count int32) {
	for _, lot := range gs.Lots {
		info := GetBuildingInfo(lot.Building)
		if info != nil && info.LevelInfos[lot.Level].Residence != nil {
			count = count +
				info.LevelInfos[lot.Level].Residence.Beds
		}
	}
	return
}

func CountTownPopulation(population *GameStatePopulation) int32 {
	if population == nil {
		return 0
	}

	return (population.Uneducated +
		population.Chefs +
		population.Salesmice +
		population.Guards +
		population.Thieves +
		population.Publicists)
}

func CountTravellingPopulation(travelQueue []Travel) int32 {
	var count int32 = 0
	for _, t := range travelQueue {
		count = count + t.Thieves
	}

	return count
}

func CountTrainingPopulation(trainingQueue []Training) int32 {
	var count int32 = 0
	for _, t := range trainingQueue {
		count = count + t.Amount
	}

	return count
}

func CountAllPopulation(gs *GameState) int32 {
	if gs == nil {
		return 0
	}

	return CountTownPopulation(&gs.Population) +
		CountTravellingPopulation(gs.TravelQueue) +
		CountTrainingPopulation(gs.TrainingQueue)

}

func GetCompletedTravels(gs *GameState) (res []*Travel) {
	now := time.Now().UnixNano()

	for _, t := range gs.TravelQueue {
		arrivalAt, err := strconv.ParseInt(t.ArrivalAt, 10, 64)
		if err != nil {
			panic("Falied to parse ArrivalAt as int64")
		}

		if arrivalAt > now {
			break
		}

		res = append(res, &t)
	}

	return res
}

func HasBuildingMinLevel(gs *GameState, b Building, minLvl int) bool {
	for _, lot := range gs.Lots {
		if lot.Building == b {
			if lot.Level >= int32(minLvl - 1) {
				return true
			}
		}
	}

	return false
}

