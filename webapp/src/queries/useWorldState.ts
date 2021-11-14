import { useQuery } from "react-query";
import { apiFetch } from "../api";
import { WorldState } from "../generated/world";

export const useWorldState = () => {
  const queryResult = useQuery(
    "worldState",
    async (): Promise<WorldState> => {
      const res = await apiFetch("/world/");
      if (!res.ok) {
        throw new Error("Failed to fetch world state");
      }
      const json = await res.json();
      return WorldState.fromJson(json);
    },
    {
      staleTime: 10 * 60 * 1000,
      refetchInterval: (worldState) => {
        if (worldState?.type.oneofKind === "starting") {
          return Math.max(
            Number(worldState.startTime) * 1e3 - Date.now() + 500,
            5_000
          );
        }
        return false;
      },
    }
  );

  return queryResult;
};
