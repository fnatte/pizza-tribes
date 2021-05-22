import React, { useEffect, useState } from "react";
import { useAsync, useMedia } from "react-use";
import { classnames } from "tailwindcss-classnames";
import { WorldZone } from "../../generated/world";
import HexGrid from "./HexGrid";
import { ReactComponent as HeartsSvg } from "../../../images/hearts.svg";
import { useStore } from "../../store";
import { getIdx } from "../getIdx";
import { useNavigate } from "react-router-dom";
import styles from "../../styles";

const unique = <T extends unknown>(arr: T[]): T[] => {
  return [...new Set(arr)];
};

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
      navigate(`/world/entry?x=${x}&y=${y}`);
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
        <div className={classnames("relative", "flex", "flex-col", "items-center")}>
          <div
            className={classnames(
              "flex",
              "justify-center",
              "lg:flex-col",
              "gap-1",
              "lg:absolute",
              "-top-5",
              "-right-1/4",
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
        </div>
      )}
    </div>
  );
}

export default MapView;
