package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/fnatte/pizza-tribes/internal"
	"github.com/rs/zerolog/log"
)

func nextPushNotification(ctx context.Context, r internal.RedisClient) (*messaging.Message, error) {
	packed, err := r.ZRangeWithScores(ctx, "push_notifications", 0, 0).Result()
	if err != nil {
		return nil, err
	}

	if len(packed) == 0 {
		return nil, nil
	}

	timestamp := int64(packed[0].Score)
	if timestamp > time.Now().Unix() {
		return nil, nil
	}

	member, ok := packed[0].Member.(string)
	if !ok {
		return nil, errors.New("failed to read push notification member")
	}

	msg := &messaging.Message{}
	json.Unmarshal([]byte(member), msg)

	removed, err := r.ZRem(ctx, "push_notifications", member).Result()
	if err != nil {
		return nil, err
	}

	if removed != 1 {
		return nil, nil
	}

	return msg, nil
}

func pushNotificationsWorker(ctx context.Context, client *messaging.Client, r internal.RedisClient) {
	for {
		msg, err := nextPushNotification(ctx, r)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get next push notification")
			time.Sleep(1 * time.Second)
			continue
		}

		if msg == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		userId, ok := msg.Data["userId"]
		if !ok || userId == "" {
			log.Error().Err(err).Msg("Encountered push notification without target user id")
			time.Sleep(10 * time.Millisecond)
			continue
		}

		fcmTokens, err := r.SMembers(ctx, fmt.Sprintf("user:%s:fcm_tokens", userId)).Result()
		if err != nil {
			log.Error().Err(err).Str("userId", userId).Msg("Failed to retrive fcm tokens for user")
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if len(fcmTokens) == 0 {
			// If there are no fcm tokens registered for the user, lets reschedule this message
			// at a later time, at which hopefully the user might have registered.
			_, err = internal.SchedulePushNotification(ctx, r, msg, time.Now().Add(3*time.Minute))
			if err != nil {
				log.Error().Err(err).Msg("Failed to reschedule push notification when there was no fcm tokens - it will be lost")
			}
			log.Info().Msg("Rescheduled push notification because the user has no fcm tokens")
			continue
		}

		mm := &messaging.MulticastMessage{
			Tokens:       fcmTokens,
			Data:         msg.Data,
			Notification: msg.Notification,
			Android:      msg.Android,
			Webpush:      msg.Webpush,
			APNS:         msg.APNS,
		}

		resp, err := client.SendMulticast(ctx, mm)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send push notification")
			_, err = internal.SchedulePushNotification(ctx, r, msg, time.Now().Add(3*time.Second))
			if err != nil {
				log.Error().Err(err).Msg("Failed to reschedule push notification when send failed - it will be lost")
			}
			continue
		}

		log.Debug().Interface("response", resp).Msg("Sent push notification")
	}
}
