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
	r           internal.RedisClient
	world       *internal.WorldService
	leaderboard *internal.LeaderboardService
}

func (u *updater) update(ctx context.Context, userId string) {
	log.Debug().Str("userId", userId).Msg("Update")

	gameStateKey := fmt.Sprintf("user:%s:gamestate", userId)

	gs := internal.GameState{}
	var changes changes
	var completedConstructions []*internal.Construction

	// Patch that will be used to notify and update clients
	gsPatch := &internal.GameStatePatch{}
	sendReports := false

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
		completedConstructions = getCompletedConstructions(&gs)

		pipeFns := []pipeFn{}

		pipeFn, err := completeTrainings(ctx, tx, &gs, gsPatch, userId)
		if err != nil {
			return err
		}
		pipeFns = append(pipeFns, pipeFn)

		pipeFn, err = completeTravels(ctx, tx, u.world, userId, &gs, gsPatch, &sendReports)
		if err != nil {
			return err
		}
		pipeFns = append(pipeFns, pipeFn)

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

			for _, pipeFn := range pipeFns {
				if pipeFn != nil {
					if err = pipeFn(pipe); err != nil {
						return err
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

	err = internal.EnsureTimeseries(ctx, u.r, userId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure timeseries")
	}
	err = internal.AddMetricCoins(ctx, u.r, userId, changes.timestamp*1000, int64(changes.coins))
	if err != nil {
		log.Error().Err(err).Msg("Failed to add coins metric")
	}
	err = internal.AddMetricPizzas(ctx, u.r, userId, changes.timestamp*1000, int64(changes.pizzas))
	if err != nil {
		log.Error().Err(err).Msg("Failed to add pizzas metric")
	}

	// Update patch
	if gsPatch.Resources == nil {
		gsPatch.Resources = &internal.GameStatePatch_ResourcesPatch{}
	}
	if gsPatch.Resources.Coins == nil {
		gsPatch.Resources.Coins = &wrapperspb.Int32Value{}
	}
	gsPatch.Resources.Coins.Value = gsPatch.Resources.Coins.Value + changes.coins
	if gsPatch.Resources.Pizzas == nil {
		gsPatch.Resources.Pizzas = &wrapperspb.Int32Value{}
	}
	gsPatch.Resources.Pizzas.Value = gsPatch.Resources.Pizzas.Value + changes.pizzas
	gsPatch.Timestamp = &wrapperspb.Int64Value{Value: changes.timestamp}

	// TODO: move to when constructions are done
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

	if err = u.leaderboard.UpdateUser(ctx, userId, int64(changes.coins)); err != nil {
		log.Error().Err(err).Msg("Failed to update leaderboard")
	}

	if sendReports {
		reports, err := internal.GetReports(ctx, u.r, userId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send updated reports")
			return
		}

		b, err = protojson.Marshal(&internal.ServerMessage{
			Id: xid.New().String(),
			Payload: &internal.ServerMessage_Reports_{
				Reports: &internal.ServerMessage_Reports{
					Reports: reports,
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to send inital reports")
		} else {
			u.r.RPush(ctx, "wsout", &internal.OutgoingMessage{
				ReceiverId: userId,
				Body:       string(b),
			})
		}
	}
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

func extrapolate(gs *internal.GameState) changes {
	// No changes if there are no population
	if gs.Population == nil {
		return changes{}
	}

	now := time.Now()
	rush, offpeak := mtime.GetRush(gs.Timestamp, now.Unix())
	dt := float64(now.Unix() - gs.Timestamp)

	stats := internal.CalculateStats(gs)

	demand := int32(stats.DemandOffpeak*float64(offpeak) +
		stats.DemandRushHour*float64(rush))

	pizzasProduced := int32(stats.PizzasProducedPerSecond * dt)
	pizzasAvailable := gs.Resources.Pizzas + pizzasProduced

	maxSellsByMice := int32(stats.MaxSellsByMicePerSecond * dt)
	pizzasSold := min(demand, min(maxSellsByMice, pizzasAvailable))

	log.Debug().
		Int32("pizzasProduced", pizzasProduced).
		Int32("maxSellsByMice", maxSellsByMice).
		Int32("pizzasSold", pizzasSold).
		Msg("Game state update")

	return changes{
		coins:     gs.Resources.Coins + pizzasSold,
		pizzas:    gs.Resources.Pizzas + pizzasProduced - pizzasSold,
		timestamp: now.Unix(),
	}
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
	world := internal.NewWorldService(rc)
	leaderboard := internal.NewLeaderboardService(rc)
	u := updater{r: rc, world: world, leaderboard: leaderboard}

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
