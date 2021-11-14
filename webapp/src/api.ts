import { API_BASE_URL, platform } from "./config";
import { Storage } from "@capacitor/storage";

let _accessToken: null | string = null;

export const initApi = async () => {
  const result = await Storage.get({ key: "accessToken" });
  _accessToken = result.value;
};

export const setAccessToken = (accessToken: string) => {
  _accessToken = accessToken;
  Storage.set({ key: "accessToken", value: _accessToken });
};

export const getAccessToken = () => _accessToken;

const getAuthHeaders = () => {
  if (platform === "android" || platform === "ios") {
    return {
      Authorization: `Bearer ${getAccessToken()}`,
    };
  }
};

export const apiFetch = (path: string, init?: RequestInit) => {
  return fetch(`${API_BASE_URL}${path[0] === "/" ? "" : "/"}${path}`, {
    credentials: "same-origin",
    headers: {
      ...getAuthHeaders(),
      ...init?.headers,
    },
    ...init,
  });
};
