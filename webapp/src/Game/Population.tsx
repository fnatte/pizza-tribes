import React from "react";
import { classnames, TArg, TClasses } from "tailwindcss-classnames";
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
  }
  return null;
};

const Population: React.FC<{
  className?: string;
  minimized: boolean;
  onToggleClick: () => void;
}> = ({ className, minimized, onToggleClick }) => {
  const population = useStore((state) => state.gameState.population);
  const lots = useStore((state) => state.gameState.lots);
  const gameData = useStore((state) => state.gameData);

  const educations = gameData?.educations ?? {};
  const maxEmployed = gameData && countMaxEmployed(lots, gameData);

  return (
    <div
      className={classnames(
        "bg-white",
        "p-2",
        "pointer-events-auto",
        className as TArg
      )}
    >
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
            <tr>
              <td className={classnames("p-2")}>Uneducated</td>
              <td className={classnames("p-2")}>{population["uneducated"]}</td>
            </tr>
            {Object.keys(educations)
              .map(Number)
              .map((id) => {
                const popKey = getPopulationKey(id);
                const pop = (popKey && population[popKey].toString()) ?? "0";
                const education = educations[id];
                const max =
                  maxEmployed &&
                  education.employer !== undefined &&
                  maxEmployed[education.employer];

                return (
                  <tr key={id}>
                    <td className={classnames("p-2")}>
                      {educations[id].titlePlural}
                    </td>
                    <td className={classnames("p-2")}>
                      {education.employer !== undefined
                        ? `${pop.toString()} / ${max}`
                        : pop.toString()}
                    </td>
                  </tr>
                );
              })}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default Population;
