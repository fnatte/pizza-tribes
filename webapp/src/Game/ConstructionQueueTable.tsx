import React from "react";
import classnames from "classnames";
import { useStore } from "../store";
import { CountDown } from "./CountDown";

const ConstructionQueue: React.FC<{
  className?: string;
}> = ({ className }) => {
  const buildings = useStore((state) => state.gameData?.buildings) ?? {};
  const constructionQueue = useStore(
    (state) => state.gameState.constructionQueue
  );

  if (constructionQueue.length === 0) {
    return (
      <div className={classnames("text-sm", "text-gray-800")}>
        There are no ongoing constructions.
      </div>
    );
  }

  return (
    <table className={className} data-cy="construction-queue-table">
      <tbody>
        {constructionQueue.map((construction) => (
          <tr key={construction.lotId}>
            <td
              className={classnames({
                "p-2": true,
                "font-bold": construction.razing,
                "text-red-700": construction.razing,
              })}
            >
              {construction.razing && "Razing "}
              {buildings[construction.building].title}
              {construction.level > 0 && (
                <span>
                  {" "}
                  {!construction.razing && "to"} level {construction.level + 1}
                </span>
              )}
            </td>
            <td className={classnames("p-2")}>
              <CountDown time={construction.completeAt} />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default ConstructionQueue;
