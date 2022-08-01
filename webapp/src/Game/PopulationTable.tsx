import React from "react";
import classnames from "classnames";
import { Education } from "../generated/education";
import { useStore } from "../store";
import { countMaxEmployedByBuilding } from "../utils";
import { sum } from "lodash";
import { useUneducatedCount } from "./useUneducatedCount";
import { useEducationCount } from "./useEducationCount";

const PopulationTable: React.FC<{
  className?: string;
  showZeroes?: boolean;
  showTotalCount?: boolean;
  showHeader?: boolean;
}> = ({ className, showZeroes, showTotalCount, showHeader }) => {
  const uneducatedCount = useUneducatedCount();
  const educationCount = useEducationCount();
  const lots = useStore((state) => state.gameState.lots);
  const gameData = useStore((state) => state.gameData);

  const educations = gameData?.educations ?? {};
  const maxEmployed = gameData && countMaxEmployedByBuilding(lots, gameData);

  const populationCount = educationCount
    ? sum(Object.values(educationCount))
    : 0;

  return (
    <table className={className} data-cy="population-table">
      {showHeader && (
        <thead>
          <tr>
            <th colSpan={2} className="p-2 text-left text-sm text-gray-700">
              Education
            </th>
          </tr>
        </thead>
      )}
      <tbody>
        <tr>
          <td className={classnames("p-2")}>Uneducated</td>
          <td
            className={classnames("p-2")}
            data-cy="population-table-uneducated-count"
          >
            {uneducatedCount}
          </td>
        </tr>
        {Object.keys(educations)
          .map((id) => Number(id) as Education)
          .map((id) => {
            const pop = educationCount[id] ?? 0;
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
              <tr key={id} data-cy="population-table-row">
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
