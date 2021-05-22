package main

import (
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/go-redis/redis/v8"
)

func completedConstructions(ctx updateContext, tx *redis.Tx) (error) {
	completedConstructions := getCompletedConstructions(ctx.gs)

	// Exit early if there are no completed constructions
	if len(completedConstructions) == 0 {
		return nil
	}

	// Update patch
	ctx.patch.gsPatch.ConstructionQueue = ctx.gs.ConstructionQueue[len(completedConstructions):]
	ctx.patch.gsPatch.ConstructionQueuePatched = true
	ctx.patch.gsPatch.Lots = map[string]*models.GameStatePatch_LotPatch{}

	for _, constr := range completedConstructions {
		if constr.Razing {
			ctx.patch.gsPatch.Lots[constr.LotId] = &models.GameStatePatch_LotPatch{
				Razed: true,
			}
		} else {
			ctx.patch.gsPatch.Lots[constr.LotId] = &models.GameStatePatch_LotPatch{
				Building: constr.Building,
				Level:    constr.Level,
			}
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
			if !constr.Razing {
				if constr.Level > 0 {
					prevLevelInfo := buildInfo.LevelInfos[constr.Level-1]
					ctx.IncrUneducated(levelInfo.Residence.Beds - prevLevelInfo.Residence.Beds)
				} else {
					ctx.IncrUneducated(levelInfo.Residence.Beds)
				}
			} else {
				rest := levelInfo.Residence.Beds

				if rest > ctx.gs.Population.Uneducated {
					ctx.IncrUneducated(-ctx.gs.Population.Uneducated)
					rest = rest - ctx.gs.Population.Uneducated
				} else {
					ctx.IncrUneducated(-rest)
					rest = 0
				}

				popKey := 0
				loopCount := 0
				for rest > 0 && loopCount < 1000 {
					switch popKey {
					case 0:
						if ctx.gs.Population.Chefs > 0 {
							ctx.IncrChefs(-1)
							rest = rest - 1
						}
					case 1:
						if ctx.gs.Population.Salesmice > 0 {
							ctx.IncrSalesmice(-1)
							rest = rest - 1
						}
					case 2:
						if ctx.gs.Population.Guards > 0 {
							ctx.IncrGuards(-1)
							rest = rest - 1
						}
					case 3:
						if ctx.gs.Population.Thieves > 0 {
							ctx.IncrThieves(-1)
							rest = rest - 1
						}
					}
					popKey = (popKey + 1) % 4
					loopCount++
				}
			}
		}
	}

	// Completion of buildings can affect the stats because we increase the
	// number of employables. E.g. if the player had 10 chefs but only 5 of them
	// were employed.
	ctx.patch.sendStats = true

	return nil
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
