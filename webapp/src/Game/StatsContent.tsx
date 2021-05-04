import React from "react";
import { classnames } from "tailwindcss-classnames";
import { useStore } from "../store";

const StatsContent: React.FC<{}> = () => {
  const stats = useStore((state) => state.gameStats);

  const rows = [
    { label: "Employed chefs", value: stats.employedChefs },
    { label: "Employed salesmice", value: stats.employedSalesmice },
    { label: "Pizzas produces", value: `${stats.pizzasProducedPerSecond}/s` },
    { label: "Max sells", value: `${stats.maxSellsByMicePerSecond}/s` },
    { label: "Pizza demand (offpeak)", value: `${stats.demandOffpeak}/s` },
    { label: "Pizza demand (rush hour)", value: `${stats.demandRushHour}/s` },
  ];

  return (
    <div
      className={classnames(
        "flex",
        "flex-col",
        "items-center",
        "justify-center",
        "mt-2"
      )}
    >
      Stats
      <table
        className={classnames(
          "w-full",
          "max-w-md",
          "my-4",
          "border-collapse",
          "border-green-400",
          "border-2",
        )}
      >
        <tbody>
          {rows.map(({ label, value }, i) => (
            <tr
              key={label}
              className={classnames({
                "bg-green-200": i % 2 === 0,
              })}
            >
              <td className={classnames("p-1")}>{label}</td>
              <td className={classnames("p-1")}>{value}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default StatsContent;
