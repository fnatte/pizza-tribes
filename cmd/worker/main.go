package main

import (
	"context"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func envOrDefault(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func main() {
	log.Info().Msg("Starting worker")

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

	world := game.NewWorldService(rc)

	gsRepo := persist.NewGameStateRepository(rc)
	reportsRepo := persist.NewReportsRepository(rc)
	userRepo := persist.NewGameUserRepository(rc)
	notifyRepo := persist.NewNotifyRepository(rc)
	worldRepo := persist.NewWorldRepository(rc)
	marketRepo := persist.NewMarketRepository(rc)
	updater := gamestate.NewUpdater(gsRepo, reportsRepo, userRepo, notifyRepo, marketRepo, worldRepo)

	speed, err := world.GetSpeed(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get game speed")
		return
	}

	game.AlterGameDataForSpeed(speed)

	h := &handler{rdb: rc, world: world, gsRepo: gsRepo, reportsRepo: reportsRepo,
		userRepo: userRepo, updater: updater, marketRepo: marketRepo, speed: speed}

	useFirebase := envOrDefault("FEATURE_FIREBASE", "0") == "1"
	usePushNotifications := envOrDefault("FEATURE_PUSH_NOTIFICATIONS", "0") == "1"
	useReminders := envOrDefault("FEATURE_REMINDERS", "0") == "1"
	mockPushNotifications := envOrDefault("MOCK_PUSH_NOTIFICATIONS", "0") == "1"

	var fapp *firebase.App
	var msgClient *messaging.Client
	if useFirebase {
		var err error
		fapp, err = firebase.NewApp(context.Background(), nil)
		if err != nil {
			log.Fatal().Err(err).Msg("Error when initializing firebase app")
		}

		msgClient, err = fapp.Messaging(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Error when initializing messaging client")
		}
	}

	if usePushNotifications {
		if !useFirebase && !mockPushNotifications {
			log.Warn().Msg("Push notifications are not supported without enabling firebase")
		} else {
			log.Info().Msg("Starting Push Notification worker")
			go pushNotificationsWorker(ctx, msgClient, rc, mockPushNotifications)
		}
	}

	if useReminders {
		log.Info().Msg("Starting reminders worker")
		startRemindersWorker(ctx, rc, userRepo, worldRepo)
	}

	for {
		res, err := rc.BLPop(ctx, 30*time.Second, "wsin").Result()
		if err != nil {
			if err != redis.Nil {
				log.Error().Err(err).Msg("Error when dequeuing message")
			}
			continue
		}

		if len(res) < 2 {
			log.Error().Err(err).Msg("This should never happend. BLPop should always return a slice with two values")
			continue
		}

		msg := &game.IncomingMessage{}
		msg.UnmarshalBinary([]byte(res[1]))

		m := &models.ClientMessage{}
		err = protojson.Unmarshal([]byte(msg.Body), m)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse incoming message")
			continue
		}

		h.Handle(ctx, msg.SenderId, m)
	}

}
