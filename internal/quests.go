package internal

import (
	. "github.com/fnatte/pizza-tribes/internal/models/gamestate"
)

func GetAvailableQuestIds(gs *GameState) []string {
	quests := []string{"1"}

	if q, ok := gs.Quests["1"]; ok && q.Completed {
		quests = append(quests, "2")
	}

	if q, ok := gs.Quests["2"]; ok && q.Completed {
		quests = append(quests, "3")
	}

	if q, ok := gs.Quests["3"]; ok && q.Completed {
		quests = append(quests, "4")
	}

	if q, ok := gs.Quests["4"]; ok && q.Completed {
		quests = append(quests, "5")
	}

	if q, ok := gs.Quests["5"]; ok && q.Completed {
		quests = append(quests, "6")
	}

	if q, ok := gs.Quests["6"]; ok && q.Completed {
		quests = append(quests, "7")
	}

	return quests
}

func GetNewCompletedQuests(gs *GameState) []string {
	solved := []string{}

	for qid, q := range gs.Quests {
		// Already solved
		if q.Completed {
			continue
		}

		switch qid {
		case "1":
			if HasBuildingMinLevel(gs, BuildingKitchen, 1) &&
				HasBuildingMinLevel(gs, BuildingShop, 1) {
				solved = append(solved, qid)
			}
			break
		case "2":
			if HasBuildingMinLevel(gs, BuildingHouse, 1) &&
				HasBuildingMinLevel(gs, BuildingSchool, 1) {
				solved = append(solved, qid)
			}
			break
		case "3":
			if gs.Population.Chefs >= 1 && gs.Population.Salesmice >= 1 {
				solved = append(solved, qid)
			}
			break
		case "4":
			// Change name quest is solved on handling changeName client message
			break
		case "5":
			if HasBuildingMinLevel(gs, BuildingHouse, 2) {
				solved = append(solved, qid)
			}
			break
		case "6":
			// "Check out help page" is solved using special message
			break
		case "7":
			if HasBuildingMinLevel(gs, BuildingKitchen, 2) &&
				HasBuildingMinLevel(gs, BuildingShop, 2) {
				solved = append(solved, qid)
			}
			break
		}
	}

	return solved
}
