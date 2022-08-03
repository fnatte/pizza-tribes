package internal

import (
	"time"

	"github.com/fnatte/pizza-tribes/internal/models"
	. "github.com/fnatte/pizza-tribes/internal/models"
)

func NewInt64(i int64) *int64    { return &i }
func NewString(s string) *string { return &s }

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

func CountEmployed(gs *GameState) (count int32) {
	e := CountTownPopulationEducations(gs)
	m := CountMaxEmployed(gs)

	for education, n := range e {
		info := FullGameData.Educations[int32(education)]
		if info.Employer != nil {
			max := m[int32(*info.Employer)]
			if n < max {
				count += n
			} else {
				count += max
			}
		}
	}

	return count
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

func CountUneducated(gs *GameState) int32 {
	n := int32(0)
	for _, m := range gs.Mice {
		if !m.IsEducated {
			n++
		}
	}

	return n
}

func CountEducated(gs *GameState) int32 {
	n := int32(0)
	for _, m := range gs.Mice {
		if m.IsEducated {
			n++
		}
	}

	return n
}

func CountTownPopulation(gs *GameState) int32 {
	return int32(len(gs.Mice))
}

func CountTravellingPopulation(travelQueue []*Travel) int32 {
	var count int32 = 0
	for _, t := range travelQueue {
		count = count + t.Thieves
	}

	return count
}

func CountTrainingPopulation(trainingQueue []*Training) int32 {
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

	return CountTownPopulation(gs) +
		CountTravellingPopulation(gs.TravelQueue) +
		CountTrainingPopulation(gs.TrainingQueue)
}

func CountTownPopulationEducations(gs *GameState) map[models.Education]int32 {
	res := map[models.Education]int32{}
	for _, m := range gs.Mice {
		if m.IsEducated {
			res[m.Education] = res[m.Education] + 1
		}
	}

	for _, t := range gs.TravelQueue {
		if t.Thieves > 0 {
			res[Education_THIEF] = res[Education_THIEF] - t.Thieves
		}
	}

	return res
}

func GetCompletedTravels(gs *GameState) (res []*Travel) {
	now := time.Now().UnixNano()

	for _, t := range gs.TravelQueue {
		if t.ArrivalAt > now {
			break
		}

		res = append(res, t)
	}

	return res
}

func HasBuildingMinLevel(gs *GameState, b Building, minLvl int) bool {
	if gs.Lots == nil {
		return false
	}

	for _, lot := range gs.Lots {
		if lot != nil && lot.Building == b {
			if lot.Level >= int32(minLvl-1) {
				return true
			}
		}
	}

	return false
}

func FindMouseIdWithEducation(mice map[string]*Mouse, edu Education) string {
	for id, m := range mice {
		if m.IsEducated && m.Education == edu {
			return id
		}
	}

	return ""
}

var validLotIds = []string{
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
	"10",
	"11",
}

func IsValidLotId(lotId string) bool {
	for _, x := range validLotIds {
		if x == lotId {
			return true
		}
	}
	return false
}

