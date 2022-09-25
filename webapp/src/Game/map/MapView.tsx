import React, { useEffect, useState } from "react";
import { useAsync } from "react-use";
import classnames from "classnames";
import { EntriesResponse, WorldEntry } from "../../generated/world";
import HexGridCanvas from "./HexGridCanvas";
import { ReactComponent as HeartsSvg } from "../../../images/hearts.svg";
import { useStore } from "../../store";
import styles from "../../styles";
import { apiFetch } from "../../api";
import { useGameNavigate } from "../useGameNavigate";

function unique<T extends unknown>(arr: T[]): T[] {
  return [...new Set(arr)];
}

const getMapKey = (x: number, y: number) => `${x}:${y}`;

function MapView() {
  const townX = useStore((state) => state.gameState.townX);
  const townY = useStore((state) => state.gameState.townY);
  const [{ x, y }, setXY] = useState({ x: townX, y: townY });
  const [{ x: loadX, y: loadY }, setLoadXY] = useState({ x: townX, y: townY });
  const [zoom, setZoom] = useState(1);
  const navigate = useGameNavigate();

  useEffect(() => {
    setXY({ x: townX, y: townY });
    setLoadXY({ x: townX, y: townY });
  }, [townX, townY]);

  const [map, setMap] = useState(new Map<string, WorldEntry | null>());
  const [selectedXY, setSelectedXY] = useState<{ x: number; y: number }>();

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
    const x = loadX;
    const y = loadY;
    const fetchSize = Math.floor((1 / zoom) * 20);
    const cornerDistance = Math.floor(fetchSize / 2);
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
  }, [loadX, loadY, zoom]);

  const onClick = (x: number, y: number) => {
    setSelectedXY({ x, y });
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

  const onOffsetChange = (x: number, y: number) => {
    setLoadXY({ x, y });
  };

  return (
    <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
      <h2 className="hidden lg:block">Map</h2>
      {map.size === 0 && (
        <div className={classnames("flex", "items-center")}>
          <HeartsSvg />
        </div>
      )}
      {map.size > 0 && (
        <div
          className={classnames(
            "relative",
            "flex",
            "flex-col",
            "items-center",
            "w-full",
            "max-w-screen-lg"
          )}
        >
          <div
            className={classnames(
              "flex",
              "justify-center",
              "gap-1",
              "2xl:absolute",
              "2xl:-top-5",
              "2xl:-right-1/4",
              "2xl:flex-col",
              "2xl:bg-green-200",
              "2xl:p-4"
            )}
          >
            <button
              className={styles.primaryButton}
              onClick={() =>
                setZoom((zoom) => Math.min((zoom * 6) / 5, Math.pow(6 / 5, 3)))
              }
            >
              Zoom In
            </button>
            <button
              className={styles.primaryButton}
              onClick={() =>
                setZoom((zoom) => Math.max((zoom * 5) / 6, Math.pow(5 / 6, 4)))
              }
            >
              Zoom Out
            </button>
          </div>

          <HexGridCanvas
            x={x}
            y={y}
            zoom={zoom}
            data={map}
            onClick={onClick}
            onOffsetChange={onOffsetChange}
            className="container h-[500px] max-h-[48vh] xxs:max-h-[60vh] my-3 xl:my-10 border-4 border-green-700"
            selection={selectedXY}
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
                "sm:-mt-5"
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
