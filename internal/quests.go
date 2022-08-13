package internal

import (
	"github.com/fnatte/pizza-tribes/internal/models"
)

func GetAvailableQuestIds(gs *models.GameState) []string {
	quests := []string{"1"}

	if q, ok := gs.Quests["1"]; ok && q.Completed {
		quests = append(quests, "2")
	}

	if q, ok := gs.Quests["2"]; ok && q.Completed {
		quests = append(quests, "3", "4")
	}

	if q, ok := gs.Quests["3"]; ok && q.Completed {
		if q, ok := gs.Quests["4"]; ok && q.Completed {
			quests = append(quests, "5", "12")
		}
	}

	if q, ok := gs.Quests["5"]; ok && q.Completed {
		quests = append(quests, "6")
	}

	if q, ok := gs.Quests["6"]; ok && q.Completed {
		quests = append(quests, "7")
	}

	if q, ok := gs.Quests["7"]; ok && q.Completed {
		quests = append(quests, "8", "9")
	}

	if q, ok := gs.Quests["8"]; ok && q.Completed {
		if q, ok := gs.Quests["9"]; ok && q.Completed {
			quests = append(quests, "10", "11")
		}
	}
	if q, ok := gs.Quests["10"]; ok && q.Completed {
		quests = append(quests, "13")
	}

	return quests
}

func GetNewCompletedQuests(gs *models.GameState) []string {
	solved := []string{}

	for qid, q := range gs.Quests {
		// Already solved
		if q.Completed {
			continue
		}

		switch qid {
		case "1":
			if HasBuildingMinLevel(gs, models.Building_KITCHEN, 1) &&
				HasBuildingMinLevel(gs, models.Building_SHOP, 1) {
				solved = append(solved, qid)
			}
			break
		case "2":
			if HasBuildingMinLevel(gs, models.Building_HOUSE, 1) &&
				HasBuildingMinLevel(gs, models.Building_SCHOOL, 1) {
				solved = append(solved, qid)
			}
			break
		case "3":
			e := CountTownPopulationEducations(gs)
			if e[models.Education_CHEF] >= 1 && e[models.Education_SALESMOUSE] >= 1 {
				solved = append(solved, qid)
			}
			break
		case "4":
			// Change name quest is solved on handling changeName client message
			break
		case "5":
			if HasBuildingMinLevel(gs, models.Building_HOUSE, 2) {
				solved = append(solved, qid)
			}
			break
		case "6":
			// "Check out help page" is solved using special message
			break
		case "7":
			if HasBuildingMinLevel(gs, models.Building_KITCHEN, 2) &&
				HasBuildingMinLevel(gs, models.Building_SHOP, 2) {
				solved = append(solved, qid)
			}
			break
		case "8":
			if CountEmployed(gs) >= 7 {
				solved = append(solved, qid)
			}
			break
		case "9":
			// "Stats quest" is solved using special message
			break
		case "10":
			e := CountTownPopulationEducations(gs)
			if e[models.Education_GUARD] >= 1 {
				solved = append(solved, qid)
			}
			break
		case "11":
			// "Leaderboard quest" is solved using special message
			break
		case "12":
			for _, lot := range gs.Lots {
				if lot.Streak >= 2 {
					solved = append(solved, qid)
				}
			}
			break
		case "13":
			if HasBuildingMinLevel(gs, models.Building_TOWN_CENTRE, 2) {
				solved = append(solved, qid)
			}
			break
		}
	}

	return solved
}
