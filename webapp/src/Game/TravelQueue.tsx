import React from "react";
import { classnames, TArg, TClasses } from "tailwindcss-classnames";
import { useStore } from "../store";
import TravelQueueTable from "./TravelQueueTable";

const TravelQueue: React.FC<{ className?: string, minimized: boolean, onToggleClick: () => void }> = ({ className, minimized, onToggleClick }) => {
  const travelQueue = useStore((state) => state.gameState.travelQueue);

  if (travelQueue.length === 0) {
    return null;
  }

  return (
    <div className={classnames("bg-white", "p-2", className as TArg)}>
      <div
        className={classnames(
          "flex",
          "items-center",
          "justify-between",
          "pointer-events-auto"
        )}
      >
        <h4
          className={classnames({
            hidden: minimized,
            ["xs:inline" as TClasses]: minimized,
            "mr-2": true,
          })}
        >
          Travels
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
        <TravelQueueTable />
      )}
    </div>
  );
};

export default TravelQueue;
