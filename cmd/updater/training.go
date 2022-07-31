package main

import (
	"errors"
	"time"

	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
)

func completeTrainings(userId string, gs *models.GameState, tx *gamestate.GameTx) error {
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
	for i, t := range gs.TrainingQueue {
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
	u := tx.Users[userId]


	q := gs.TrainingQueue
	for _, c := range completions {
		// Remove completion by index from training queue
		q = append(
			q[:c.queueIdx],
			q[c.queueIdx+1:]...,
		)

		var mouseId string
		for id, m := range u.Gs.Mice {
			if m.IsBeingEducated {
				mouseId = id
				break
			}
		}
		if mouseId == "" {
			return errors.New("could not find mouse being educated")
		}
		u.SetMouseEducation(mouseId, c.education)
	}

	u.SetTrainingQueue(q)

	// Since we have changed the population we should send a new stats message
	u.StatsInvalidated = true

	return nil
}

