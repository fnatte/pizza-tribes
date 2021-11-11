import {
  ActionPerformed,
  PushNotificationSchema,
  PushNotifications,
  Token,
} from "@capacitor/push-notifications";
import { API_BASE_URL } from "./config";

// Request permission to use push notifications
// iOS will prompt user and return if they granted permission or not
// Android will just grant without prompting
PushNotifications.requestPermissions().then((result) => {
  if (result.receive === "granted") {
    // Register with Apple / Google to receive push via APNS/FCM
    PushNotifications.register();
  } else {
    // Show some error
  }
});

// On success, we should be able to receive notifications
PushNotifications.addListener("registration", async (token: Token) => {
  console.log("Push registration success, token: ", token.value);
  const response = await fetch(API_BASE_URL + "/push_notifications/register", {
    method: "POST",
    body: JSON.stringify({ token: token.value }),
    credentials: "include",
  });
  if (!response.ok) {
    console.error("Failed to register fcm token");
  }
});

// Some issue with our setup and push will not work
PushNotifications.addListener("registrationError", (error: any) => {
  console.log("Error on registration: ", error);
});

// Show us the notification payload if the app is open on our device
PushNotifications.addListener(
  "pushNotificationReceived",
  (notification: PushNotificationSchema) => {
    console.log("Push received: ", notification);
  }
);

// Method called when tapping on a notification
PushNotifications.addListener(
  "pushNotificationActionPerformed",
  (notification: ActionPerformed) => {
    console.log("Push action performed: ", notification);
  }
);
