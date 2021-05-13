package main

import (
	"fmt"
	"math"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
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
	ctx.gsPatch.Lots = map[string]*models.GameStatePatch_LotPatch{}

	var incrUneducated int32 = 0

	for _, constr := range completedConstructions {
		ctx.gsPatch.Lots[constr.LotId] = &models.GameStatePatch_LotPatch{
			Building: constr.Building,
			Level:    constr.Level,
		}

		buildInfo := internal.FullGameData.Buildings[int32(constr.Building)]
		if buildInfo == nil {
			continue
		}

		levelInfo := buildInfo.LevelInfos[constr.Level]
		if levelInfo == nil {
			continue
		}

		if levelInfo.Residence != nil {
			if constr.Level > 0 {
				prevLevelInfo := buildInfo.LevelInfos[constr.Level - 1]
				incrUneducated = levelInfo.Residence.Beds - prevLevelInfo.Residence.Beds
			} else {
				incrUneducated = levelInfo.Residence.Beds
			}

			if ctx.gsPatch.Population == nil {
				ctx.gsPatch.Population = &models.GameStatePatch_PopulationPatch{}
			}
			ctx.gsPatch.Population.Uneducated = &wrapperspb.Int32Value{
				Value: ctx.gs.Population.Uneducated + incrUneducated,
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

				if constr.Level == 0 {
					// Place building (level 0)
					err = internal.RedisJsonSet(
						pipe, ctx, gsKey,
						fmt.Sprintf(".lots[\"%s\"]", constr.LotId),
						fmt.Sprintf("{ \"building\": %d }", constr.Building)).Err()
					if err != nil {
						return fmt.Errorf("failed to set building on lot: %w", err)
					}
				} else {
					// Upgrade building (level >= 1)
					err = internal.RedisJsonSet(
						pipe, ctx, gsKey,
						fmt.Sprintf(".lots[\"%s\"].level", constr.LotId),
						constr.Level).Err()
					if err != nil {
						return fmt.Errorf("failed to set building on lot: %w", err)
					}
				}

				if incrUneducated > 0 {
					_, err = internal.RedisJsonNumIncrBy(
						pipe, ctx, gsKey,
						".population.uneducated",
						int64(incrUneducated)).Result()
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}, nil
}

func getCompletedConstructions(gs *models.GameState) (res []*models.Construction) {
	now := time.Now().UnixNano()

	for _, t := range gs.ConstructionQueue {
		if t.CompleteAt > now {
			break
		}

		res = append(res, t)
	}

	return res
}
