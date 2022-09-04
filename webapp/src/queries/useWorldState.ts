import { useQuery } from "react-query";
import { apiFetch } from "../api";
import { WorldState } from "../generated/world";
import { checkError } from "./error";
import { retry } from "./retry";

export const useWorldState = () => {
  const queryResult = useQuery(
    "worldState",
    async (): Promise<WorldState> => {
      const res = await apiFetch("/world/");
      checkError(res)
      const json = await res.json();
      return WorldState.fromJson(json);
    },
    {
      retry,
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
