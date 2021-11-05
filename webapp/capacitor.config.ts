import { CapacitorConfig } from "@capacitor/cli";

const baseConfig: CapacitorConfig = {
  appId: "dev.teus.pizzatribes",
  appName: "Pizza Tribes",
  webDir: "dist",
  bundledWebRuntime: false,
}

const devConfig: CapacitorConfig = {
  server: {
    url: "http://10.0.2.2:3000",
    cleartext: true,
  },
  android: {
   flavor: "dev",
 },
};

const prodConfig: CapacitorConfig = {
  android: {
   flavor: "prod",
 },
};

const config: CapacitorConfig = {
  ...baseConfig,
  ...(process.env.NODE_ENV === "development" ? devConfig : prodConfig)
}

export default config;
