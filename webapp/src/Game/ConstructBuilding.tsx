import React from "react";
import { useNavigate } from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import { Building } from "../generated/building";
import { useStore } from "../store";
import styles from "../styles";
import {
  countBuildings,
  countBuildingsUnderConstruction,
  isNotNull,
} from "../utils";
import PlaceholderImage from "./PlaceholderImage";

const title = classnames("text-lg", "md:text-xl", "mb-2");
const label = classnames("text-xs", "md:text-sm", "mr-1");
const value = classnames("text-sm", "md:text-lg", "ml-1");

type Props = {
  lotId: string;
};

const toBuildingId = (key: string) => {
  const n = Number(key);
  if (n in Building) {
    return n as Building;
  }

  return null;
};

const ConstructBuilding = ({ lotId }: Props) => {
  const constructBuilding = useStore((state) => state.constructBuilding);
  const buildings = useStore((state) => state.gameData?.buildings) ?? [];
  const coins = useStore((state) => state.gameState.resources.coins);
  const lots = useStore((state) => state.gameState.lots);
  const constructionQueue = useStore(
    (state) => state.gameState.constructionQueue
  );
  const navigate = useNavigate();

  const buildingCounts = countBuildings(lots);
  const buildingConstrCounts = countBuildingsUnderConstruction(
    constructionQueue
  );

  const onSelectClick = (e: React.MouseEvent, building: Building) => {
    e.preventDefault();
    constructBuilding(lotId, building);
    navigate("/town");
  };

  return (
    <div className={classnames("container", "mx-auto", "mt-4", "px-1")}>
      <h2>Construct Building</h2>
      {Object.keys(buildings)
        .map(toBuildingId)
        .filter(isNotNull)
        .map((id) => {
          let discountText: string | null = null;
          let discountCost: number | null = null;
          let reducedTime: number | null = null;

          // First construction of building type is free
          if (buildingCounts[id] + buildingConstrCounts[id] === 0) {
            discountCost = 0;
            discountText = "First one is free and fast!";
            reducedTime = Math.ceil(buildings[id].constructionTime / 100);
          }

          const canAfford = coins >= (discountCost ?? buildings[id].cost);

          return (
            <div className={classnames("flex", "mb-8")} key={id}>
              <PlaceholderImage />
              <div className={classnames("ml-4")}>
                <div className={title}>{buildings[id].title}</div>
                <div
                  className={classnames("grid", "grid-cols-2", "items-center")}
                >
                  <span className={label}>Cost:</span>
                  <span className={value}>
                    {discountCost !== null ? (
                      <>
                        <span
                          className={classnames(
                            "line-through",
                            "mr-2",
                            "text-sm"
                          )}
                        >
                          {buildings[id].cost} coins
                        </span>
                        <span>{discountCost} coins</span>
                      </>
                    ) : (
                      <span>{buildings[id].cost} coins</span>
                    )}
                  </span>
                  <span className={label}>Build time:</span>
                  <span className={value}>
                    {reducedTime !== null ? (
                      <>
                        <span
                          className={classnames(
                            "line-through",
                            "mr-2",
                            "text-sm"
                          )}
                        >
                          {buildings[id].constructionTime}s
                        </span>
                        <span>{reducedTime}s</span>
                      </>
                    ) : (
                      <span>{buildings[id].constructionTime}s</span>
                    )}
                  </span>
                </div>
                <div className={classnames("my-2")}>
                  <button
                    className={classnames(styles.button)}
                    onClick={(e) => onSelectClick(e, id)}
                    disabled={!canAfford}
                  >
                    Place Building
                  </button>
                  {discountText && (
                    <span className={classnames("mx-4", "text-sm")}>
                      {discountText}
                    </span>
                  )}
                  {!canAfford && (
                    <span className={classnames("mx-4", "text-sm")}>
                      Not enough coins
                    </span>
                  )}
                </div>
              </div>
            </div>
          );
        })}
    </div>
  );
};

export default ConstructBuilding;
