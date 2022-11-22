import { useQuery, UseQueryOptions } from "react-query";
import { apiFetch } from "../api";
import { DemandLeaderboard } from "../generated/leaderboard";
import { checkError } from "./error";
import { retry } from "./retry";

export const useDemandLeaderboard = (
  skip?: number,
  options?: UseQueryOptions<DemandLeaderboard, unknown, DemandLeaderboard, string[]>
) => {
  const queryResult = useQuery(
    ["demand-leaderboard", "skip", skip?.toString() ?? ""],
    async (): Promise<DemandLeaderboard> => {
      let skipParam = skip;

      if (skipParam === undefined) {
        const res = await apiFetch("/demand_leaderboard/me/rank");
        checkError(res);
        const rank = await res.json();
        skipParam = Math.max(0, rank - 10);
      }

      const res = await apiFetch(`/demand_leaderboard/?skip=${skipParam}`);
      checkError(res);
      return DemandLeaderboard.fromJson(await res.json());
    },
    {
      ...options,
      retry,
    }
  );

  return queryResult;
};
