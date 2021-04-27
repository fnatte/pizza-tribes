package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/fnatte/pizza-mouse/internal"
	"github.com/fnatte/pizza-mouse/internal/mtime"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func getPopulationKey(edu internal.Education) (string, error) {
	switch edu {
	case internal.Education_CHEF:
		return "chefs", nil
	case internal.Education_SALESMOUSE:
		return "salesmice", nil
	case internal.Education_GUARD:
		return "guards", nil
	case internal.Education_THIEF:
		return "thieves", nil
	default:
		return "", fmt.Errorf("Invalid education: %s", edu)
	}
}

type updater struct {
	r internal.RedisClient
}

func (u *updater) update(ctx context.Context, userId string) {
	log.Info().Str("userId", userId).Msg("Update")

	gameStateKey := fmt.Sprintf("user:%s:gamestate", userId)

	gs := internal.GameState{}
	var changes changes
	var completedTrainings []*internal.Training
	var completedConstructions []*internal.Construction

	txf := func(tx *redis.Tx) error {
		// Get current game state
		b, err := internal.RedisJsonGet(tx, ctx, gameStateKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		err = gs.LoadProtoJson([]byte(b))
		if err != nil {
			return err
		}

		// Calculate changes
		changes = extrapolate(&gs)
		completedTrainings = getCompletedTrainings(&gs)
		completedConstructions = getCompletedConstructions(&gs)

		// Write changes to Redis
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// Write timestamp
			err = internal.RedisJsonSet(
				pipe, ctx, gameStateKey,
				".timestamp", changes.timestamp).Err()
			if err != nil {
				return err
			}

			// Write coins
			err = internal.RedisJsonSet(
				pipe, ctx, gameStateKey,
				".resources.coins", changes.coins).Err()
			if err != nil {
				return err
			}

			// Write pizzas
			err = internal.RedisJsonSet(
				pipe, ctx, gameStateKey,
				".resources.pizzas", changes.pizzas).Err()
			if err != nil {
				return err
			}

			if len(completedTrainings) > 0 {
				// Remove completed trainings from queue
				err = internal.RedisJsonArrTrim(
					pipe, ctx, gameStateKey,
					".trainingQueue",
					len(completedTrainings),
					math.MaxInt32,
				).Err()
				if err != nil {
					return err
				}

				// Complete trainings by increasing the corresponding population
				for _, t := range completedTrainings {
					populationKey, err := getPopulationKey(t.Education)
					if err != nil {
						return err
					}

					_, err = internal.RedisJsonNumIncrBy(
						pipe, ctx, gameStateKey,
						fmt.Sprintf(".population.%s", populationKey),
						int64(t.Amount)).Result()
					if err != nil {
						return err
					}
				}
			}

			if len(completedConstructions) > 0 {
				// Remove completed constructions from queue
				err = internal.RedisJsonArrTrim(
					pipe, ctx, gameStateKey,
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
						pipe, ctx, gameStateKey,
						fmt.Sprintf(".lots[\"%s\"]", constr.LotId),
						fmt.Sprintf("{ \"building\": %d }", constr.Building)).Err()
					if err != nil {
						log.Error().Err(err).Msg("Failed to handle construct building message")
					}
				}
			}

			return nil
		})
		return err
	}

	err := u.r.Watch(ctx, txf, gameStateKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to train")
		return
	}

	// Notify
	gsPatch := &internal.GameStatePatch{
		Resources: &internal.GameStatePatch_ResourcesPatch{
			Coins:  internal.NewInt64(changes.coins),
			Pizzas: internal.NewInt64(changes.pizzas),
		},
		Timestamp: internal.NewInt64(changes.timestamp),
	}

	if len(completedTrainings) > 0 {
		gsPatch.TrainingQueue = gs.TrainingQueue[len(completedTrainings):]
		gsPatch.TrainingQueuePatched = true
		gsPatch.Population = &internal.GameStatePatch_PopulationPatch{}
		for _, t := range completedTrainings {
			switch t.Education {
			case internal.Education_CHEF:
				gsPatch.Population.Chefs = internal.NewInt64(
					gs.Population.Chefs + int64(t.Amount),
				)
			case internal.Education_SALESMOUSE:
				gsPatch.Population.Salesmice = internal.NewInt64(
					gs.Population.Salesmice + int64(t.Amount),
				)
			case internal.Education_GUARD:
				gsPatch.Population.Guards = internal.NewInt64(
					gs.Population.Guards + int64(t.Amount),
				)
			case internal.Education_THIEF:
				gsPatch.Population.Thieves = internal.NewInt64(
					gs.Population.Thieves + int64(t.Amount),
				)
			}
		}
	}

	if len(completedConstructions) > 0 {
		gsPatch.ConstructionQueue = gs.ConstructionQueue[len(completedConstructions):]
		gsPatch.ConstructionQueuePatched = true
		gsPatch.Lots = map[string]*internal.GameStatePatch_LotPatch{}
		for _, constr := range completedConstructions {
			gsPatch.Lots[constr.LotId] = &internal.GameStatePatch_LotPatch{
				Building: &constr.Building,
			}
		}
	}

	msg := internal.ServerMessage{
		Id: xid.New().String(),
		Payload: &internal.ServerMessage_StateChange{
			StateChange: gsPatch,
		},
	}

	b, err := protojson.Marshal(&msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send full state update")
		return
	}

	u.r.RPush(ctx, "wsout", &internal.OutgoingMessage{
		ReceiverId: userId,
		Body:       string(b),
	})
}

func (u *updater) next(ctx context.Context) string {
	return "c20okiv8ecspveu0v5hg"
}

type changes struct {
	timestamp int64
	coins     int64
	pizzas    int64
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func countBuildings(gs *internal.GameState) (counts map[int32]int32) {
	counts = map[int32]int32{}
	for _, lot := range gs.Lots {
		counts[int32(lot.Building)] = counts[int32(lot.Building)] + 1
	}
	return counts
}

func countMaxEmployed(buildingCount map[int32]int32) (counts map[int32]int32) {
	counts = map[int32]int32{}
	buildings := internal.FullGameData.Buildings
	for k := range internal.Building_name {
		employer := buildings[k].Employer
		if employer != nil {
			maxWorkforce := employer.MaxWorkforce
			counts[k] = buildingCount[k] * maxWorkforce
		}
	}
	return counts
}

func extrapolate(gs *internal.GameState) changes {
	now := time.Now()
	rush, offpeak := mtime.GetRush(gs.Timestamp, now.Unix())
	dt := float64(now.Unix() - gs.Timestamp)

	buildingCount := countBuildings(gs)
	maxEmployed := countMaxEmployed(buildingCount)

	employedChefs := min(gs.Population.Chefs, int64(maxEmployed[int32(internal.Building_KITCHEN)]))
	pizzasProduced := int64(float64(employedChefs) * 0.2 * dt)
	pizzasAvailable := gs.Resources.Pizzas + pizzasProduced

	demand := int64(float64(rush)*0.75 + float64(offpeak)*0.2)
	employedSalesmice := min(gs.Population.Salesmice, int64(maxEmployed[int32(internal.Building_SHOP)]))
	maxSellsByMice := int64(float64(employedSalesmice) * 0.5 * dt)
	pizzasSold := min(demand, min(maxSellsByMice, pizzasAvailable))

	log.Info().
		Int64("chefs", gs.Population.Chefs).
		Int64("employedChefs", employedChefs).
		Int64("salesmice", gs.Population.Salesmice).
		Int64("employedSalesmice", employedSalesmice).
		Float64("dt", dt).
		Int64("rush", rush).
		Int64("offpeak", offpeak).
		Int64("demand", demand).
		Int64("maxSellsByMice", maxSellsByMice).
		Int64("pizzasProduced", pizzasProduced).
		Int64("pizzasSold", pizzasSold).
		Msg("Game state update")

	return changes{
		coins:     gs.Resources.Coins + pizzasSold,
		pizzas:    gs.Resources.Pizzas + pizzasProduced - pizzasSold,
		timestamp: now.Unix(),
	}
}

func getCompletedTrainings(gs *internal.GameState) (res []*internal.Training) {
	now := time.Now().UnixNano()

	for _, t := range gs.TrainingQueue {
		if t.CompleteAt > now {
			break
		}

		res = append(res, t)
	}

	return res
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

func main() {
	log.Info().Msg("Starting updater")

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rc := internal.NewRedisClient(rdb)
	u := updater{r: rc}

	ctx := context.Background()

	for {
		userId := u.next(ctx)
		u.update(ctx, userId)
		time.Sleep(5 * time.Second)
	}

}
