import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { platform } from "./config";
import { initApi } from "./api";

document.body.classList.add(`platform-${platform}`);

if (platform === "ios" || platform === "android") {
  import("./push-notifications");
}

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
  render();
};

run();
