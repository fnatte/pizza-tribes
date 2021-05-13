import React, { useState } from "react";
import { classnames } from "tailwindcss-classnames";
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
              className={classnames("w-full")}
              onClick={() => setSelection("population")}
            >
              Population &gt;
            </button>
          </li>
          <li className={classnames("border-b", "my-1")}>
            <button
              className={classnames("w-full")}
              onClick={() => setSelection("constructions")}
            >
              Constructions &gt;
            </button>
          </li>
          <li className={classnames("border-b", "my-1")}>
            <button
              className={classnames("w-full")}
              onClick={() => setSelection("travels")}
            >
              Travels &gt;
            </button>
          </li>
        </ul>
      )}
      {!minimized && selection === "population" && (
        <>
          <h4>Population</h4>
          <PopulationTable />
        </>
      )}
      {!minimized && selection === "constructions" && (
        <>
          <h4>Constructions</h4>
          <ConstructionQueueTable />
        </>
      )}
      {!minimized && selection === "travels" && (
        <>
          <h4>Travels</h4>
          <TravelQueueTable />
        </>
      )}
    </div>
  );
};

export default TownExpandMenu;
