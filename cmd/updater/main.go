package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/mtime"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

					// Increase population (uneducated) if a house was completed
					if constr.Building == internal.Building_HOUSE {
						_, err = internal.RedisJsonNumIncrBy(
							pipe, ctx, gameStateKey,
							".population.uneducated",
							internal.MicePerHouse).Result()
						if err != nil {
							return err
						}
					}
				}
			}

			return nil
		})
		return err
	}

	err := u.r.Watch(ctx, txf, gameStateKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update")
		return
	}

	_, err = internal.UpdateTimestamp(u.r, ctx, userId, &gs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update timestamp")
	}

	// Notify
	gsPatch := &internal.GameStatePatch{
		Resources: &internal.GameStatePatch_ResourcesPatch{
			Coins:  &wrapperspb.Int32Value{Value: changes.coins},
			Pizzas: &wrapperspb.Int32Value{Value: changes.pizzas},
		},
		Timestamp: &wrapperspb.Int64Value{Value: changes.timestamp},
	}

	if len(completedTrainings) > 0 {
		gsPatch.TrainingQueue = gs.TrainingQueue[len(completedTrainings):]
		gsPatch.TrainingQueuePatched = true
		gsPatch.Population = &internal.GameStatePatch_PopulationPatch{}
		for _, t := range completedTrainings {
			switch t.Education {
			case internal.Education_CHEF:
				gsPatch.Population.Chefs = &wrapperspb.Int32Value{
					Value: gs.Population.Chefs + t.Amount,
				}
			case internal.Education_SALESMOUSE:
				gsPatch.Population.Salesmice = &wrapperspb.Int32Value{
					Value: gs.Population.Salesmice + t.Amount,
				}
			case internal.Education_GUARD:
				gsPatch.Population.Guards = &wrapperspb.Int32Value{
					Value: gs.Population.Guards + t.Amount,
				}
			case internal.Education_THIEF:
				gsPatch.Population.Thieves = &wrapperspb.Int32Value{
					Value: gs.Population.Thieves + t.Amount,
				}
			}
		}
	}

	if len(completedConstructions) > 0 {
		gsPatch.ConstructionQueue = gs.ConstructionQueue[len(completedConstructions):]
		gsPatch.ConstructionQueuePatched = true
		gsPatch.Lots = map[string]*internal.GameStatePatch_LotPatch{}
		for _, constr := range completedConstructions {
			gsPatch.Lots[constr.LotId] = &internal.GameStatePatch_LotPatch{
				Building: constr.Building,
			}

			if constr.Building == internal.Building_HOUSE {
				if gsPatch.Population == nil {
					gsPatch.Population = &internal.GameStatePatch_PopulationPatch{}
				}
				gsPatch.Population.Uneducated = &wrapperspb.Int32Value{
					Value: gs.Population.Uneducated + internal.MicePerHouse,
				}
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

func (u *updater) next(ctx context.Context) (string, error) {
	packed, err := u.r.ZRangeWithScores(ctx, "user_updates", 0, 0).Result()
	if err != nil {
		return "", err
	}

	if len(packed) == 0 {
		return "", nil
	}

	timestamp := int64(packed[0].Score)
	if timestamp > time.Now().UnixNano() {
		return "", nil
	}

	userId, ok := packed[0].Member.(string)
	if !ok {
		return "", errors.New("Failed to parse user update member")
	}

	removed, err := u.r.ZRem(ctx, "user_updates", userId).Result()
	if err != nil {
		return "", err
	}

	if removed != 1 {
		return "", nil
	}

	return userId, nil
}

type changes struct {
	timestamp int64
	coins     int32
	pizzas    int32
}

func min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
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

func countPopulation(gs *internal.GameState) int32 {
	if gs.Population == nil {
		return 0
	}

	return (gs.Population.Uneducated +
		gs.Population.Chefs +
		gs.Population.Salesmice +
		gs.Population.Guards +
		gs.Population.Thieves)
}

func extrapolate(gs *internal.GameState) changes {
	// No changes if there are no population
	if gs.Population == nil {
		return changes{}
	}

	now := time.Now()
	rush, offpeak := mtime.GetRush(gs.Timestamp, now.Unix())
	dt := float64(now.Unix() - gs.Timestamp)

	buildingCount := internal.CountBuildings(gs)
	maxEmployed := countMaxEmployed(buildingCount)

	employedChefs := min(gs.Population.Chefs, maxEmployed[int32(internal.Building_KITCHEN)])
	pizzasProduced := int32(float64(employedChefs) * 0.2 * dt)
	pizzasAvailable := gs.Resources.Pizzas + pizzasProduced

	popularity := 1.0 * float64(countPopulation(gs))
	demand := int32((float64(rush)*0.75 + float64(offpeak)*0.2) * popularity)
	employedSalesmice := min(gs.Population.Salesmice, maxEmployed[int32(internal.Building_SHOP)])
	maxSellsByMice := int32(float64(employedSalesmice) * 0.5 * dt)
	pizzasSold := min(demand, min(maxSellsByMice, pizzasAvailable))

	log.Info().
		Int32("chefs", gs.Population.Chefs).
		Int32("employedChefs", employedChefs).
		Int32("salesmice", gs.Population.Salesmice).
		Int32("employedSalesmice", employedSalesmice).
		Float64("dt", dt).
		Int64("rush", rush).
		Int64("offpeak", offpeak).
		Int32("demand", demand).
		Int32("maxSellsByMice", maxSellsByMice).
		Int32("pizzasProduced", pizzasProduced).
		Int32("pizzasSold", pizzasSold).
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

func envOrDefault(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func main() {
	log.Info().Msg("Starting updater")

	rdb := redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	rc := internal.NewRedisClient(rdb)
	u := updater{r: rc}

	ctx := context.Background()

	for {
		userId, err := u.next(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get next")
			time.Sleep(1 * time.Second)
			continue
		}
		if userId == "" {
			// TODO: could be smarter about how long to sleep here
			time.Sleep(10 * time.Millisecond)
			continue
		}

		log.Info().Msg("Updating")

		u.update(ctx, userId)
		time.Sleep(1 * time.Millisecond)
	}

}
