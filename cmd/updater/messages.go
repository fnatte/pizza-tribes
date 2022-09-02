package main

import "firebase.google.com/go/messaging"

func makeGameStartedMessage(userId string) *messaging.Message {
	return &messaging.Message{
		Data: map[string]string{
			"userId": userId,
		},
		Notification: &messaging.Notification{
			Title: "And we're off. Game on!",
		},
		Android: &messaging.AndroidConfig{
			CollapseKey: "gameround",
		},
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Tag: "gameround",
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-collapse-id": "gameround",
			},
		},
	}
}
