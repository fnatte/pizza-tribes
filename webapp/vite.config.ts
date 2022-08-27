import { defineConfig } from "vite";
import reactRefresh from "@vitejs/plugin-react-refresh";
import svgr from "./plugins/svgr";
import fs from "fs";
import path from "path";

// We set clearScreen to false in case of a configuration error occur, so that the user have
// a chance to see any logged error messages.
let clearScreen = true;

const defaultOrigin = "http://localhost:3000";

const readOriginsFromEnvFile = (envFile: string) => {
  try {
    const data = fs.readFileSync(envFile, "utf8");
    const lines = data.split("\n");
    const origins = lines
      .find((line) => line.startsWith("ORIGIN="))
      .substr("ORIGIN=".length)
      .split(" ");

    return origins;
  } catch (e) {
    console.error(
      `Failed to read ORIGIN value from env file at: ${envFile}. Using the default value "http://localhost:3000".`
    );
    clearScreen = false;
    return defaultOrigin;
  }
};

// Read origin from .env-file in parent folder so that we use the same value that the real backend is using.
const origin = readOriginsFromEnvFile("../.env");

export default defineConfig({
  resolve: {
    alias: {
      images: path.resolve(__dirname, "./images"),
    },
  },
  clearScreen,
  server: {
    proxy: {
      "/api/game": {
        target: "http://localhost:8080",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/game/, ""),
        ws: true,
      },
      "/api/central": {
        target: "http://localhost:8083",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/central/, ""),
      },
    },
    cors: {
      origin,
      allowedHeaders: [
        "Accept",
        "Accept-Language",
        "Authorization",
        "Content-Language",
        "Content-Type",
        "Origin",
      ],
      methods: ["GET", "HEAD", "POST"],
      credentials: true,
    },
  },
  plugins: [reactRefresh(), svgr()],
});
