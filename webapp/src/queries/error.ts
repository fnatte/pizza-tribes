export type ApiErrorCode = "unauthorized" | "unknown";

export class ApiError extends Error {
  readonly code: ApiErrorCode;

  constructor(code: ApiErrorCode, message: string) {
    super(message);
    Object.setPrototypeOf(this, ApiError.prototype);

    this.code = code;
  }
}

export function checkError(res: Response) {
  if (!res.ok) {
    if (res.status === 401 || res.status === 403) {
      throw new ApiError(
        "unauthorized",
        "Failed to fetch because request was unauthorized"
      );
    }

    throw new ApiError("unknown", "Failed to fetch");
  }
}
