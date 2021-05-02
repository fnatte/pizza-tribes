import React, { useCallback } from "react";
import { useParams } from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import {Building} from "../generated/building";
import { useStore } from "../store";
import ConstructBuilding from "./ConstructBuilding";
import School from "./School";

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
        "mt-2",
      )}
    >
      {!lot && <ConstructBuilding lotId={id} />}
      {lot?.building === Building.KITCHEN && <h2>Kitchen</h2>}
      {lot?.building === Building.HOUSE && <h2>House</h2>}
      {lot?.building === Building.SHOP && <h2>Shop</h2>}
      {lot?.building === Building.SCHOOL && <School />}
    </div>
  );
}

export default TownLot;
