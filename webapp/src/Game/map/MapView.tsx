import React, { useEffect, useState } from "react";
import { useAsync, useMedia } from "react-use";
import { classnames } from "tailwindcss-classnames";
import { WorldEntry, WorldZone } from "../../generated/world";
import HexGrid from "./HexGrid";
import { ReactComponent as HeartsSvg } from "../../../images/hearts.svg";
import { useStore } from "../../store";
import { getIdx } from "../getIdx";
import { useNavigate } from "react-router-dom";
import styles from "../../styles";

function unique<T extends unknown>(arr: T[]): T[] {
  return [...new Set(arr)];
}

function MapView() {
  const isMinLg = useMedia("(min-width: 1024px)", false);
  const townX = useStore((state) => state.gameState.townX);
  const townY = useStore((state) => state.gameState.townY);
  const [{ x, y }, setXY] = useState({ x: -1, y: -1 });
  const [size, setSize] = useState(isMinLg ? 9 : 5);
  const navigate = useNavigate();

  useEffect(() => {
    setXY({ x: townX, y: townY });
  }, [townX, townY]);

  useEffect(() => {
    setSize(isMinLg ? 9 : 5);
  }, [isMinLg]);

  // TODO: store zones in store
  const [zones, setZones] = useState<WorldZone[]>([]);

  const [selectedEntry, setSelectedEntry] = useState<{
    entry: WorldEntry;
    myTown: boolean;
    x: number;
    y: number;
  } | null>(null);

  const selected = useAsync(async () => {
    if (selectedEntry === null) {
      return null;
    }

    if (selectedEntry.myTown) {
      return { myTown: true };
    }

    const town =
      selectedEntry.entry.object.oneofKind === "town" &&
      selectedEntry.entry.object.town;
    if (!town) {
      return null;
    }

    const response = await fetch(`/api/user/${town.userId}`);
    if (
      !response.ok ||
      response.headers.get("Content-Type") !== "application/json"
    ) {
      throw new Error("Failed to get user");
    }

    const data = await response.json();
    return { username: data.username as string, myTown: false };
  }, [selectedEntry]);

  useAsync(async () => {
    const missingZones =
      x >= 0 &&
      y >= 0 &&
      unique(
        [
          [x - 5, y - 5],
          [x - 5, y + 5],
          [x + 5, y - 5],
          [x + 5, y + 5],
        ]
          .filter(([a, b]) => a > 0 && b > 0)
          .map(([a, b]) => getIdx(a, b).zidx)
          .filter((idx) => zones[idx] === undefined)
      );

    if (!missingZones || missingZones.length === 0) {
      return;
    }

    for (let idx of missingZones) {
      const response = await fetch(`/api/world/zone?idx=${idx}`);
      if (
        !response.ok ||
        response.headers.get("Content-Type") !== "application/json"
      ) {
        throw new Error("Failed to get zone");
      }
      const zone = WorldZone.fromJson(await response.json());
      setZones((s) => {
        const newArr = [...s];
        newArr[idx] = zone;
        return newArr;
      });
    }
  }, [x, y]);

  const onNavigate = (x: number, y: number) => {
    setXY((p) => ({ x: p.x + x, y: p.y + y }));
  };

  const onClick = (x: number, y: number, zidx: number, eidx: number) => {
    if (zones[zidx].entries[eidx]?.object?.oneofKind === "town") {
      const entry = zones[zidx].entries[eidx];
      setSelectedEntry({ entry, x, y, myTown: x === townX && y === townY });
    } else {
      setSelectedEntry(null);
    }
  };

  const onClickVisit = () => {
    if (selectedEntry) {
      navigate(`/world/entry?x=${selectedEntry.x}&y=${selectedEntry.y}`);
    }
  };

  return (
    <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
      <h2>Map</h2>
      {zones.length === 0 && (
        <div className={classnames("flex", "items-center")}>
          <HeartsSvg />
        </div>
      )}
      {zones.length > 0 && (
        <div
          className={classnames("relative", "flex", "flex-col", "items-center")}
        >
          <div
            className={classnames(
              "flex",
              "justify-center",
              "lg:flex-col",
              "gap-1",
              "lg:absolute",
              "-top-5",
              "-right-1/4",
              "lg:bg-green-200",
              "lg:p-4"
            )}
          >
            <button
              className={styles.primaryButton}
              onClick={() =>
                setSize((size) => Math.min(Math.max(size - 1, 4), 12))
              }
            >
              Zoom In
            </button>
            <button
              className={styles.primaryButton}
              onClick={() =>
                setSize((size) => Math.min(Math.max(size + 1, 4), 12))
              }
            >
              Zoom Out
            </button>
          </div>
          <HexGrid
            x={x}
            y={y}
            size={size}
            data={zones}
            onNavigate={onNavigate}
            onClick={onClick}
          />
          {selected.value && (
            <div
              className={classnames(
                "p-2",
                "bg-green-200",
                "w-full",
                "max-w-sm",
                "flex",
                "justify-between",
                "items-center",
                "sm:-mt-5",
                "lg:absolute",
                "lg:-top-5",
                "lg:-left-1/4",
                "lg:max-w-none",
                "lg:w-auto",
                "lg:p-4",
                "lg:mt-0"
              )}
            >
              {selected.value?.myTown && "Your town"}
              {selected.value?.username && (
                <>Player: {selected.value.username}</>
              )}
              <button
                className={classnames(styles.primaryButton, "ml-4")}
                onClick={onClickVisit}
              >
                Visit
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default MapView;
