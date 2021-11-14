import { CapacitorConfig } from "@capacitor/cli";

const baseConfig: CapacitorConfig = {
  appId: "dev.teus.pizzatribes",
  appName: "Pizza Tribes",
  webDir: "dist",
  bundledWebRuntime: false,
  appendUserAgent: "pizzatribes",
};

const liveReload = {
  url: "http://localhost:3000",
  cleartext: true,
};

const devConfig: CapacitorConfig = {
  server: {
    ...(process.env.LIVE_RELOAD === "1" ? liveReload : {}),
  },
  android: {
    flavor: "dev",
  },
  plugins: {
    Keyboard: {
      resize: "native",
    },
  },
};

const prodConfig: CapacitorConfig = {
  android: {
    flavor: "prod",
  },
};

const config: CapacitorConfig = {
  ...baseConfig,
  ...(process.env.NODE_ENV === "development" ? devConfig : prodConfig),
};

export default config;
