package main

import (
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
)

func completeQuests(userId string, gs *models.GameState, tx *gamestate.GameTx) error {
	for _, questId := range internal.GetNewCompletedQuests(gs) {
		tx.Users[userId].SetQuestCompleted(questId)
	}

	for _, questId := range internal.GetAvailableQuestIds(gs) {
		if gs.Quests[questId] != nil {
			continue
		}

		tx.Users[userId].SetQuestAvailable(questId)
	}

	return nil
}
