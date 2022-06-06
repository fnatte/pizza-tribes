import React from "react";
import classnames from "classnames";
import { useStore } from "../store";
import ConstructionQueueTable from "./ConstructionQueueTable";

const ConstructionQueue: React.FC<{
  className?: string;
  minimized: boolean;
  onToggleClick: () => void;
}> = ({ className, minimized, onToggleClick }) => {
  const constructionQueue = useStore(
    (state) => state.gameState.constructionQueue
  );

  if (constructionQueue.length === 0) {
    return <div />;
  }

  return (
    <div
      className={classnames(
        "bg-white",
        "p-2",
        className,
        "pointer-events-auto"
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
          Construction Queue
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
      {!minimized && <ConstructionQueueTable />}
    </div>
  );
};

export default ConstructionQueue;
