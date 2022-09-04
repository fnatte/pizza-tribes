import React, { useEffect, useState } from "react";
import { useAsync, useMedia } from "react-use";
import classnames from "classnames";
import { EntriesResponse, WorldEntry } from "../../generated/world";
import HexGrid from "./HexGrid";
import { ReactComponent as HeartsSvg } from "../../../images/hearts.svg";
import { useStore } from "../../store";
import { useNavigate } from "react-router-dom";
import styles from "../../styles";
import { apiFetch } from "../../api";
import { useGameNavigate } from "../useGameNavigate";

function unique<T extends unknown>(arr: T[]): T[] {
  return [...new Set(arr)];
}

const getMapKey = (x: number, y: number) => `${x}:${y}`;

function MapView() {
  const isMinLg = useMedia("(min-width: 1024px)", false);
  const townX = useStore((state) => state.gameState.townX);
  const townY = useStore((state) => state.gameState.townY);
  const [{ x, y }, setXY] = useState({ x: townX, y: townY });
  const [size, setSize] = useState(isMinLg ? 9 : 5);
  const navigate = useGameNavigate();

  useEffect(() => {
    setXY({ x: townX, y: townY });
  }, [townX, townY]);

  useEffect(() => {
    setSize(isMinLg ? 9 : 5);
  }, [isMinLg]);

  const [map, setMap] = useState(new Map<string, WorldEntry | null>());

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

    const response = await apiFetch(`/user/${town.userId}`);
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
    const fetchSize = size * 2;
    const cornerDistance = Math.floor(size * 1.2);
    const missingCorners = unique(
      [
        [x - cornerDistance, y - cornerDistance],
        [x - cornerDistance, y + cornerDistance],
        [x + cornerDistance, y - cornerDistance],
        [x + cornerDistance, y + cornerDistance],
      ].filter(([x, y]) => !map.has(getMapKey(x, y)))
    );

    if (!missingCorners || missingCorners.length === 0) {
      return;
    }

    for (let [x, y] of missingCorners) {
      const response = await apiFetch(
        `/world/entries?x=${x}&y=${y}&r=${fetchSize}`
      );
      if (
        !response.ok ||
        response.headers.get("Content-Type") !== "application/json"
      ) {
        throw new Error("Failed to get zone");
      }
      const resp = EntriesResponse.fromJson(await response.json());
      const newMap = new Map();

      // Add all points in circle with radius r
      const points = [];
      const r = fetchSize;
      for (let i = y - r; i < y + r; i++) {
        for (
          let j = x;
          Math.pow(j - x, 2) + Math.pow(i - y, 2) < Math.pow(r, 2);
          j--
        ) {
          points.push([j, i]);
        }
        for (
          let j = x + 1;
          (j - x) * (j - x) + (i - y) * (i - y) < r * r;
          j++
        ) {
          points.push([j, i]);
        }
      }

      for (let [x, y] of points) {
        const entry = resp.entries[getMapKey(x, y)];
        newMap.set(getMapKey(x, y), entry ?? null);
      }

      setMap((m) => new Map([...m, ...newMap]));
    }
  }, [x, y, size]);

  const onNavigate = (x: number, y: number) => {
    setXY((p) => ({ x: p.x + x, y: p.y + y }));
  };

  const onClick = (x: number, y: number) => {
    const entry = map.get(getMapKey(x, y));
    if (entry?.object?.oneofKind === "town") {
      setSelectedEntry({ entry, x, y, myTown: x === townX && y === townY });
    } else {
      setSelectedEntry(null);
    }
  };

  const onClickVisit = () => {
    if (selectedEntry) {
      navigate("world-entry", selectedEntry.x, selectedEntry.y);
    }
  };

  return (
    <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
      <h2>Map</h2>
      {map.size === 0 && (
        <div className={classnames("flex", "items-center")}>
          <HeartsSvg />
        </div>
      )}
      {map.size > 0 && (
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
            data={map}
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
