import { formatDuration, intervalToDuration } from "date-fns";
import JSBI from "jsbi";
import { Building } from "./generated/building";
import { GameState_Population } from "./generated/gamestate";
import { GameState, Lot } from "./store";

export type RemoveIndex<T> = {
  [P in keyof T as string extends P
    ? never
    : number extends P
    ? never
    : P]: T[P];
};

export function isNotNull<T>(v: T | null): v is T {
  return v !== null;
}

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

export const countPopulation = (population: GameState_Population): number => {
  return (
    population.chefs +
    population.guards +
    population.thieves +
    population.salesmice +
    population.uneducated
  );
};

const shortDuration = (str: string) => {
  return str
    .replace("hours", "h")
    .replace("hour", "h")
    .replace("minutes", "min")
    .replace("minute", "min")
    .replace("seconds", "sec")
    .replace("second", "sec");
};

export const formatDurationShort = (time: number) =>
  shortDuration(
    formatDuration(
      intervalToDuration({
        start: 0,
        end: time * 1000,
      }),
      { delimiter: ", " }
    )
  );

export const generateId = () => {
  return (
    Array(16)
      .fill(0)
      .map(() => String.fromCharCode(Math.floor(Math.random() * 26) + 97))
      .join("") + Date.now().toString(36)
  );
};

export const parseDateNano = (ns: string) => {
  return new Date(
    JSBI.toNumber(JSBI.divide(JSBI.BigInt(ns), JSBI.BigInt(1e6)))
  );
};
