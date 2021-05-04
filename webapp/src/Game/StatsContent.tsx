import { format, fromUnixTime } from "date-fns";
import React from "react";
import { useAsync } from "react-use";
import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { classnames } from "tailwindcss-classnames";
import { TimeseriesData } from "../generated/timeseries";
import { useStore } from "../store";

const StatsContent: React.FC<{}> = () => {
  const stats = useStore((state) => state.gameStats);

  const rows =
    (stats && [
      { label: "Employed chefs", value: stats.employedChefs },
      { label: "Employed salesmice", value: stats.employedSalesmice },
      { label: "Pizzas produces", value: `${stats.pizzasProducedPerSecond}/s` },
      { label: "Max sells", value: `${stats.maxSellsByMicePerSecond}/s` },
      { label: "Pizza demand (offpeak)", value: `${stats.demandOffpeak}/s` },
      { label: "Pizza demand (rush hour)", value: `${stats.demandRushHour}/s` },
    ]) ??
    [];

  const data = useAsync(async () => {
    const response = await fetch("/api/timeseries/data");
    if (
      !response.ok ||
      response.headers.get("Content-Type") !== "application/json"
    ) {
      throw new Error("Failed to get timeseries data");
    }
    const data = await response.json();
    return data as TimeseriesData;
  });

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
      <h2>Stats</h2>
      <table
        className={classnames(
          "w-full",
          "max-w-md",
          "my-4",
          "border-collapse",
          "border-green-400",
          "border-2"
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
      <h3>Resource History</h3>
      {data.value && (
        <LineChart width={600} height={300} data={data.value.dataPoints}>
          <Line type="monotone" dataKey="coins" stroke="#F59E0B" />
          <Line type="monotone" dataKey="pizzas" stroke="#991B1B" />
          <CartesianGrid stroke="#ccc" />
          <XAxis
            dataKey="timestamp"
            domain={["auto", "auto"]}
            name="Time"
            tickFormatter={(t) => format(new Date(Number(t)), "dd/MM HH:mm")}
            type="number"
          />
          <YAxis />
          <Tooltip
            labelFormatter={(l) => format(new Date(Number(l)), "dd/MM HH:mm")}
          />
          <Legend verticalAlign="top" />
        </LineChart>
      )}
    </div>
  );
};

export default StatsContent;
