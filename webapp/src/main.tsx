import React, { StrictMode } from "react";
import { createRoot } from "react-dom/client";
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
  const element = document.getElementById("root");
  if (!element) {
    console.warn("Could not find target root element");
    return;
  }

  const root = createRoot(element);
  root.render(
    <StrictMode>
      <App />
    </StrictMode>
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
