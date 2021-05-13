import JSBI from "jsbi";
import React, { useState } from "react";
import { useInterval } from "react-use";
import { classnames, TArg, TClasses } from "tailwindcss-classnames";
import { useStore } from "../store";

const TravelQueue: React.FC<{ className?: string, minimized: boolean, onToggleClick: () => void }> = ({ className, minimized, onToggleClick }) => {
  const travelQueue = useStore((state) => state.gameState.travelQueue);
  const [now, setNow] = useState(Date.now());
  useInterval(() => {
    setNow(Date.now());
  }, 1000);

  /*
  if (travelQueue.length === 0) {
    return null;
  }
   */

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
        <table>
          <tbody>
            {travelQueue.map((travel) => (
              <tr
                key={
                  travel.arrivalAt.toString() + travel.coins + travel.thieves
                }
              >
                <td className={classnames("p-2")}>
                  {travel.thieves}{" "}
                  {travel.returning ? "returning" : "travelling"} thieves
                </td>
                <td className={classnames("p-2")}>
                  {`in ${JSBI.divide(
                    JSBI.subtract(
                      JSBI.BigInt(travel.arrivalAt),
                      JSBI.multiply(JSBI.BigInt(now), JSBI.BigInt(1e6))
                    ),
                    JSBI.BigInt(1e9)
                  )}s`}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default TravelQueue;
