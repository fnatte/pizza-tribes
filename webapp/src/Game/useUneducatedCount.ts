import { useMemo } from "react";
import { useStore } from "../store";

export function useUneducatedCount() {
  const mice = useStore((state) => state.gameState.mice);
  const counts = useMemo(() => {
    return Object.values(mice).reduce((n, mouse) => {
      if (!mouse.isEducated) {
        n++;
      }
      return n;
    }, 0);
  }, [mice]);

  return counts;
}
