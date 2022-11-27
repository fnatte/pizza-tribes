package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	"github.com/fnatte/pizza-tribes/internal/game/ws"
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
	world       *game.WorldService
	leaderboard *game.LeaderboardService
	gsRepo      persist.GameStateRepository
	reportsRepo persist.ReportsRepository
	worldRepo   persist.WorldRepository
	marketRepo  persist.MarketRepository
	userRepo    persist.GameUserRepository
	updater     gamestate.Updater
	speed       float64
}

// Update the game state for the specified user
func (u *updater) update(ctx context.Context, userId string) {
	worldState, err := u.worldRepo.GetState(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform update")
		u.scheduleNextUpdateAfterError(ctx, userId)
		return
	}

	globalDemandScore, err := u.marketRepo.GetGlobalDemandScore(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform update")
		u.scheduleNextUpdateAfterError(ctx, userId)
		return
	}

	userCount, err := u.userRepo.GetUserCount(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform update")
		u.scheduleNextUpdateAfterError(ctx, userId)
		return
	}

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
		if err := extrapolate(userId, gs, tx, worldState, globalDemandScore, userCount, u.speed); err != nil {
			return err
		}
		if err := completedConstructions(userId, gs, tx); err != nil {
			return err
		}
		if err := completeTrainings(userId, gs, tx); err != nil {
			return err
		}
		if err := completeTravels(ctx, userId, gs, tx, u.r, u.world, u.speed); err != nil {
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
		u.scheduleNextUpdateAfterError(ctx, userId)
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

func (u *updater) scheduleNextUpdateAfterError(ctx context.Context, userId string) {
	_, err := game.SetNextUpdateTo(u.r, ctx, userId, time.Now().Add(2*time.Second).UnixNano())
	if err != nil {
		log.Error().Err(err).Msg("Failed to set next update after error")
	}
}

func (u *updater) scheduleNextUpdate(ctx context.Context, userId string, gs *models.GameState) {
	_, err := game.SetNextUpdate(u.r, ctx, userId, gs)
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
		log.Error().Err(err).
			Str("userId", userId).
			Msg("Failed to add timeseries data points")
	}

	// Send reports
	if txu.ReportsInvalidated {
		if err = sendReports(ctx, u.r, userId); err != nil {
			log.Error().Err(err).Msg("Failed to send inital reports")
		}
	}

	if txu.StatsInvalidated {
		// Send stats
		var stats *models.Stats
		if stats, err = u.getStats(ctx, userId); err != nil {
			log.Error().Err(err).Msg("Failed to get stats")
		}

		msg := stats.ToServerMessage()
		if err = ws.Send(ctx, u.r, userId, msg); err != nil {
			log.Error().Err(err).Msg("Failed to send stats")
		}

		// Update user demand score
		demandScore := game.CalculateDemandScore(txu.Gs)
		if err = u.marketRepo.SetUserDemandScore(ctx, userId, demandScore); err != nil {
			log.Error().Err(err).Msg("Failed to set user demand score")
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

	err := game.EnsureTimeseries(ctx, r, userId)
	if err != nil {
		return fmt.Errorf("failed to ensure timeseries: %w", err)
	}

	time := txu.Gs.Timestamp * 1000
	coins := txu.Gs.Resources.Coins
	pizzas := txu.Gs.Resources.Pizzas

	if txu.CoinsChanged {
		err = game.AddMetricCoins(ctx, r, userId, time, int64(coins))
		if err != nil {
			return fmt.Errorf("failed to add coins metric: %w", err)
		}
	}

	if txu.PizzasChanged {
		err = game.AddMetricPizzas(ctx, r, userId, time, int64(pizzas))
		if err != nil {
			return fmt.Errorf("failed to add pizzas metric: %w", err)
		}
	}

	return nil
}

func sendReports(ctx context.Context, r redis.RedisClient, userId string) error {
	reports, err := game.GetReports(ctx, r, userId)
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

func (u *updater) getStats(ctx context.Context, userId string) (*models.Stats, error) {
	gs, err := u.gsRepo.Get(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to send stats: %w", err)
	}

	worldState, err := u.worldRepo.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to send full state update: %w", err)
	}

	globalDemandScore, err := u.marketRepo.GetGlobalDemandScore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to send full state update: %w", err)
	}

	userCount, err := u.userRepo.GetUserCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to send full state update: %w", err)
	}

	speed, err := u.world.GetSpeed(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get game speed: %w", err)
	}

	return game.CalculateStats(gs, globalDemandScore, worldState, userCount, speed), nil
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

func handleStarting(ctx context.Context, world *game.WorldService, worldState *models.WorldState, userRepo persist.GameUserRepository, r redis.RedisClient) {
	if time.Now().Unix() >= worldState.StartTime {
		log.Info().Msg("Starting new round")
		// Let's get this game started.
		err := world.Start(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to start world")
			time.Sleep(1 * time.Second)
			return
		}

		worldState, err = world.GetState(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get world state after started so it is not possible to annonuce")
			time.Sleep(1 * time.Second)
			return
		}

		// Announce new world state to all connected players
		err = ws.Send(ctx, r, "everyone", &models.ServerMessage{
			Id: xid.New().String(),
			Payload: &models.ServerMessage_WorldState{
				WorldState: worldState,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to annonuce started world state")
			return
		}

		// Send notifications
		users, err := userRepo.GetAllUsers(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get all users - not able to send notifications")
			return
		}
		log.Debug().Int("numberOfUsers", len(users)).Msg("Scheduling activity reminder push notifications")
		for _, u := range users {
			// TODO: might be able to batch this into 500 notification chunks
			game.SchedulePushNotification(ctx, r, makeGameStartedMessage(u), time.Now())
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

	ctx := context.Background()

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

	redisDebug := envOrDefault("REDIS_DEBUG", "0") == "1"
	if redisDebug {
		rc.AddDebugHook()
	}

	world := game.NewWorldService(rc)
	leaderboard := game.NewLeaderboardService(rc)
	gsRepo := persist.NewGameStateRepository(rc)
	reportsRepo := persist.NewReportsRepository(rc)
	userRepo := persist.NewGameUserRepository(rc)
	notifyRepo := persist.NewNotifyRepository(rc)
	worldRepo := persist.NewWorldRepository(rc)
	marketRepo := persist.NewMarketRepository(rc)
	u2 := gamestate.NewUpdater(gsRepo, reportsRepo, userRepo, notifyRepo, marketRepo, worldRepo)

	speed, err := world.GetSpeed(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get world speed")
		return
	}

	game.AlterGameDataForSpeed(speed)

	u := updater{
		r: rc, world: world, leaderboard: leaderboard, gsRepo: gsRepo,
		reportsRepo: reportsRepo, worldRepo: worldRepo, marketRepo: marketRepo,
		userRepo: userRepo, updater: u2, speed: speed}

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
			handleStarting(ctx, world, worldState, userRepo, rc)
			break
		case *models.WorldState_Ended_:
			time.Sleep(1 * time.Second)
			break
		}
	}
}
