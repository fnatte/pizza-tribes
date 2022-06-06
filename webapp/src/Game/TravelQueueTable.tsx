import React, { useState } from "react";
import { useInterval } from "react-use";
import classnames from "classnames";
import { useStore } from "../store";
import { formatNanoTimestampToNowShort } from "../utils";

const TravelQueue: React.FC<{ className?: string }> = ({ className }) => {
  const travelQueue = useStore((state) => state.gameState.travelQueue);
  const [now, setNow] = useState(Date.now());
  useInterval(() => {
    setNow(Date.now());
  }, 1000);

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
          <tr key={travel.arrivalAt.toString() + travel.coins + travel.thieves}>
            <td className={classnames("p-2")}>
              {travel.thieves} {travel.returning ? "returning" : "travelling"}{" "}
              thieves
            </td>
            <td className={classnames("p-2")}>
              {formatNanoTimestampToNowShort(travel.arrivalAt)}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default TravelQueue;
