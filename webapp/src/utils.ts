import { Building } from "./generated/building";
import { GameState, Lot } from "./store";

export type RemoveIndex<T> = {
  [P in keyof T as string extends P
    ? never
    : number extends P
    ? never
    : P]: T[P];
};

export const countBuildings = (
  lots: Record<string, Lot | undefined>
): Record<Building, number> => {
  const counts = {
    0: 0,
    1: 0,
    2: 0,
    3: 0,
  };

  Object.keys(lots).forEach((lotId) => {
    const lot = lots[lotId];
    if (lot) {
      counts[lot.building]++;
    }
  });

  return counts;
};

export const countBuildingsUnderConstruction = (
  constructionQueue: GameState["constructionQueue"]
): Record<Building, number> => {
  return constructionQueue.reduce(
    (counts, item) => {
      counts[item.building]++;
      return counts;
    },
    {
      0: 0,
      1: 0,
      2: 0,
      3: 0,
    }
  );
};
