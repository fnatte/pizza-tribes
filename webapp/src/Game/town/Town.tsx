import React, { useMemo, useState } from "react";
import classnames from "classnames";
import townImg from "images/v2/town.png";
import { useStore } from "../../store";
import { Building } from "../../generated/building";
import { getTapInfo } from "../../utils";
import { useInterval } from "react-use";
import { useGameNavigate } from "../useGameNavigate";
import { GameState_Lot } from "../../generated/gamestate";
import BuildingImage from "../components/BuildingImage";

const lotData: Record<string, { centerX: number; centerY: number }> = {
  "1": {
    centerX: 746 / 2048,
    centerY: 444 / 2048,
  },
  "2": {
    centerX: 1312 / 2048,
    centerY: 442 / 2048,
  },
  "3": {
    centerX: 203 / 1024,
    centerY: 346 / 1024,
  },
  "4": {
    centerX: 500 / 1024,
    centerY: 372 / 1024,
  },
  "5": {
    centerX: 1610 / 2048,
    centerY: 720 / 2048,
  },
  "6": {
    centerX: 742 / 2048,
    centerY: 1056 / 2048,
  },
  "7": {
    centerX: 1384 / 2048,
    centerY: 1000 / 2048,
  },
  "8": {
    centerX: 568 / 1024,
    centerY: 604 / 1024,
  },
  "9": {
    centerX: 1656 / 2048,
    centerY: 1352 / 2048,
  },
  "10": {
    centerX: 316 / 1024,
    centerY: 759 / 1024,
  },
  "11": {
    centerX: 1278 / 2048,
    centerY: 1606 / 2048,
  },
};

type NotificationPosition = Partial<
  Record<"top" | "left" | "right" | "bottom", number | string>
>;

type BuildingRenderData = {
  centerX: number;
  centerY: number;
  notifications?: Partial<Record<"topright" | "topleft", NotificationPosition>>;
};

const buildingData: Record<Building, BuildingRenderData> = {
  [Building.HOUSE]: {
    centerX: 0.45,
    centerY: 0.7,
    notifications: {
      topright: {
        top: "2%",
        right: "15%",
      },
    },
  },
  [Building.KITCHEN]: {
    centerX: 210 / 512,
    centerY: 350 / 536,
    notifications: {
      topleft: {
        top: "8%",
        left: "20%",
      },
      topright: {
        top: "14.5%",
        right: "14.5%",
      },
    },
  },
  [Building.SHOP]: {
    centerX: 249 / 512,
    centerY: 483 / 636,
    notifications: {
      topleft: {
        top: "3%",
        left: "16%",
      },
      topright: {
        top: "0%",
        right: "15%",
      },
    },
  },
  [Building.SCHOOL]: {
    centerX: 0.4,
    centerY: 0.75,
    notifications: {
      topleft: {
        top: "5%",
        left: "5%",
      },
      topright: {
        top: "2%",
        right: "9%",
      },
    },
  },
  [Building.TOWN_CENTRE]: {
    centerX: 800 / 2048,
    centerY: 1355 / 2048,
    notifications: {
      topright: {
        top: "5%",
        right: "5%",
      },
    },
  },
  [Building.MARKETINGHQ]: {
    centerX: 800 / 2048,
    centerY: 1500 / 2048,
    notifications: {
      topleft: {
        top: "6%",
        left: "6%",
      },
      topright: {
        top: "3%",
        right: "7%",
      },
    },
  },
  [Building.RESEARCH_INSTITUTE]: {
    centerX: 916 / 2048,
    centerY: 1280 / 2048,
    notifications: {
      topleft: {
        top: "6%",
        left: "6%",
      },
      topright: {
        top: "7%",
        right: "15%",
      },
    },
  },
};

function Badge({
  position,
  animation,
  children,
  background = "red",
  size = "normal",
  ...rest
}: {
  position: NotificationPosition;
  animation?: "bounce" | undefined;
  background?: "red" | "white" | undefined;
  size?: "normal" | "big";
} & React.HTMLAttributes<HTMLDivElement>) {
  const offset = useMemo(() => Math.random(), []);

  return (
    <div
      className={classnames(
        "absolute rounded-full flex justify-center items-center",
        {
          "w-3 h-3 md:w-6 md:h-6": size === "normal",
          "w-5 h-5 md:w-8 md:h-8": size === "big",
          "text-xs md:text-md": size === "normal",
          "text-sm md:text-xl": size === "big",
          "bg-red-700": background === "red",
          "bg-gray-300": background === "white",
          "text-gray-50": background === "red",
          "text-black": background === "white",
          "animate-bounce-loop": animation === "bounce",
        }
      )}
      style={{ animationDelay: `${offset}s`, ...position }}
      {...rest}
    >
      {children}
    </div>
  );
}

function NotificationBadge({ position }: { position: NotificationPosition }) {
  return <Badge position={position} animation="bounce" />;
}

function LevelBadge({
  level,
  position,
}: {
  level: number;
  position: NotificationPosition;
}) {
  return (
    <Badge
      position={position}
      background="white"
      size="big"
      data-cy="level-badge"
    >
      {level + 1}
    </Badge>
  );
}

function TownLot({ id }: { id: string }) {
  const buildingInfos = useStore((state) => state.gameData?.buildings);
  const lots = useStore((state) => state.gameState.lots);
  const discoveries = useStore((state) => state.gameState.discoveries);
  const constructionQueue = useStore(
    (state) => state.gameState.constructionQueue
  );
  const gameNavigate = useGameNavigate();
  const [now, setNow] = useState(new Date());
  useInterval(() => setNow(new Date()), 10_000);

  if (!lotData[id]) {
    console.warn("Could not find lot data for id", id);
    return null;
  }

  const lot = lots[id] as GameState_Lot | undefined;
  const { centerX, centerY } = lotData[id];

  const construction = constructionQueue.find((x) => x.lotId === id);
  const isUnderConstruction =
    construction !== undefined &&
    !construction.razing &&
    construction.level <= 0;

  const level = lot?.level;
  const building = lot?.building ?? construction?.building;
  const buildingInfo =
    buildingInfos && building !== undefined ? buildingInfos[building] : null;

  const showNotification =
    lot !== undefined && getTapInfo(lot, discoveries, now).canTap;
  const buildingRenderData =
    building !== undefined ? buildingData[building] : undefined;

  return (
    <div
      className="absolute top-0 left-0 inline-block w-[20%] h-[16%] cursor-pointer"
      style={{
        left: `${centerX * 100}%`,
        top: `${centerY * 100}%`,
        transform: "translate(-50%, -50%)",
      }}
      title={buildingInfo?.title}
      data-id={id}
      data-cy={`lot${id}`}
      onClick={() => gameNavigate("town-lot", id)}
    >
      {building !== undefined && (
        <>
          <BuildingImage
            building={building}
            className={classnames(
              "absolute top-1/2 left-1/2 w-[65%] h-[65%] object-contain",
              {
                "opacity-50": isUnderConstruction,
              }
            )}
            style={{
              transform: `translate(${
                -(buildingRenderData?.centerX ?? 0.5) * 100
              }%, ${-(buildingRenderData?.centerY ?? 0.5) * 100}%)`,
            }}
          />
          {showNotification && (
            <NotificationBadge
              position={
                buildingRenderData?.notifications?.topleft ?? {
                  top: 0,
                  left: 0,
                }
              }
            />
          )}
          {level !== undefined && (
            <LevelBadge
              level={level}
              position={
                buildingRenderData?.notifications?.topright ?? {
                  top: 0,
                  right: 0,
                }
              }
            />
          )}
        </>
      )}
    </div>
  );
}

function Town() {
  return (
    <div className="relative">
      <img src={townImg} />
      <TownLot id="1" />
      <TownLot id="2" />
      <TownLot id="3" />
      <TownLot id="4" />
      <TownLot id="5" />
      <TownLot id="6" />
      <TownLot id="7" />
      <TownLot id="8" />
      <TownLot id="9" />
      <TownLot id="10" />
      <TownLot id="11" />
    </div>
  );
}

export default Town;
