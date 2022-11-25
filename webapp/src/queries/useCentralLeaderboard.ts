import { useQuery, UseQueryOptions } from "react-query";
import { centralApiFetch } from "../api";
import { checkError } from "./error";
import { retry } from "./retry";

type CentralLeaderboardRow = { 
  userId: string,
  username: string,
  coins: string
}

type CentralLeaderboard = { 
  gameId: string,
  skip: number,
  limit: number,
  rows: CentralLeaderboardRow[],
};

export const useCentralLeaderboard = (gameId: string, skip: number, options?: UseQueryOptions<CentralLeaderboard, unknown>) => {
  const queryResult = useQuery<CentralLeaderboard, unknown>(
    ["games", gameId, "leaderboard", skip],
    async (): Promise<CentralLeaderboard> => {
      const res = await centralApiFetch(`/games/${gameId}/leaderboard?skip=${skip}`);
      checkError(res);
      const json = await res.json();
      return json as CentralLeaderboard;
    },
    {
      ...options,
      retry,
      staleTime: 10 * 60 * 1000,
    }
  );

  return queryResult;
};
