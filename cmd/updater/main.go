package main

import (
	"context"
	"errors"
	"fmt"
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

type updateContext struct {
	context.Context
	userId      string
	gs          *internal.GameState
	gsPatch     *internal.GameStatePatch
	sendReports *bool
	sendStats   *bool
}

func (u *updater) update(ctx context.Context, userId string) {
	log.Debug().Str("userId", userId).Msg("Update")

	gameStateKey := fmt.Sprintf("user:%s:gamestate", userId)

	sendReports := false
	sendStats2 := false

	uctx := updateContext{
		ctx,
		userId,
		&internal.GameState{},
		&internal.GameStatePatch{},
		&sendReports,
		&sendStats2}
	var changes changes

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

		// Calculate changes
		changes = extrapolate(uctx.gs)

		// Prepare changes and setup pipeline funcs
		pipeFn1, err := completedConstructions(uctx, tx)
		if err != nil {
			return err
		}

		pipeFn2, err := completeTrainings(uctx, tx)
		if err != nil {
			return err
		}

		pipeFn3, err := completeTravels(uctx, tx, u.world)
		if err != nil {
			return err
		}

		pipeFns := []pipeFn{pipeFn1, pipeFn2, pipeFn3}

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

	_, err = internal.UpdateTimestamp(u.r, ctx, userId, uctx.gs)
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
	if uctx.gsPatch.Resources == nil {
		uctx.gsPatch.Resources = &internal.GameStatePatch_ResourcesPatch{}
	}
	if uctx.gsPatch.Resources.Coins == nil {
		uctx.gsPatch.Resources.Coins = &wrapperspb.Int32Value{}
	}
	uctx.gsPatch.Resources.Coins.Value = uctx.gsPatch.Resources.Coins.Value + changes.coins
	if uctx.gsPatch.Resources.Pizzas == nil {
		uctx.gsPatch.Resources.Pizzas = &wrapperspb.Int32Value{}
	}
	uctx.gsPatch.Resources.Pizzas.Value = uctx.gsPatch.Resources.Pizzas.Value + changes.pizzas
	uctx.gsPatch.Timestamp = &wrapperspb.Int64Value{Value: changes.timestamp}

	err = send(ctx, u.r, userId, &internal.ServerMessage{
		Id: xid.New().String(),
		Payload: &internal.ServerMessage_StateChange{
			StateChange: uctx.gsPatch,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send state change")
	}

	if err = u.leaderboard.UpdateUser(ctx, userId, int64(changes.coins)); err != nil {
		log.Error().Err(err).Msg("Failed to update leaderboard")
	}

	if *uctx.sendReports {
		reports, err := internal.GetReports(ctx, u.r, userId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send updated reports")
			return
		}

		err = send(ctx, u.r, userId, &internal.ServerMessage{
			Id: xid.New().String(),
			Payload: &internal.ServerMessage_Reports_{
				Reports: &internal.ServerMessage_Reports{
					Reports: reports,
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to send inital reports")
		}
	}

	if *uctx.sendStats {
		if err = sendStats(ctx, u.r, userId); err != nil {
			log.Error().Err(err).Msg("Failed to send stats")
		}
	}
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
