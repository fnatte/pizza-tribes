import React, { useState } from "react";
import { classnames } from "tailwindcss-classnames";
import { Navigate, useLocation } from "react-router-dom";
import { WorldZone } from "../../generated/world";
import { getIdx } from "../getIdx";
import { useAsync } from "react-use";
import { ReactComponent as HeartsSvg } from "../../../images/hearts.svg";
import WorldTownView from "./WorldTownView";
import { useStore } from "../../store";
import WorldMyTownView from "./WorldMyTownView";

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

function WorldEntryView() {
  const query = useQuery();
  const x = parseInt(query.get("x") ?? "");
  const y = parseInt(query.get("y") ?? "");
  const { zidx, eidx } = getIdx(x, y);

  const townX = useStore((state) => state.gameState.townX);
  const townY = useStore((state) => state.gameState.townY);

  // TODO: store zones in store
  const [zoneData, setZoneData] = useState<{
    zone: WorldZone;
    zidx: number;
  } | null>(null);

  useAsync(async () => {
    if (isNaN(x) || isNaN(y)) {
      setZoneData(null);
      return;
    }

    const response = await fetch(`/api/world/zone?idx=${zidx}`);
    if (
      !response.ok ||
      response.headers.get("Content-Type") !== "application/json"
    ) {
      throw new Error("Failed to get zone");
    }
    const zone = WorldZone.fromJson(await response.json());
    setZoneData({ zone, zidx });
  }, [x, y]);

  const entry = zoneData?.zone.entries[eidx];
  const town =
    (entry?.object.oneofKind === "town" && entry.object.town) || null;

  if (isNaN(x) || isNaN(y)) {
    return <div>Bad coordinates</div>;
  }

  return (
    <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
      {zoneData === null || zoneData.zidx !== zidx ? (
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
