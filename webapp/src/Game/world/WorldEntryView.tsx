import React, { useState } from "react";
import { classnames } from "tailwindcss-classnames";
import { Navigate, useLocation } from "react-router-dom";
import { useAsync } from "react-use";
import { ReactComponent as HeartsSvg } from "../../../images/hearts.svg";
import WorldTownView from "./WorldTownView";
import { useStore } from "../../store";
import WorldMyTownView from "./WorldMyTownView";
import { EntriesResponse, WorldEntry } from "../../generated/world";

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

function WorldEntryView() {
  const query = useQuery();
  const x = parseInt(query.get("x") ?? "");
  const y = parseInt(query.get("y") ?? "");

  const townX = useStore((state) => state.gameState.townX);
  const townY = useStore((state) => state.gameState.townY);

  // TODO: store zones in store
  const [data, setData] = useState<{ entry: WorldEntry|null, x: number, y: number} | null>(null);

  useAsync(async () => {
    if (isNaN(x) || isNaN(y)) {
      setData(null);
      return;
    }

    const response = await fetch(`/api/world/entries?x=${x}&y=${y}`);
    if (
      !response.ok ||
      response.headers.get("Content-Type") !== "application/json"
    ) {
      throw new Error("Failed to get zone");
    }
    const resp = EntriesResponse.fromJson(await response.json());
    setData({ entry: resp.entries[`${x}:${y}`] ?? null, x, y });
  }, [x, y]);

  const town =
    (data?.entry?.object.oneofKind === "town" && data.entry.object.town) || null;

  if (isNaN(x) || isNaN(y)) {
    return <div>Bad coordinates</div>;
  }

  return (
    <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
      {data === null || data.x !== x || data.y !== y ? (
        <HeartsSvg />
      ) : (
        <div>
          {townX === x && townY === y ? (
            <WorldMyTownView />
          ) : town !== null ? (
            <WorldTownView town={town} x={x} y={y} />
          ) : (
            <Navigate to="/" />
          )}
        </div>
      )}
    </div>
  );
}

export default WorldEntryView;
