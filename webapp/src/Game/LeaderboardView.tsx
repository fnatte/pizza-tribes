import React from "react";
import { classnames } from "tailwindcss-classnames";
import { useLocation } from "react-router-dom";
import { useAsync } from "react-use";
import { ReactComponent as HeartsSvg } from "../../images/hearts.svg";
import { Leaderboard } from "../generated/leaderboard";

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

function LeaderboardView() {
  const query = useQuery();
  const skip = parseInt(query.get("skip") ?? "") || 0;

  const data = useAsync(async () => {
    const response = await fetch(`/api/leaderboard/?skip=${skip}`);
    if (
      !response.ok ||
      response.headers.get("Content-Type") !== "application/json"
    ) {
      throw new Error("Failed to get leaderboard");
    }
    return Leaderboard.fromJson(await response.json());
  }, [skip]);

  return (
    <div
      className={classnames("flex", "items-center", "flex-col", "mt-2", "p-2")}
    >
      {data.loading === null && <HeartsSvg />}
      {data.value && (
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
          <thead>
            <tr>
              <th className={classnames("p-1", "text-left", "w-1")}>Position</th>
              <th className={classnames("p-1", "pl-8", "text-left")}>User</th>
              <th className={classnames("p-1", "text-right")}>Coins</th>
            </tr>
          </thead>
          <tbody>
            {data.value.rows.map((row, i) => (
              <tr
                key={row.userId}
                className={classnames({
                  "bg-green-200": i % 2 === 0,
                })}
              >
                <td className={classnames("p-1")}>
                  {i + 1 + (data.value?.skip ?? 0)}
                </td>
                <td className={classnames("p-1", "pl-8")}>{row.username}</td>
                <td className={classnames("text-right", "p-1")}>{row.coins}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default LeaderboardView;
