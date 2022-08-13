package main

import (
	"context"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/persist"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"github.com/rs/zerolog/log"
)

func makeKeepStreakMessage(userId string) *messaging.Message {
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

func makeActivityReminderMessage(userId string) *messaging.Message {
	return &messaging.Message{
		Data: map[string]string{
			"userId": userId,
		},
		Notification: &messaging.Notification{
			Title: "Your tribe asks for your guidance",
			Body:  "Boss! Things are getting out of hand. Can you help us?",
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

func makeAvailableTapsMessage(userId string) *messaging.Message {
	return &messaging.Message{
		Data: map[string]string{
			"userId": userId,
		},
		Notification: &messaging.Notification{
			Title: "Pizza Tribes",
			Body:  "You have available taps!",
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

func handleTapReminder(ctx context.Context, rc redis.RedisClient, u persist.UserRepository, r *internal.Reminder) {
	users, err := u.GetAllUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all users when handling tap reminder")
	}

	for _, user := range users {
		lastActivity, err := u.GetUserLatestActivity(ctx, user)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user's latest activity when handling tap reminder")
			continue
		}

		t := time.Unix(0, lastActivity)

		// Don't send tap notifications if user has been inactive for over 12 hours
		if t.Before(time.Now().Add(-12 * time.Hour)) {
			continue
		}

		// Check if user has not been active this hour
		if t.Before(time.Now().Truncate(time.Hour).Add(-time.Minute)) {
			log.Debug().Str("userId", user).Msg("Scheduling tap reminder push notification")

			// If user was not active last hour either the streak is over,
			// so send a more suitable push notification.
			if t.Before(time.Now().Truncate(time.Hour).Add(-time.Hour - time.Minute)) {
				internal.SchedulePushNotification(ctx, rc, makeAvailableTapsMessage(user), time.Now())
			} else {
				internal.SchedulePushNotification(ctx, rc, makeKeepStreakMessage(user), time.Now())
			}
		}
	}
}

func handleActivityReminder(ctx context.Context, rc redis.RedisClient, u persist.UserRepository, r *internal.Reminder) {
	users, err := u.GetAllUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all users when handling activity reminder")
	}

	for _, user := range users {
		lastActivity, err := u.GetUserLatestActivity(ctx, user)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get user's latest activity when handling activity reminder")
			continue
		}

		t := time.Unix(0, lastActivity)

		// Don't send activity reminder if user has been inactive for over 72 hours
		if t.Before(time.Now().Add(-72 * time.Hour)) {
			continue
		}

		// Check if user has not been active for over 24 hours
		if t.Before(time.Now().Add(-24 * time.Hour)) {
			log.Debug().Str("userId", user).Msg("Scheduling activity reminder push notification")
			internal.SchedulePushNotification(ctx, rc, makeActivityReminderMessage(user), time.Now())
		}
	}
}

func startRemindersWorker(ctx context.Context, rc redis.RedisClient, u persist.UserRepository) {
	internal.ScheduleReminder(ctx, rc, &internal.Reminder{
		Id:       "tap-reminder",
		Interval: time.Hour,
		Offset:   20 * time.Minute,
	})
	internal.ScheduleReminder(ctx, rc, &internal.Reminder{
		Id:       "activity-reminder",
		Interval: time.Hour,
		Offset:   40 * time.Minute,
	})

	go internal.HandleReminders(ctx, rc, func(r *internal.Reminder) {
		log.Debug().Str("id", r.Id).Msg("Handle reminder")
		switch r.Id {
		case "tap-reminder":
			handleTapReminder(ctx, rc, u, r)
			break
		case "activity-reminder":
			handleActivityReminder(ctx, rc, u, r)
			break
		}
	})
}
