import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { platform } from "./config";
import { initApi } from "./api";
import {
  initializePushNotifications,
  isPushNotificationsSupported,
} from "./push-notifications";
import("./push-notifications");

document.body.classList.add(`platform-${platform}`);

const render = () => {
  ReactDOM.render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
    document.getElementById("root")
  );
};

const run = async () => {
  await initApi();
  if (isPushNotificationsSupported()) {
    initializePushNotifications();
  }
  render();
};

run();
