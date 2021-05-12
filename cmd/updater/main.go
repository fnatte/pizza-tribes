package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
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

type updateContext struct {
	context.Context
	userId      string
	gs          *internal.GameState
	gsPatch     *internal.GameStatePatch
	sendReports *bool
	sendStats   *bool
}

// Update the game state for the specified user
func (u *updater) update(ctx context.Context, userId string) {
	log.Debug().Str("userId", userId).Msg("Update")

	/*
	 * There's a lot of stuff happening in this function, but this is
	 * also the heart of the game. In short it does:
	 *   - calculate changes
	 *   - apply changes in a Redis pipeline
	 *   - send changes to the client (so that UI can be updated in the web app)
	 *   - update the leaderboard
	 *   - update the timeseries data
	 *   - set the next update time
	 */

	gameStateKey := fmt.Sprintf("user:%s:gamestate", userId)

	sendReportsBoolRef := false
	sendStats2 := false

	uctx := updateContext{
		ctx,
		userId,
		&internal.GameState{},
		&internal.GameStatePatch{
			Resources:  &internal.GameStatePatch_ResourcesPatch{},
			Population: &internal.GameStatePatch_PopulationPatch{},
		},
		&sendReportsBoolRef,
		&sendStats2}

	txf := func(tx *redis.Tx) error {
		// Get current game state
		b, err := internal.RedisJsonGet(tx, ctx, gameStateKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		err = uctx.gs.LoadProtoJson([]byte(b))
		if err != nil {
			return err
		}

		// Prepare changes and setup pipeline funcs.
		// Each of the following functions returns a pipeline function
		// that is executed serially in a single Redis pipeline below.
		pipeFns := make([]pipeFn, 4)
		if pipeFns[0], err = extrapolate(uctx, tx); err != nil {
			return err
		}
		if pipeFns[1], err = completedConstructions(uctx, tx); err != nil {
			return err
		}
		if pipeFns[2], err = completeTrainings(uctx, tx); err != nil {
			return err
		}
		if pipeFns[3], err = completeTravels(uctx, tx, u.world); err != nil {
			return err
		}

		// Write changes to Redis (execute pipelines)
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
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

	_, err = internal.SetNextUpdate(u.r, ctx, userId, uctx.gs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set next update")
	}

	if err = addTimeseriesDataPoints(uctx, u.r); err != nil {
		log.Error().Err(err).Msg("Failed to add timeseries data points")
	}

	// Send patch
	err = send(ctx, u.r, userId, &internal.ServerMessage{
		Id: xid.New().String(),
		Payload: &internal.ServerMessage_StateChange{
			StateChange: uctx.gsPatch,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send state change")
	}

	// Update leaderboard
	if uctx.gsPatch.Resources.Coins != nil {
		coins := int64(uctx.gsPatch.Resources.Coins.Value)
		if err = u.leaderboard.UpdateUser(ctx, userId, coins); err != nil {
			log.Error().Err(err).Msg("Failed to update leaderboard")
		}
	}

	if *uctx.sendReports {
		if err = sendReports(ctx, u.r, userId); err != nil {
			log.Error().Err(err).Msg("Failed to send inital reports")
		}
	}

	if *uctx.sendStats {
		if err = sendStats(ctx, u.r, userId); err != nil {
			log.Error().Err(err).Msg("Failed to send stats")
		}
	}

	if err = u.restorePopulation(ctx, userId); err != nil {
		log.Error().Err(err).Msg("Failed to restore population")
	}
}

func addTimeseriesDataPoints(ctx updateContext, r internal.RedisClient) error {
	if ctx.gsPatch.Timestamp == nil {
		return nil
	}

	err := internal.EnsureTimeseries(ctx, r, ctx.userId)
	if err != nil {
		return fmt.Errorf("failed to ensure timeseries: %w", err)
	}

	time := ctx.gsPatch.Timestamp.Value * 1000
	coins := ctx.gsPatch.Resources.Coins
	pizzas := ctx.gsPatch.Resources.Pizzas

	if coins != nil {
		err = internal.AddMetricCoins(ctx, r, ctx.userId, time, int64(coins.Value))
		if err != nil {
			return fmt.Errorf("failed to add coins metric: %w", err)
		}
	}

	if pizzas != nil {
		err = internal.AddMetricPizzas(ctx, r, ctx.userId, time, int64(pizzas.Value))
		if err != nil {
			return fmt.Errorf("failed to add pizzas metric: %w", err)
		}
	}

	return nil
}

func sendReports(ctx context.Context, r internal.RedisClient, userId string) error {
	reports, err := internal.GetReports(ctx, r, userId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send updated reports")
		return fmt.Errorf("failed to get reports: %w", err)
	}

	return send(ctx, r, userId, &internal.ServerMessage{
		Id: xid.New().String(),
		Payload: &internal.ServerMessage_Reports_{
			Reports: &internal.ServerMessage_Reports{
				Reports: reports,
			},
		},
	})
}

func sendStats(ctx context.Context, r internal.RedisClient, userId string) error {
	gs := internal.GameState{}
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)
	b, err := internal.RedisJsonGet(r, ctx, gsKey, ".").Result()
	if err != nil && err != redis.Nil {
		return err
	}
	err = gs.LoadProtoJson([]byte(b))
	if err != nil {
		return err
	}

	msg := internal.CalculateStats(&gs).ToServerMessage()

	return send(ctx, r, userId, msg)
}

// Restore population will add any missing population to the town.
// If the user lost thieves in a heist, those mice will be replaced
// by uneducated mice in the town.
//
// This process in currently not implemented as a part of the normal
// game state update pipeline because it does not handle modification
// of the same variables twice with ease. Once that is fixed, this
// func should be part of the update pipeline.
func (u *updater) restorePopulation(ctx context.Context, userId string) (error) {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)

	gs := &internal.GameState{}
	var addedPop int32 = 0

	txf := func(tx *redis.Tx) error {
		// Get current game state
		b, err := internal.RedisJsonGet(tx, ctx, gsKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		err = gs.LoadProtoJson([]byte(b))
		if err != nil {
			return err
		}

		pop := internal.CountAllPopulation(gs)
		maxPop := internal.CountMaxPopulation(gs)

		if pop < maxPop {
			addedPop = maxPop - pop
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				return internal.RedisJsonNumIncrBy(
					u.r, ctx, gsKey,
					".population.uneducated",
					int64(addedPop)).Err()
			})
		}

		return err
	}

	if err := u.r.Watch(ctx, txf, gsKey); err != nil {
		return fmt.Errorf("failed to restore population: %w", err)
	}

	if addedPop != 0 {
		return send(ctx, u.r, userId, &internal.ServerMessage{
			Id: xid.New().String(),
			Payload: &internal.ServerMessage_StateChange{
				StateChange: &internal.GameStatePatch{
					Population: &internal.GameStatePatch_PopulationPatch{
						Uneducated: &wrapperspb.Int32Value{
							Value: gs.Population.Uneducated + addedPop,
						},
					},
				},
			},
		})
	}

	return nil
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

// Send a message to the specified userId
func send(ctx context.Context, r redis.Cmdable, userId string, msg *internal.ServerMessage) error {
	b, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	return r.RPush(ctx, "wsout", &internal.OutgoingMessage{
		ReceiverId: userId,
		Body:       string(b),
	}).Err()
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
