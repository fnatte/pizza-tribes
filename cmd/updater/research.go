package main

import (
	"time"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/go-redis/redis/v8"
)

func completeResearchs(ctx updateContext, tx *redis.Tx) error {
	completedResearchs := getCompletedResearchs(ctx.gs)

	// Exit early if there are no completed researchs
	if len(completedResearchs) == 0 {
		return nil
	}

	// Update patch
	ctx.patch.gsPatch.ResearchQueue = ctx.gs.ResearchQueue[len(completedResearchs):]
	ctx.patch.gsPatch.ResearchQueuePatched = true

	if ctx.patch.gsPatch.Discoveries == nil {
		ctx.patch.gsPatch.Discoveries = ctx.gs.Discoveries
	}

	for _, r := range completedResearchs {
		if !ctx.patch.gsPatch.HasDiscovery(r.Discovery) {
			ctx.patch.gsPatch.DiscoveriesPatched = true
			ctx.patch.gsPatch.Discoveries =
				append(ctx.patch.gsPatch.Discoveries, r.Discovery)
		}
	}

	// Completion of research can affect the stats
	ctx.patch.sendStats = true

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
