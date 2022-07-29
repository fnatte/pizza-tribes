package main

import (
	"context"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/persist"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
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

	debug := envOrDefault("DEBUG", "0") == "1"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: envOrDefault("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})

	rc := internal.NewRedisClient(rdb)

	world := internal.NewWorldService(rc)

	gsRepo := persist.NewGameStateRepository(rc)
	reportsRepo := persist.NewReportsRepository(rc)
	userRepo := persist.NewUserRepository(rc)
	notifyRepo := persist.NewNotifyRepository(rc)
	updater := gamestate.NewUpdater(gsRepo, reportsRepo, userRepo, notifyRepo)

	h := &handler{rdb: rc, world: world, gsRepo: gsRepo, reportsRepo: reportsRepo, userRepo: userRepo, updater: updater}

	ctx := context.Background()

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
		startRemindersWorker(ctx, rc, userRepo)
	}

	for {
		res, err := rdb.BLPop(ctx, 30*time.Second, "wsin").Result()
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

		msg := &internal.IncomingMessage{}
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
