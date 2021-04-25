import React, { useCallback } from "react";
import { useParams } from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import { useStore } from "../store";
import ConstructBuilding from "./ConstructBuilding";

function TownLot() {
  const { id } = useParams();

  const lot = useStore(useCallback((state) => state.gameState.lots[id], [id]));

  return (
    <div
      className={classnames(
        "flex",
        "flex-col",
        "items-center",
        "justify-center",
        "mt-2"
      )}
    >
      {!lot && <ConstructBuilding lotId={id} />}
      {lot?.building === "kitchen" && <h2>Kitchen</h2>}
      {lot?.building === "house" && <h2>House</h2>}
      {lot?.building === "shop" && <h2>Shop</h2>}
    </div>
  );
}

export default TownLot;
