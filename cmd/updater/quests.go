package main

import (
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
)

func completeQuests(userId string, gs *models.GameState, tx *gamestate.GameTx) error {
	for _, questId := range game.GetNewCompletedQuests(gs) {
		tx.Users[userId].SetQuestCompleted(questId)
	}

	for _, questId := range game.GetAvailableQuestIds(gs) {
		if gs.Quests[questId] != nil {
			continue
		}

		tx.Users[userId].SetQuestAvailable(questId)
	}

	return nil
}
