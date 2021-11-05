const isNativePlatform =
  typeof (window as any).Capacitor !== "undefined" &&
  (window as any).Capacitor.isNativePlatform();

const getWebWsUrl = () => {
  const isSecure = window.location.protocol === "https:";
  return `${isSecure ? "wss" : "ws"}://${window.location.host}/api/ws`;
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

const API_BASE_URL = isNativePlatform
  ? getString(import.meta.env.VITE_APP_BASE_URL)
  : "/api";
const WS_URL = isNativePlatform
  ? getString(import.meta.env.VITE_APP_WS_URL)
  : getWebWsUrl();

export { API_BASE_URL, WS_URL };
