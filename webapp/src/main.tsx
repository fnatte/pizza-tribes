import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { platform } from "./config";

document.body.classList.add(`platform-${platform}`);

if (platform === "ios") {
  import("./push-notifications");
}

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root")
);
