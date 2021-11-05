import { useQuery, UseQueryOptions } from "react-query";
import {API_BASE_URL} from "../config";

type User = { username: string };

export const useUser = (userId: string, options: UseQueryOptions<User, unknown, User, string[]>) => {
  const queryResult = useQuery(
    ["user", userId],
    async (): Promise<User> => {
      const res = await fetch(`${API_BASE_URL}/user/${userId}`);
      const json = await res.json();
      return { username: json.username };
    },
    {
      ...options,
      staleTime: 10 * 60 * 1000,
    }
  );

  return queryResult;
};
