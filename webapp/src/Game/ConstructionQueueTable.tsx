import React, { useState } from "react";
import { useInterval } from "react-use";
import { classnames } from "tailwindcss-classnames";
import { useStore } from "../store";
import { formatNanoTimestampToNowShort } from "../utils";

const ConstructionQueue: React.FC<{
  className?: string;
}> = ({ className }) => {
  const buildings = useStore((state) => state.gameData?.buildings) ?? {};
  const constructionQueue = useStore(
    (state) => state.gameState.constructionQueue
  );
  const [now, setNow] = useState(Date.now());
  useInterval(() => {
    setNow(Date.now());
  }, 1000);

  if (constructionQueue.length === 0) {
    return (
      <div className={classnames("text-sm", "text-gray-800")}>
        There are no ongoing constructions.
      </div>
    );
  }

  return (
    <table className={className}>
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
                <span> {!construction.razing && 'to'} level {construction.level + 1}</span>
              )}
            </td>
            <td className={classnames("p-2")}>
              {formatNanoTimestampToNowShort(construction.completeAt)}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default ConstructionQueue;
