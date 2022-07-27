package main

import (
	"context"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/persist"
	"github.com/rs/zerolog/log"
)

func makeTapReminderMessage(userId string) *messaging.Message {
	return &messaging.Message{
		Data: map[string]string{
			"userId": userId,
		},
		Notification: &messaging.Notification{
			Title: "Keep your tap streak alive!",
		},
		Android: &messaging.AndroidConfig{
			CollapseKey: "reminder",
		},
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Tag: "reminder",
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-collapse-id": "reminder",
			},
		},
	}
}

func handleTapReminder(ctx context.Context, rc internal.RedisClient, u persist.UserRepository, r *internal.Reminder) {
	users, err := u.GetAllUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all users when handling tap reminder")
	}

	for _, user := range users {
		t, err := u.GetUserLatestActivity(ctx, user)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user's latest activity when handling tap reminder")
			continue
		}

		// Check if user has not been active this hour
		if time.Unix(0, t).Before(time.Now().Truncate(time.Hour).Add(-time.Minute)) {

			log.Debug().Str("userId", user).Msg("Scheduling tap reminder push notification")

			internal.SchedulePushNotification(ctx, rc, makeTapReminderMessage(user), time.Now())
		}
	}
}

func startRemindersWorker(ctx context.Context, rc internal.RedisClient, u persist.UserRepository) {
	internal.ScheduleReminder(ctx, rc, &internal.Reminder{
		Id:     "tap-reminder",
		Interval: time.Hour,
		Offset:   20 * time.Minute,
	})

	go internal.HandleReminders(ctx, rc, func(r *internal.Reminder) {
		log.Debug().Str("id", r.Id).Msg("Handle reminder")
		switch r.Id {
		case "tap-reminder":
			handleTapReminder(ctx, rc, u, r)
			break
		}
	})
}
