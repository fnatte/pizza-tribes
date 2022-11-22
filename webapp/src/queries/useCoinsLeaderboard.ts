import { useQuery, UseQueryOptions } from "react-query";
import { apiFetch } from "../api";
import { Leaderboard } from "../generated/leaderboard";
import { checkError } from "./error";
import { retry } from "./retry";

export const useCoinsLeaderboard = (
  skip?: number,
  options?: UseQueryOptions<Leaderboard, unknown, Leaderboard, string[]>
) => {
  const queryResult = useQuery(
    ["coins-leaderboard", "skip", skip?.toString() ?? ""],
    async (): Promise<Leaderboard> => {
      let skipParam = skip;

      if (skipParam === undefined) {
        const res = await apiFetch("/leaderboard/me/rank");
        checkError(res);
        const rank = await res.json();
        skipParam = Math.max(0, rank - 10);
      }

      const res = await apiFetch(`/leaderboard/?skip=${skipParam}`);
      checkError(res);
      return Leaderboard.fromJson(await res.json());
    },
    {
      ...options,
      retry,
      staleTime: 1 * 60 * 1000,
    }
  );

  return queryResult;
};
