package main

import (
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func completeTrainings(ctx updateContext, tx *redis.Tx) (pipeFn, error) {
	gsKey := fmt.Sprintf("user:%s:gamestate", ctx.userId)

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
			return nil, err
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
		return nil, nil
	}

	// Update patch
	ctx.gsPatch.TrainingQueue = ctx.gs.TrainingQueue
	ctx.gsPatch.TrainingQueuePatched = true
	if ctx.gsPatch.Population == nil {
		ctx.gsPatch.Population = &models.GameStatePatch_PopulationPatch{}
	}
	for _, c := range completions {
		// Remove completion index from training queue
		ctx.gsPatch.TrainingQueue = append(
			ctx.gsPatch.TrainingQueue[:c.queueIdx],
			ctx.gsPatch.TrainingQueue[c.queueIdx+1:]...,
		)

		// TODO:
		// fix bug that will happen if thieves return at the same time as
		// thieves return
		increasePopulation(ctx.gs, ctx.gsPatch, c.education, c.amount)
	}

	// Since we have changed the population we should send a new stats message
	*ctx.sendStats = true

	return func(pipe redis.Pipeliner) error {
		for _, c := range completions {
			err := internal.RedisJsonArrPop(
				pipe, ctx, gsKey,
				".trainingQueue", c.queueIdx).Err()
			if err != nil {
				return fmt.Errorf("failed to remove from training queue: %w", err)
			}

			_, err = internal.RedisJsonNumIncrBy(
				pipe, ctx, gsKey,
				fmt.Sprintf(".population.%s", c.popKey),
				int64(c.amount)).Result()
			if err != nil {
				return fmt.Errorf("failed to increase population: %w", err)
			}
		}
		return nil
	}, nil
}

func increasePopulation(gs *models.GameState, gsPatch *models.GameStatePatch, education models.Education, amount int32) {
	switch education {
	case models.Education_CHEF:
		gsPatch.Population.Chefs = &wrapperspb.Int32Value{
			Value: gs.Population.Chefs + amount,
		}
	case models.Education_SALESMOUSE:
		gsPatch.Population.Salesmice = &wrapperspb.Int32Value{
			Value: gs.Population.Salesmice + amount,
		}
	case models.Education_GUARD:
		gsPatch.Population.Guards = &wrapperspb.Int32Value{
			Value: gs.Population.Guards + amount,
		}
	case models.Education_THIEF:
		gsPatch.Population.Thieves = &wrapperspb.Int32Value{
			Value: gs.Population.Thieves + amount,
		}
	}
}
