export const platform: "web" | "ios" | "android" =
  typeof (window as any).Capacitor !== "undefined"
    ? (window as any).Capacitor.getPlatform()
    : "web";

const getWebWsUrl = () => {
  const isSecure = window.location.protocol === "https:";
  return `${isSecure ? "wss" : "ws"}://${window.location.host}/api/game/ws`;
};

const getString = (value: string | boolean | undefined) => {
  switch (typeof value) {
    case "string":
      return value;
    case "boolean":
      return value.toString();
    case "undefined":
      return "";
  }
};

const API_BASE_URL =
  platform === "android" || platform === "ios"
    ? getString(import.meta.env.VITE_APP_BASE_URL)
    : "/api";

const WS_URL =
  platform === "android"
    ? getString(import.meta.env.VITE_ANDROID_WS_URL)
    : platform === "ios"
    ? getString(import.meta.env.VITE_IOS_WS_URL)
    : getWebWsUrl();

export { API_BASE_URL, WS_URL };
