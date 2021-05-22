import { formatDuration, intervalToDuration } from "date-fns";
import JSBI from "jsbi";
import { Building } from "./generated/building";
import { Education } from "./generated/education";
import { GameState_Population } from "./generated/gamestate";
import { GameData } from "./generated/game_data";
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
    4: 0,
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
      4: 0,
    }
  );
};

export const countPopulation = (population: GameState_Population): number => {
  return (
    population.chefs +
    population.guards +
    population.thieves +
    population.salesmice +
    population.publicists +
    population.uneducated
  );
};

export const countMaxEmployed = (
  lots: GameState["lots"],
  gameData: GameData
): Record<Education, number | undefined> => {
  const counts: Record<Education, number|undefined> = {
    [Building.SHOP]: undefined,
    [Building.HOUSE]: undefined,
    [Building.SCHOOL]: undefined,
    [Building.KITCHEN]: undefined,
    [Building.MARKETINGHQ]: undefined,
  };

  Object.keys(lots).map((lotId) => {
    const lot = lots[lotId];
    const building = lot?.building;
    if (!lot || building === undefined) {
      return;
    }

    const info = gameData.buildings[building];
    const levelInfo = info.levelInfos[lot.level];
    if (levelInfo?.employer !== undefined) {
      counts[lot.building] =
        (counts[lot.building] || 0) + levelInfo.employer.maxWorkforce;
    }
  });

  return counts;
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

export const formatNanoTimestampToNowShort = (time: string) => {
  const now = Date.now();
  const totalSeconds = JSBI.toNumber(
    JSBI.divide(
      JSBI.subtract(
        JSBI.BigInt(time),
        JSBI.multiply(JSBI.BigInt(now), JSBI.BigInt(1e6))
      ),
      JSBI.BigInt(1e9)
    )
  );

  if (totalSeconds <= 0) {
    return 'now';
  }

  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds - hours * 3600) / 60);
  const seconds = Math.floor(totalSeconds - hours * 3600 - minutes * 60);

  if (hours > 0) {
    return `in ${hours} h ${minutes} min`;
  } else if (minutes > 0) {
    return `in ${minutes} min ${seconds} sec`;
  } else {
    return `in ${seconds} sec`;
  }
};

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

const numberFormat = new Intl.NumberFormat();
export const formatNumber = (n: number) => numberFormat.format(n);

export const getTapInfo = (lot: Lot) => {
  if (lot.building !== Building.KITCHEN && lot.building !== Building.SHOP) {
    return { canTap: false, nextTapAt: 0 };
  }

  // convert lot.tappedAt from ns ms
  const tappedAt = JSBI.toNumber(
    JSBI.divide(JSBI.BigInt(lot.tappedAt), JSBI.BigInt(1e6))
  );
  const nextTapAt = tappedAt + 60_000 * 60;
  const canTap = nextTapAt < Date.now();

  return { canTap, nextTapAt };
}

