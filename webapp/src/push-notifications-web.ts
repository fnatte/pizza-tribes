import { WebPlugin } from "@capacitor/core";
import {
  Channel,
  DeliveredNotifications,
  ListChannelsResult,
  PermissionStatus,
  PushNotificationsPlugin,
} from "@capacitor/push-notifications";

/**
 * This is an unused mock implementation. Push notifications are currently disabled for the web platform.
 * The class below may be used in the future to add support for web platform push notifications.
 */
export class PushNotificationsWeb
  extends WebPlugin
  implements PushNotificationsPlugin {
  register(): Promise<void> {
    console.log("register");
    this.notifyListeners("registration", {
      value: "test",
    });
    return Promise.resolve();
  }
  getDeliveredNotifications(): Promise<DeliveredNotifications> {
    console.log("getDeliveredNotifications");
    return Promise.resolve({
      notifications: [],
    });
  }
  removeDeliveredNotifications(
    delivered: DeliveredNotifications
  ): Promise<void> {
    console.log("removeDeliveredNotifications", delivered);
    return Promise.resolve();
  }
  removeAllDeliveredNotifications(): Promise<void> {
    console.log("removeAllDeliveredNotifications");
    return Promise.resolve();
  }
  createChannel(channel: Channel): Promise<void> {
    console.log("createChannel", channel);
    return Promise.resolve();
  }
  deleteChannel(channel: Channel): Promise<void> {
    console.log("deleteChannel", channel);
    return Promise.resolve();
  }
  listChannels(): Promise<ListChannelsResult> {
    console.log("listChannels");
    return Promise.resolve({
      channels: [],
    });
  }
  checkPermissions(): Promise<PermissionStatus> {
    console.log("checkPermissions");
    return Promise.resolve({
      receive: "denied",
    });
  }
  requestPermissions(): Promise<PermissionStatus> {
    console.log("requestPermissions");
    return Promise.resolve({
      receive: "granted",
    });
  }
}
