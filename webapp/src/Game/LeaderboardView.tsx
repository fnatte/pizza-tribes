import React, { useState } from "react";
import classnames from "classnames";
import { useLocation } from "react-router-dom";
import { ReactComponent as HeartsSvg } from "../../images/hearts.svg";
import { DemandLeaderboard, Leaderboard } from "../generated/leaderboard";
import { formatNumber } from "../utils";
import { useCoinsLeaderboard } from "../queries/useCoinsLeaderboard";
import { useDemandLeaderboard } from "../queries/useDemandLeaderboard";

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

export function TopThree({ rows }: { rows: { username: string; coins: string }[] }) {
  return (
    <div className={classnames("max-w-lg", "w-full")}>
      <p className={classnames("text-gray-700", "text-center", "mb-2")}>
        The first tribe to reach 10,000,000 coins wins the round.
      </p>
      <div
        className={classnames(
          "grid",
          "gap-4",
          "content-center",
          "items-center"
        )}
        style={{ gridTemplateColumns: "fit-content(33%) 1fr" }}
      >
        {rows.map(({ username, coins }) => (
          <React.Fragment key={username}>
            <div className={classnames("overflow-hidden", "overflow-ellipsis")}>
              {username}
            </div>
            <div
              className={classnames(
                "w-full",
                "h-8",
                "p-1",
                "bg-green-200",
                "rounded"
              )}
            >
              <div className="relative w-full h-full">
                <div
                  className={classnames(
                    "h-6",
                    "bg-green-700",
                    "rounded",
                    "absolute",
                    "top-0",
                    "left-0"
                  )}
                  style={{
                    width: `${Math.min(
                      100,
                      (Number(coins) / 10_000_000) * 100
                    )}%`,
                    minWidth: 5,
                  }}
                />
              </div>
            </div>
          </React.Fragment>
        ))}
        <div
          className={classnames(
            "-mt-1",
            "col-start-2",
            "flex",
            "justify-between"
          )}
        >
          <div>0</div>
          <div>10,000,000</div>
        </div>
      </div>
    </div>
  );
}

export function LeaderboardTable({ leaderboard }: { leaderboard: Leaderboard }) {
  return (
    <table
      data-cy="coins-leaderboard-table"
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
          <th className={classnames("p-1", "text-right")}>Win %</th>
        </tr>
      </thead>
      <tbody>
        {leaderboard.rows.map((row, i) => (
          <tr
            key={row.userId}
            className={classnames({
              "bg-green-200": i % 2 === 0,
            })}
            data-cy="leaderboard-row"
          >
            <td className={classnames("p-1")}>
              {i + 1 + (leaderboard.skip ?? 0)}
            </td>
            <td className={classnames("p-1", "pl-8")}>{row.username}</td>
            <td className={classnames("text-right", "p-1")}>
              {formatNumber(Number(row.coins))}
            </td>
            <td className={classnames("text-right", "p-1")}>
              {new Intl.NumberFormat(undefined, {
                style: "percent",
                maximumFractionDigits: 1,
              }).format(Math.min(Number(row.coins) / 10_000_000, 1))}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function DemandLeaderboardTable({
  leaderboard,
}: {
  leaderboard: DemandLeaderboard;
}) {
  return (
    <table
      data-cy="demand-leaderboard-table"
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
          <th className={classnames("p-1", "text-right")}>Demand</th>
          <th className={classnames("p-1", "text-right")}>Market Share</th>
        </tr>
      </thead>
      <tbody>
        {leaderboard.rows.map((row, i) => (
          <tr
            key={row.userId}
            className={classnames({
              "bg-green-200": i % 2 === 0,
            })}
            data-cy="leaderboard-row"
          >
            <td className={classnames("p-1")}>
              {i + 1 + (leaderboard.skip ?? 0)}
            </td>
            <td className={classnames("p-1", "pl-8")}>{row.username}</td>
            <td className={classnames("text-right", "p-1")}>
              {new Intl.NumberFormat(undefined, {
                maximumFractionDigits: 2,
              }).format(row.demand)}
            </td>
            <td className={classnames("text-right", "p-1")}>
              {new Intl.NumberFormat(undefined, {
                style: "percent",
                maximumFractionDigits: 1,
              }).format(row.marketShare)}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
function LeaderboardView() {
  const query = useQuery();
  const skip = parseInt(query.get("skip") ?? "") || 0;

  const [selectedLeaderboard, setSelectedLeaderboard] = useState<
    "coins" | "demand"
  >("coins");

  const coinsLeaderboard = useCoinsLeaderboard(skip);
  const demandLeaderboard = useDemandLeaderboard(undefined, {
    enabled: selectedLeaderboard === "demand",
  });

  const showTopThree =
    skip === 0 &&
    coinsLeaderboard.data?.rows.some((x) => Number(x.coins) >= 300_000);

  return (
    <div
      className={classnames("flex", "items-center", "flex-col", "mt-2", "p-2")}
    >
      {showTopThree ? (
        <>
          <h3 className={classnames("mt-8", "text-center")}>Top 3 Tribes</h3>
          {coinsLeaderboard.data && (
            <TopThree rows={coinsLeaderboard.data.rows.slice(0, 3)} />
          )}
        </>
      ) : null}

      <h3 className={classnames("mt-8", "text-center")}>Leaderboard</h3>
      <button className="flex justify-center items-center my-2 bg-white border-2 border-green-600">
        <div
          className={classnames("py-1 px-2 w-20", {
            "bg-green-600 text-white": selectedLeaderboard === "coins",
          })}
          onClick={() => setSelectedLeaderboard("coins")}
        >
          Coins
        </div>
        <div
          className={classnames("py-1 px-2 w-20", {
            "bg-green-600 text-white": selectedLeaderboard === "demand",
          })}
          onClick={() => setSelectedLeaderboard("demand")}
        >
          Demand
        </div>
      </button>
      {selectedLeaderboard === "coins" && coinsLeaderboard.data && (
        <LeaderboardTable leaderboard={coinsLeaderboard.data} />
      )}
      {selectedLeaderboard === "coins" && coinsLeaderboard.isLoading && (
        <HeartsSvg className="mt-4 fill-green-600" />
      )}
      {selectedLeaderboard === "demand" && demandLeaderboard.data && (
        <DemandLeaderboardTable leaderboard={demandLeaderboard.data} />
      )}
      {selectedLeaderboard === "demand" && demandLeaderboard.isLoading && (
        <HeartsSvg className="mt-4 fill-green-600" />
      )}
    </div>
  );
}

export default LeaderboardView;
