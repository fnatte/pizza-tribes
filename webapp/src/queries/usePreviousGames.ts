import { useQuery, UseQueryOptions } from "react-query";
import { centralApiFetch } from "../api";
import { checkError } from "./error";
import { retry } from "./retry";

type Game = { id: string; title: string; status: string; joined: boolean };

export const usePreviousGames = (options?: UseQueryOptions<Game[], unknown>) => {
  const queryResult = useQuery<Game[], unknown>(
    "games/previous",
    async (): Promise<Game[]> => {
      const res = await centralApiFetch(`/games/previous`);
      checkError(res);
      const json = await res.json();
      return json as Game[];
    },
    {
      ...options,
      retry,
      staleTime: 10 * 60 * 1000,
    }
  );

  return queryResult;
};
