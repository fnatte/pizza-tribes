import React from "react";
import { classnames } from "tailwindcss-classnames";
import { Education } from "../generated/education";
import { GameState_Population } from "../generated/gamestate";
import { useStore } from "../store";
import { countMaxEmployed } from "../utils";

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

const Population: React.FC<{ className?: string }> = ({ className }) => {
  const population = useStore((state) => state.gameState.population);
  const lots = useStore((state) => state.gameState.lots);
  const gameData = useStore((state) => state.gameData);

  const educations = gameData?.educations ?? {};
  const maxEmployed = gameData && countMaxEmployed(lots, gameData);

  return (
    <table className={className}>
      <tbody>
        <tr>
          <td className={classnames("p-2")}>Uneducated</td>
          <td className={classnames("p-2")}>{population["uneducated"]}</td>
        </tr>
        {Object.keys(educations)
          .map(Number)
          .map((id) => {
            const popKey = getPopulationKey(id);
            const pop = (popKey && population[popKey]) ?? 0;
            const education = educations[id];
            const max =
              (maxEmployed &&
                education.employer !== undefined &&
                maxEmployed[education.employer]) ||
              0;
            const isOverStaffed = pop > max;

            if (pop === 0 && max === 0) {
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
      </tbody>
    </table>
  );
};

export default Population;
