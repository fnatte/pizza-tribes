package main

import (
	"fmt"
	"math"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func completedConstructions(ctx updateContext, tx *redis.Tx) (pipeFn, error) {
	gsKey := fmt.Sprintf("user:%s:gamestate", ctx.userId)

	completedConstructions := getCompletedConstructions(ctx.gs)

	// Exit early if there are no completed constructions
	if len(completedConstructions) == 0 {
		return nil, nil
	}

	// Update patch
	ctx.gsPatch.ConstructionQueue = ctx.gs.ConstructionQueue[len(completedConstructions):]
	ctx.gsPatch.ConstructionQueuePatched = true
	ctx.gsPatch.Lots = map[string]*internal.GameStatePatch_LotPatch{}

	for _, constr := range completedConstructions {
		ctx.gsPatch.Lots[constr.LotId] = &internal.GameStatePatch_LotPatch{
			Building: constr.Building,
		}

		if constr.Building == internal.Building_HOUSE {
			if ctx.gsPatch.Population == nil {
				ctx.gsPatch.Population = &internal.GameStatePatch_PopulationPatch{}
			}
			ctx.gsPatch.Population.Uneducated = &wrapperspb.Int32Value{
				Value: ctx.gs.Population.Uneducated + internal.MicePerHouse,
			}

		}
	}

	// Completion of buildings can affect the stats because we increase the
	// number of employables. E.g. if the player had 10 chefs but only 5 of them
	// were employed.
	*ctx.sendStats = true

	return func(pipe redis.Pipeliner) error {
		if len(completedConstructions) > 0 {
			// Remove completed constructions from queue
			err := internal.RedisJsonArrTrim(
				pipe, ctx, gsKey,
				".constructionQueue",
				len(completedConstructions),
				math.MaxInt32,
			).Err()
			if err != nil {
				return err
			}

			// Complete constructions
			for _, constr := range completedConstructions {
				err = internal.RedisJsonSet(
					pipe, ctx, gsKey,
					fmt.Sprintf(".lots[\"%s\"]", constr.LotId),
					fmt.Sprintf("{ \"building\": %d }", constr.Building)).Err()
				if err != nil {
					return fmt.Errorf("failed to set building on lot: %w", err)
				}

				// Increase population (uneducated) if a house was completed
				if constr.Building == internal.Building_HOUSE {
					_, err = internal.RedisJsonNumIncrBy(
						pipe, ctx, gsKey,
						".population.uneducated",
						internal.MicePerHouse).Result()
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}, nil
}

func getCompletedConstructions(gs *internal.GameState) (res []*internal.Construction) {
	now := time.Now().UnixNano()

	for _, t := range gs.ConstructionQueue {
		if t.CompleteAt > now {
			break
		}

		res = append(res, t)
	}

	return res
}
