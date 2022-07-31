import React from "react";
import classnames from "classnames";
import { useStore } from "../store";
import { CountDown } from "./CountDown";

const TravelQueue: React.FC<{ className?: string }> = ({ className }) => {
  const travelQueue = useStore((state) => state.gameState.travelQueue);
  if (travelQueue.length === 0) {
    return (
      <div className={classnames("text-sm", "text-gray-800")}>
        There are no ongoing travels.
      </div>
    );
  }

  return (
    <table className={className}>
      <tbody>
        {travelQueue.map((travel) => (
          <tr
            key={travel.arrivalAt.toString() + travel.coins + travel.thieves}
            data-cy="travel-queue-row"
          >
            <td className={classnames("p-2")}>
              {travel.thieves} {travel.returning ? "returning" : "travelling"}{" "}
              {travel.thieves > 1 ? 'thieves' : 'thief'}
            </td>
            <td className={classnames("p-2")}>
              <CountDown time={travel.arrivalAt} />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default TravelQueue;
