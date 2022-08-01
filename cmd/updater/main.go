package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/fnatte/pizza-tribes/internal/ws"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	case models.Education_PUBLICIST:
		return "publicists", nil
	default:
		return "", fmt.Errorf("Invalid education: %s", edu)
	}
}

type updater struct {
	r           redis.RedisClient
	world       *internal.WorldService
	leaderboard *internal.LeaderboardService
	gsRepo      persist.GameStateRepository
	reportsRepo persist.ReportsRepository
	updater     gamestate.Updater
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
	 *   - schedule the next update
	 */

	tx, err := u.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		if err := extrapolate(userId, gs, tx); err != nil {
			return err
		}
		if err := completedConstructions(userId, gs, tx); err != nil {
			return err
		}
		if err := completeTrainings(userId, gs, tx); err != nil {
			return err
		}
		if err := completeTravels(ctx, userId, gs, tx, u.r, u.world); err != nil {
			return err
		}
		if err := completeResearchs(userId, gs, tx); err != nil {
			return err
		}
		if err := completeQuests(userId, gs, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to perform update")
		return
	}

	for txUserId, txu := range tx.Users {
		if err := u.postProcessPatch(ctx, txUserId, txu); err != nil {
			log.Error().Err(err).
				Bool("wasSender", txUserId == userId).
				Str("userId", txUserId).
				Msg("Failed to process patch")
		}
	}

	u.scheduleNextUpdate(ctx, userId, tx.Users[userId].Gs)
}

func (u *updater) scheduleNextUpdate(ctx context.Context, userId string, gs *models.GameState) {
	_, err := internal.SetNextUpdate(u.r, ctx, userId, gs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set next update")
	}
}

func (u *updater) postProcessPatch(ctx context.Context, userId string, txu *gamestate.GameTx_User) error {
	var err error

	// Send game state patch
	err = ws.Send(ctx, u.r, userId, txu.ToServerMessage())
	if err != nil {
		log.Error().Err(err).Msg("Failed to send state change")
	}

	// Update leaderboard
	if txu.CoinsChanged {
		coins := int64(txu.Gs.Resources.Coins)
		if err = u.leaderboard.UpdateUser(ctx, userId, coins); err != nil {
			log.Error().Err(err).Msg("Failed to update leaderboard")
		}
	}

	if err := addTimeseriesDataPoints(ctx, userId, txu, u.r); err != nil {
		log.Error().Err(err).Msg("Failed to add timeseries data points")
	}

	// Send reports
	if txu.ReportsInvalidated {
		if err = sendReports(ctx, u.r, userId); err != nil {
			log.Error().Err(err).Msg("Failed to send inital reports")
		}
	}

	// Send stats
	if txu.StatsInvalidated {
		if err = sendStats(ctx, u.r, userId); err != nil {
			log.Error().Err(err).Msg("Failed to send stats")
		}
	}

	// Handle win
	if txu.CoinsChanged {
		if txu.Gs.Resources.Coins >= 10_000_000 {
			// End the world
			worldState, err := u.world.End(ctx, userId)
			if err != nil {
				log.Error().Err(err).Msg("Failed to end game")
			}

			// Announce new world state to all connected players
			err = ws.Send(ctx, u.r, "everyone", &models.ServerMessage{
				Id: xid.New().String(),
				Payload: &models.ServerMessage_WorldState{
					WorldState: worldState,
				},
			})
		}
	}

	return nil
}

func addTimeseriesDataPoints(ctx context.Context, userId string, txu *gamestate.GameTx_User, r redis.RedisClient) error {
	if !txu.CoinsChanged && !txu.PizzasChanged {
		return nil
	}

	err := internal.EnsureTimeseries(ctx, r, userId)
	if err != nil {
		return fmt.Errorf("failed to ensure timeseries: %w", err)
	}

	time := txu.Gs.Timestamp * 1000
	coins := txu.Gs.Resources.Coins
	pizzas := txu.Gs.Resources.Pizzas

	if txu.CoinsChanged {
		err = internal.AddMetricCoins(ctx, r, userId, time, int64(coins))
		if err != nil {
			return fmt.Errorf("failed to add coins metric: %w", err)
		}
	}

	if txu.PizzasChanged {
		err = internal.AddMetricPizzas(ctx, r, userId, time, int64(pizzas))
		if err != nil {
			return fmt.Errorf("failed to add pizzas metric: %w", err)
		}
	}

	return nil
}

func sendReports(ctx context.Context, r redis.RedisClient, userId string) error {
	reports, err := internal.GetReports(ctx, r, userId)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send updated reports")
		return fmt.Errorf("failed to get reports: %w", err)
	}

	return ws.Send(ctx, r, userId, &models.ServerMessage{
		Id: xid.New().String(),
		Payload: &models.ServerMessage_Reports_{
			Reports: &models.ServerMessage_Reports{
				Reports: reports,
			},
		},
	})
}

func sendStats(ctx context.Context, r redis.RedisClient, userId string) error {
	gs := models.GameState{}
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)
	s, err := redis.RedisJsonGet(r, ctx, gsKey, ".").Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
		return err
	}

	msg := internal.CalculateStats(&gs).ToServerMessage()

	return ws.Send(ctx, r, userId, msg)
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

func envOrDefault(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func handleStarted(ctx context.Context, u *updater) {
	userId, err := u.next(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get next")
		time.Sleep(1 * time.Second)
		return
	}
	if userId == "" {
		// TODO: could be smarter about how long to sleep here
		time.Sleep(10 * time.Millisecond)
		return
	}

	u.update(ctx, userId)
	time.Sleep(1 * time.Millisecond)
}

func handleStarting(ctx context.Context, world *internal.WorldService, worldState *models.WorldState) {
	if time.Now().Unix() >= worldState.StartTime {
		log.Info().Msg("Starting new round")
		// Let's get this game started.
		err := world.Start(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to start world")
			time.Sleep(1 * time.Second)
			return
		}
	} else {
		time.Sleep(1 * time.Second)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Info().Msg("Starting updater")

	debug := envOrDefault("DEBUG", "0") == "1"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	rc := redis.NewRedisClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	world := internal.NewWorldService(rc)
	leaderboard := internal.NewLeaderboardService(rc)
	gsRepo := persist.NewGameStateRepository(rc)
	reportsRepo := persist.NewReportsRepository(rc)
	userRepo := persist.NewUserRepository(rc)
	notifyRepo := persist.NewNotifyRepository(rc)
	u2 := gamestate.NewUpdater(gsRepo, reportsRepo, userRepo, notifyRepo)
	u := updater{r: rc, world: world, leaderboard: leaderboard, gsRepo: gsRepo, reportsRepo: reportsRepo, updater: u2}

	ctx := context.Background()

	for {
		worldState, err := world.GetState(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get world state")
			time.Sleep(1 * time.Second)
			continue
		}

		switch worldState.Type.(type) {
		case *models.WorldState_Started_:
			handleStarted(ctx, &u)
			break
		case *models.WorldState_Starting_:
			handleStarting(ctx, world, worldState)
			break
		case *models.WorldState_Ended_:
			time.Sleep(1 * time.Second)
			break
		}
	}
}
