import { ApiError } from "./error";

export function retry(count: number, err: unknown): boolean {
  if (err instanceof ApiError && err.code === "unauthorized") {
    return false;
  }
  return count < 3;
}
