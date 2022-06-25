package main

import (
	"time"

	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
)

func completeResearchs(userId string, gs *models.GameState, tx *gamestate.GameTx) error {
	completedResearchs := getCompletedResearchs(gs)

	// Exit early if there are no completed researchs
	if len(completedResearchs) == 0 {
		return nil
	}

	// Update patch
	u := tx.Users[userId]
	u.SetResearchQueue(gs.ResearchQueue[len(completedResearchs):])

	for _, r := range completedResearchs {
		if !gs.HasDiscovery(r.Discovery) {
			u.AppendDiscovery(r.Discovery)
		}
	}

	// Completion of research can affect the stats
	u.StatsInvalidated = true

	return nil
}

func getCompletedResearchs(gs *models.GameState) (res []*models.OngoingResearch) {
	now := time.Now().UnixNano()

	for _, t := range gs.ResearchQueue {
		if t.CompleteAt > now {
			break
		}

		res = append(res, t)
	}

	return res
}
