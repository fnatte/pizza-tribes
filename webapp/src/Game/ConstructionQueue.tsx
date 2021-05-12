import JSBI from "jsbi";
import React, { useEffect, useState } from "react";
import { useInterval, useMedia } from "react-use";
import { classnames, TArg, TClasses } from "tailwindcss-classnames";
import { useStore } from "../store";
import {formatNanoTimestampToNowShort} from "../utils";

const ConstructionQueue: React.FC<{ className?: string }> = ({ className }) => {
  const isMinLg = useMedia("(min-width: 1024px)", false);
  const buildings = useStore((state) => state.gameData?.buildings) ?? {};
  const constructionQueue = useStore(
    (state) => state.gameState.constructionQueue
  );
  const [minimized, setMinimized] = useState(isMinLg);

  const onToggleClick = (e: React.MouseEvent) => {
    e.preventDefault();
    setMinimized((value) => !value);
  };

  useEffect(() => setMinimized(!isMinLg), [isMinLg, setMinimized]);

  const [now, setNow] = useState(Date.now());
  useInterval(() => {
    setNow(Date.now());
  }, 1000);

  if (constructionQueue.length === 0) {
    return <div />;
  }

  return (
    <div
      className={classnames(
        "bg-white",
        "p-2",
        className as TArg,
        "pointer-events-auto"
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
      {!minimized && (
        <table>
          <tbody>
            {constructionQueue.map((construction) => (
              <tr key={construction.lotId}>
                <td className={classnames("p-2")}>
                  {buildings[construction.building].title}
                  {construction.level > 0 && (
                    <span> to level {construction.level + 1}</span>
                  )}
                </td>
                <td className={classnames("p-2")}>
                  {/*`in ${JSBI.divide(
                    JSBI.subtract(
                      JSBI.BigInt(construction.completeAt),
                      JSBI.multiply(JSBI.BigInt(now), JSBI.BigInt(1e6))
                    ),
                    JSBI.BigInt(1e9)
                  )}s`*/}
                  {formatNanoTimestampToNowShort(construction.completeAt)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default ConstructionQueue;
