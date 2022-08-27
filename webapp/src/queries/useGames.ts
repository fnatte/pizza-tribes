import { useQuery, UseQueryOptions } from "react-query";
import { centralApiFetch } from "../api";

type Game = { id: string, title: string, status: string, joined: boolean };

export const useGames = (
  options?: UseQueryOptions<Game[], unknown>
) => {
  const queryResult = useQuery<Game[], unknown>(
    "games",
    async (): Promise<Game[]> => {
      const res = await centralApiFetch(`/games/`);
      const json = await res.json();
      return json as Game[];
    },
    {
      ...options,
      staleTime: 10 * 60 * 1000,
    }
  );

  return queryResult;
};
