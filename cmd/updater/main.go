package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func getPopulationKey(edu models.Education) (string, error) {
	switch edu {
	case models.Education_CHEF:
		return "chefs", nil
	case models.Education_SALESMOUSE:
		return "salesmice", nil
	case models.Education_GUARD:
		return "guards", nil
	case models.Education_THIEF:
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

type patch struct {
	gsPatch     *models.GameStatePatch
	sendStats   bool
	sendReports bool
}

type updateContext struct {
	context.Context
	userId  string
	gs      *models.GameState
	patch   *patch
	patches map[string]*patch
}

func (ctx *updateContext) initPatch(userId string) {
	if ctx.patches[userId] == nil {
		ctx.patches[userId] = &patch{
			&models.GameStatePatch{
				Resources: &models.GameStatePatch_ResourcesPatch{},
				Population: &models.GameStatePatch_PopulationPatch{},
			},
			false,
			false,
		}
	}
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

	uctx := updateContext{
		ctx,
		userId,
		&models.GameState{},
		&patch{
			&models.GameStatePatch{
				Resources:  &models.GameStatePatch_ResourcesPatch{},
				Population: &models.GameStatePatch_PopulationPatch{},
			},
			false,
			false,
		},
		map[string]*patch{},
	}
	uctx.patches[userId] = uctx.patch

	txf := func(tx *redis.Tx) error {
		// Get current game state
		s, err := internal.RedisJsonGet(tx, ctx, gameStateKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), uctx.gs); err != nil {
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

	if err = u.restorePopulation(ctx, userId); err != nil {
		log.Error().Err(err).Msg("Failed to restore population")
	}

	for patchUserId, p := range uctx.patches {
		if err = u.processPatch(uctx, patchUserId, p); err != nil {
			log.Error().Err(err).
				Bool("wasSender", userId == patchUserId).
				Str("userId", patchUserId).
				Msg("Failed to process patch")
		}
	}
}

func (u *updater) processPatch(ctx updateContext, userId string, p *patch) error {
	if p == nil {
		return nil
	}

	var err error

	// Send game state patch
	if p.gsPatch != nil {
		err = send(ctx, u.r, userId, &models.ServerMessage{
			Id: xid.New().String(),
			Payload: &models.ServerMessage_StateChange{
				StateChange: p.gsPatch,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to send state change")
		}
	}

	// Update leaderboard
	if p.gsPatch != nil && p.gsPatch.Resources.Coins != nil {
		coins := int64(p.gsPatch.Resources.Coins.Value)
		if err = u.leaderboard.UpdateUser(ctx, userId, coins); err != nil {
			log.Error().Err(err).Msg("Failed to update leaderboard")
		}
	}

	// Send reports
	if p.sendReports {
		if err = sendReports(ctx, u.r, userId); err != nil {
			log.Error().Err(err).Msg("Failed to send inital reports")
		}
	}

	// Send stats
	if p.sendStats {
		if err = sendStats(ctx, u.r, userId); err != nil {
			log.Error().Err(err).Msg("Failed to send stats")
		}
	}

	return nil
}

func addTimeseriesDataPoints(ctx updateContext, r internal.RedisClient) error {
	if ctx.patch.gsPatch.Timestamp == nil {
		return nil
	}

	err := internal.EnsureTimeseries(ctx, r, ctx.userId)
	if err != nil {
		return fmt.Errorf("failed to ensure timeseries: %w", err)
	}

	time := ctx.patch.gsPatch.Timestamp.Value * 1000
	coins := ctx.patch.gsPatch.Resources.Coins
	pizzas := ctx.patch.gsPatch.Resources.Pizzas

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

	return send(ctx, r, userId, &models.ServerMessage{
		Id: xid.New().String(),
		Payload: &models.ServerMessage_Reports_{
			Reports: &models.ServerMessage_Reports{
				Reports: reports,
			},
		},
	})
}

func sendStats(ctx context.Context, r internal.RedisClient, userId string) error {
	gs := models.GameState{}
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)
	s, err := internal.RedisJsonGet(r, ctx, gsKey, ".").Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
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
func (u *updater) restorePopulation(ctx context.Context, userId string) error {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)

	gs := &models.GameState{}
	var addedPop int32 = 0

	txf := func(tx *redis.Tx) error {
		// Get current game state
		s, err := internal.RedisJsonGet(tx, ctx, gsKey, ".").Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if err = protojson.Unmarshal([]byte(s), gs); err != nil {
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
		return send(ctx, u.r, userId, &models.ServerMessage{
			Id: xid.New().String(),
			Payload: &models.ServerMessage_StateChange{
				StateChange: &models.GameStatePatch{
					Population: &models.GameStatePatch_PopulationPatch{
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
func send(ctx context.Context, r redis.Cmdable, userId string, msg *models.ServerMessage) error {
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

		u.update(ctx, userId)
		time.Sleep(1 * time.Millisecond)
	}

}
