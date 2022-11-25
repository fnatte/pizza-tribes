import { useQuery, UseQueryOptions } from "react-query";
import { centralApiFetch } from "../api";
import { checkError } from "./error";
import { retry } from "./retry";

export const useCentralLeaderboardMeRank = (gameId: string, options?: UseQueryOptions<number, unknown>) => {
  const queryResult = useQuery<number, unknown>(
    ["games", gameId, "leaderboard", "me", "rank"],
    async (): Promise<number> => {
      const res = await centralApiFetch(`/games/${gameId}/leaderboard/me/rank`);
      checkError(res);
      const json = await res.json();
      return json as number;
    },
    {
      ...options,
      retry,
      staleTime: 10 * 60 * 1000,
    }
  );

  return queryResult;
};
