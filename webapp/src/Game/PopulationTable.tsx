import React from "react";
import classnames from "classnames";
import { Education } from "../generated/education";
import { GameState_Population } from "../generated/gamestate";
import { useStore } from "../store";
import { countMaxEmployedByBuilding } from "../utils";
import { sum } from "lodash";

const getPopulationKey = (id: number): keyof GameState_Population | null => {
  switch (id) {
    case Education.CHEF:
      return "chefs";
    case Education.GUARD:
      return "guards";
    case Education.THIEF:
      return "thieves";
    case Education.SALESMOUSE:
      return "salesmice";
    case Education.PUBLICIST:
      return "publicists";
  }
  return null;
};

const PopulationTable: React.FC<{
  className?: string;
  showZeroes?: boolean;
  showTotalCount?: boolean;
  showHeader?: boolean;
}> = ({ className, showZeroes, showTotalCount, showHeader }) => {
  const population = useStore((state) => state.gameState.population);
  const lots = useStore((state) => state.gameState.lots);
  const gameData = useStore((state) => state.gameData);

  const educations = gameData?.educations ?? {};
  const maxEmployed = gameData && countMaxEmployedByBuilding(lots, gameData);

  const populationCount = population ? sum(Object.values(population)) : 0;

  return (
    <table className={className}>
      {showHeader && (
        <thead>
          <tr>
            <th colSpan={2} className="p-2 text-left text-sm text-gray-700">Education</th>
          </tr>
        </thead>
      )}
      <tbody>
        <tr>
          <td className={classnames("p-2")}>Uneducated</td>
          <td className={classnames("p-2")}>{population?.uneducated ?? 0}</td>
        </tr>
        {Object.keys(educations)
          .map(Number)
          .map((id) => {
            const popKey = getPopulationKey(id);
            const pop = (popKey && population?.[popKey]) ?? 0;
            const education = educations[id];
            const max =
              (maxEmployed &&
                education.employer !== undefined &&
                maxEmployed[education.employer]) ||
              0;
            const isOverStaffed = pop > max;

            if (!showZeroes && pop === 0 && max === 0) {
              return null;
            }

            return (
              <tr key={id}>
                <td className={classnames("p-2")}>
                  {educations[id].titlePlural}
                </td>
                <td className={classnames("p-2")}>
                  {education.employer !== undefined && max !== undefined ? (
                    <>
                      {pop.toString()} /{" "}
                      <span
                        className={classnames({
                          "text-red-800": isOverStaffed,
                        })}
                      >
                        {max}
                      </span>
                    </>
                  ) : (
                    pop.toString()
                  )}
                </td>
              </tr>
            );
          })}
        {showTotalCount && (
          <>
            <tr className="h-2 border-b border-gray-600"></tr>
            <tr className="h-2"></tr>
            <tr>
              <td className={classnames("p-2")}>Total Count</td>
              <td className={classnames("p-2")}>{populationCount}</td>
            </tr>
          </>
        )}
      </tbody>
    </table>
  );
};

export default PopulationTable;
