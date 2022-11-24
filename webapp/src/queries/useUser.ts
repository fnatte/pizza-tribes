import { useQuery, UseQueryOptions } from "react-query";
import { apiFetch } from "../api";
import { checkError } from "./error";
import { retry } from "./retry";

type User = { username: string };

export const useUser = (
  userId: string,
  options?: UseQueryOptions<User, unknown, User, string[]>
) => {
  const queryResult = useQuery(
    ["user", userId],
    async (): Promise<User> => {
      const res = await apiFetch(`/user/${userId}`);
      checkError(res);
      const json = await res.json();
      return { username: json.username };
    },
    {
      ...options,
      retry,
      staleTime: 10 * 60 * 1000,
    }
  );

  return queryResult;
};
