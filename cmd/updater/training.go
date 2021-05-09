package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func completeTrainings(ctx context.Context, tx *redis.Tx, gs *internal.GameState, gsPatch *internal.GameStatePatch, userId string) (pipeFn, error) {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)

	type completion struct {
		queueIdx  int
		popKey    string
		amount    int32
		education internal.Education
	}

	completions := []completion{}

	// Append a completion for every completed training
	now := time.Now().UnixNano()
	for i, t := range gs.TrainingQueue {
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
	gsPatch.TrainingQueue = gs.TrainingQueue
	gsPatch.TrainingQueuePatched = true
	if gsPatch.Population == nil {
		gsPatch.Population = &internal.GameStatePatch_PopulationPatch{}
	}
	for _, c := range completions {
		// Remove completion index from training queue
		gsPatch.TrainingQueue = append(
			gsPatch.TrainingQueue[:c.queueIdx],
			gsPatch.TrainingQueue[c.queueIdx+1:]...,
		)

		// TODO:
		// fix bug that will happen if thieves return at the same time as
		// thieves return
		increasePopulation(gs, gsPatch, c.education, c.amount)
	}

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

func increasePopulation(gs *internal.GameState, gsPatch *internal.GameStatePatch, education internal.Education, amount int32) {
	switch education {
	case internal.Education_CHEF:
		gsPatch.Population.Chefs = &wrapperspb.Int32Value{
			Value: gs.Population.Chefs + amount,
		}
	case internal.Education_SALESMOUSE:
		gsPatch.Population.Salesmice = &wrapperspb.Int32Value{
			Value: gs.Population.Salesmice + amount,
		}
	case internal.Education_GUARD:
		gsPatch.Population.Guards = &wrapperspb.Int32Value{
			Value: gs.Population.Guards + amount,
		}
	case internal.Education_THIEF:
		gsPatch.Population.Thieves = &wrapperspb.Int32Value{
			Value: gs.Population.Thieves + amount,
		}
	}
}
