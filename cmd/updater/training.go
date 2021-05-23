package main

import (
	"time"

	"github.com/fnatte/pizza-tribes/internal/models"
)

func completeTrainings(ctx updateContext) error {
	// Setup a internal completion struct to hold completed trainings.
	// By using the internal data structure it will be easier to apply
	// the changes in the Redis pipeline. Also, we avoid doing stuff
	// that can return errors in the pipeline.
	type completion struct {
		queueIdx  int
		popKey    string
		amount    int32
		education models.Education
	}
	completions := []completion{}

	// Append a completion for every completed training
	now := time.Now().UnixNano()
	for i, t := range ctx.gs.TrainingQueue {
		if t.CompleteAt > now {
			continue
		}

		popKey, err := getPopulationKey(t.Education)
		if err != nil {
			return err
		}

		// Prepend to completions!
		// This ensures that when ARRPOP runs we pop the highest
		// indexes first - so that we don't shift the indexes that
		// we want to remove.
		completions = append([]completion{{
			queueIdx:  i,
			popKey:    popKey,
			amount:    t.Amount,
			education: t.Education,
		}}, completions...)
	}

	// Exit early if there are no completed trainings
	if len(completions) == 0 {
		return nil
	}

	// Update patch
	ctx.patch.gsPatch.TrainingQueue = ctx.gs.TrainingQueue
	ctx.patch.gsPatch.TrainingQueuePatched = true
	for _, c := range completions {
		// Remove completion index from training queue
		ctx.patch.gsPatch.TrainingQueue = append(
			ctx.patch.gsPatch.TrainingQueue[:c.queueIdx],
			ctx.patch.gsPatch.TrainingQueue[c.queueIdx+1:]...,
		)

		switch c.education {
		case models.Education_CHEF:
			ctx.IncrChefs(c.amount)
		case models.Education_SALESMOUSE:
			ctx.IncrSalesmice(c.amount)
		case models.Education_GUARD:
			ctx.IncrGuards(c.amount)
		case models.Education_THIEF:
			ctx.IncrThieves(c.amount)
		case models.Education_PUBLICIST:
			ctx.IncrPublicists(c.amount)
		}
	}

	// Since we have changed the population we should send a new stats message
	ctx.patch.sendStats = true

	return nil
}

