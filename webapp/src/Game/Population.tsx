import React from "react";
import classnames from "classnames";
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
        className
      )}
    >
      <div className={classnames("flex", "items-center", "justify-between")}>
        <h4
          className={classnames({
            hidden: minimized,
            "xs:inline": minimized,
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
            aria-expanded={!minimized}
            data-cy="population-table-toggle-button"
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
