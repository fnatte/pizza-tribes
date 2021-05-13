import React, { useState } from "react";
import { classnames } from "tailwindcss-classnames";
import { useStore } from "../store";
import ConstructionQueueTable from "./ConstructionQueueTable";
import PopulationTable from "./PopulationTable";
import TravelQueueTable from "./TravelQueueTable";

const TownExpandMenu: React.VFC<{}> = () => {
  const [minimized, setMinimized] = useState(true);

  const [selection, setSelection] = useState<
    "population" | "constructions" | "travels" | "none"
  >("none");

  const onToggleClick = () => {
    setMinimized((val) => !val);
    setSelection("none");
  };

  const constructionQueueLength = useStore(
    (state) => state.gameState.constructionQueue.length
  );

  const travelQueueLength = useStore(
    (state) => state.gameState.travelQueue.length
  );

  return (
    <div className={classnames("bg-white", "p-2", "pointer-events-auto")}>
      <div className={classnames("flex", "items-center", "justify-end")}>
        {!minimized && selection !== "none" && (
          <button
            className={classnames(
              "border",
              "p-2",
              "flex",
              "justify-center",
              "items-center",
              "mr-auto"
            )}
            onClick={() => setSelection("none")}
          >
            Back
          </button>
        )}
        <button
          className={classnames(
            "border",
            "p-2",
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
      {!minimized && selection === "none" && (
        <ul className={classnames("bg-white", "p-4")}>
          <li className={classnames("border-b", "my-1")}>
            <button
              className={classnames("w-full", "p-2")}
              onClick={() => setSelection("population")}
            >
              Population &gt;
            </button>
          </li>
          <li className={classnames("border-b", "my-1")}>
            <button
              className={classnames("w-full", "p-2")}
              onClick={() => setSelection("constructions")}
            >
              Constructions
              {constructionQueueLength > 0 && (
                <strong> ({constructionQueueLength}) </strong>
              )}{" "}
              &gt;
            </button>
          </li>
          <li className={classnames("border-b", "my-1")}>
            <button
              className={classnames("w-full", "p-2")}
              onClick={() => setSelection("travels")}
            >
              Travels
              {travelQueueLength > 0 && (
                <strong> ({travelQueueLength}) </strong>
              )}{" "}
              &gt;
            </button>
          </li>
        </ul>
      )}
      {!minimized && selection === "population" && (
        <div className={classnames("p-4")}>
          <h4>Population</h4>
          <PopulationTable />
        </div>
      )}
      {!minimized && selection === "constructions" && (
        <div className={classnames("p-4")}>
          <h4>Constructions</h4>
          <ConstructionQueueTable />
        </div>
      )}
      {!minimized && selection === "travels" && (
        <div className={classnames("p-4")}>
          <h4>Travels</h4>
          <TravelQueueTable />
        </div>
      )}
    </div>
  );
};

export default TownExpandMenu;
