package main

import (
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func completeQuests(ctx updateContext) error {
	for _, q := range internal.GetNewCompletedQuests(ctx.gs) {
		if ctx.patch.gsPatch.Quests == nil {
			ctx.patch.gsPatch.Quests = map[string]*models.GameStatePatch_QuestStatePatch{}
		}

		ctx.patch.gsPatch.Quests[q] = &models.GameStatePatch_QuestStatePatch{
			Completed: wrapperspb.Bool(true),
		}
		ctx.gs.Quests[q].Completed = true
	}

	for _, qid := range internal.GetAvailableQuestIds(ctx.gs) {
		if ctx.gs.Quests[qid] != nil {
			continue
		}

		q := &models.QuestState{}
		if ctx.gs.Quests == nil {
			ctx.gs.Quests = map[string]*models.QuestState{}
		}
		if ctx.patch.gsPatch.Quests == nil {
			ctx.patch.gsPatch.Quests = map[string]*models.GameStatePatch_QuestStatePatch{}
		}
		ctx.patch.gsPatch.Quests[qid] = q.ToPatch(true)
		ctx.gs.Quests[qid] = q
	}

	return nil
}
