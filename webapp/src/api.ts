import { API_BASE_URL, platform } from "./config";
import { Preferences } from "@capacitor/preferences";

let _accessToken: null | string = null;

export const initApi = async () => {
  const result = await Preferences.get({ key: "accessToken" });
  _accessToken = result.value;
};

export const setAccessToken = (accessToken: string) => {
  _accessToken = accessToken;
  Preferences.set({ key: "accessToken", value: _accessToken });
};

export const removeAccessToken = () => {
  _accessToken = null;
  Preferences.remove({ key: "accessToken" });
};

export const getAccessToken = () => _accessToken;

const getAuthHeaders = () => {
  if (platform === "android" || platform === "ios" || isCypress()) {
    const accessToken = getAccessToken();
    if (accessToken) {
      return {
        Authorization: `Bearer ${accessToken}`,
      };
    }
  }
};

export const centralApiFetch = (path: string, init?: RequestInit) => {
  return fetch(`${API_BASE_URL}/central${path[0] === "/" ? "" : "/"}${path}`, {
    credentials: "same-origin",
    headers: {
      ...getAuthHeaders(),
      ...init?.headers,
    },
    ...init,
  });
};

export const apiFetch = (path: string, init?: RequestInit) => {
  return fetch(`${API_BASE_URL}/game${path[0] === "/" ? "" : "/"}${path}`, {
    credentials: "same-origin",
    headers: {
      ...getAuthHeaders(),
      ...init?.headers,
    },
    ...init,
  });
};

function isCypress(): boolean {
  return "Cypress" in window;
}
