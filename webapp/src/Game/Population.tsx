import React, { useEffect, useState } from "react";
import { useMedia } from "react-use";
import { classnames, TArg, TClasses } from "tailwindcss-classnames";
import { Building } from "../generated/building";
import { Lot, useStore } from "../store";
import { roles } from "./data";

const countBuildings = (
  lots: Record<string, Lot | undefined>
): Record<Building, number> => {
  const counts = {
    0: 0,
    1: 0,
    2: 0,
    3: 0,
  };

  Object.keys(lots).forEach((lotId) => {
    const lot = lots[lotId];
    if (lot) {
      counts[lot.building]++;
    }
  });

  return counts;
};

const Population: React.FC<{ className?: string }> = ({ className }) => {
  const isMinLg = useMedia("(min-width: 1024px)", false);
  const population = useStore((state) => state.gameState.population);
  const lots = useStore((state) => state.gameState.lots);
  const gameData = useStore((state) => state.gameData);
  const [minimized, setMinimized] = useState(isMinLg);

  const onToggleClick = (e: React.MouseEvent) => {
    e.preventDefault();
    setMinimized((value) => !value);
  };

  useEffect(() => setMinimized(!isMinLg), [isMinLg, setMinimized]);

  const buildingCount = countBuildings(lots);

  return (
    <div className={classnames("bg-white", "p-2", className as TArg)}>
      <div className={classnames("flex", "items-center", "justify-between")}>
        <h4
          className={classnames({
            hidden: minimized,
            ["xs:inline" as TClasses]: minimized,
            "mr-2": true,
          })}
        >
          Population
        </h4>
        <div>
          <button
            className={classnames(
              "p-2",
              "border",
              "w-8",
              "h-8",
              "flex",
              "justify-center",
              "items-center"
            )}
            onClick={onToggleClick}
          >
            {minimized ? "➕" : "➖"}
          </button>
        </div>
      </div>
      {!minimized && (
        <table>
          <tbody>
            {roles.map((role) => (
              <tr key={role.id}>
                <td className={classnames("p-2")}>{role.titlePlural}</td>
                <td className={classnames("p-2")}>
                  {role.employer !== undefined
                    ? `${population[role.id].toString()} / ${
                        (gameData?.buildingInfos[role.employer].employer
                          ?.maxWorkforce ?? 0) * buildingCount[role.employer]
                      }`
                    : population[role.id].toString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default Population;
