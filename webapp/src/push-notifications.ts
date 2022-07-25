import {
  ActionPerformed,
  PushNotifications,
  PushNotificationSchema,
  Token,
} from "@capacitor/push-notifications";
import { apiFetch } from "./api";
import { State, useStore } from "./store";
import { platform } from "./config";

export function isPushNotificationsSupported() {
  return platform === "ios" || platform === "android";
}

async function postRegistrationToken(token: Token) {
  console.log("Push registration success, token: ", token.value);
  const response = await apiFetch("/push_notifications/register", {
    method: "POST",
    body: JSON.stringify({ token: token.value }),
  });
  if (!response.ok) {
    console.error("Failed to register fcm token");
  }
}

function isLoggedIn(state?: State): boolean {
  return (state ?? useStore.getState()).user !== null;
}

export async function initializePushNotifications() {
  let token: Token | null = null;

  // On success, we should be able to receive notifications
  await PushNotifications.addListener("registration", async (t: Token) => {
    token = t;
    if (token !== null && isLoggedIn()) {
      postRegistrationToken(token);
    }
  });

  // Some issue with our setup and push will not work
  await PushNotifications.addListener("registrationError", (error: any) => {
    console.log("Error on registration: ", error);
    token = null;
  });

  // Show us the notification payload if the app is open on our device
  await PushNotifications.addListener(
    "pushNotificationReceived",
    (notification: PushNotificationSchema) => {
      console.log("Push received: ", notification);
    }
  );

  // Method called when tapping on a notification
  await PushNotifications.addListener(
    "pushNotificationActionPerformed",
    (notification: ActionPerformed) => {
      console.log("Push action performed: ", notification);
    }
  );

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

  useStore.subscribe((state, previousState) => {
    if (
      token !== null &&
      state.user !== previousState.user &&
      isLoggedIn(state)
    ) {
      postRegistrationToken(token);
    }
  });
}
