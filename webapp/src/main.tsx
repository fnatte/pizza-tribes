import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { platform } from "./config";

document.body.classList.add(`platform-${platform}`);

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById("root")
);
