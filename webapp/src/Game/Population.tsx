import React from "react";
import { classnames, TArg, TClasses } from "tailwindcss-classnames";
import PopulationTable from "./PopulationTable";

const Population: React.FC<{
  className?: string;
  minimized: boolean;
  onToggleClick: () => void;
}> = ({ className, minimized, onToggleClick }) => {

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
        <PopulationTable />
      )}
    </div>
  );
};

export default Population;
