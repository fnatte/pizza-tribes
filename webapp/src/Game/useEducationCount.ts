import { useMemo } from "react";
import { Education } from "../generated/education";
import { useStore } from "../store";

export function useEducationCount() {
  const mice = useStore((state) => state.gameState.mice);
  const counts = useMemo(() => {
    return Object.values(mice).reduce((o, mouse) => {
      if (mouse.isEducated) {
        const val = o[mouse.education];
        if (typeof val !== "undefined") {
          o[mouse.education] = val + 1;
        } else {
          o[mouse.education] = 1;
        }
      }
      return o;
    }, {} as Record<Education, number | undefined>);
  }, [mice]);

  return counts;
}
