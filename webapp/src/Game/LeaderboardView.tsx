import React from "react";
import classnames from "classnames";
import { useLocation } from "react-router-dom";
import { useAsync } from "react-use";
import { ReactComponent as HeartsSvg } from "../../images/hearts.svg";
import { Leaderboard } from "../generated/leaderboard";
import { formatNumber } from "../utils";
import { apiFetch } from "../api";

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

function TopThree({ rows }: { rows: { username: string; coins: string }[] }) {
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
                "px-2",
                "bg-green-200",
                "rounded",
                "relative"
              )}
            >
              <div
                className={classnames(
                  "h-6",
                  "bg-green-700",
                  "rounded",
                  "absolute",
                  "top-0",
                  "left-0",
                  "m-1"
                )}
                style={{
                  width: `${(Number(coins) / 10_000_000) * 100}%`,
                  minWidth: 5,
                }}
              />
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

function LeaderboardTable({ leaderboard }: { leaderboard: Leaderboard }) {
  return (
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
          >
            <td className={classnames("p-1")}>
              {i + 1 + (leaderboard.skip ?? 0)}
            </td>
            <td className={classnames("p-1", "pl-8")}>{row.username}</td>
            <td className={classnames("text-right", "p-1")}>
              {formatNumber(Number(row.coins))}
            </td>
            <td className={classnames("text-right", "p-1")}>
              {formatNumber((Number(row.coins) / 10_000_000) * 100)}%
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

  const data = useAsync(async () => {
    let skipParam = skip;

    if (!query.has("skip")) {
      const response = await apiFetch("/leaderboard/me/rank");
      if (
        !response.ok ||
        response.headers.get("Content-Type") !== "application/json"
      ) {
        throw new Error("Failed to get leaderboard rank");
      }
      const rank = await response.json();
      skipParam = Math.max(0, rank - 10);
    }

    const response = await apiFetch(`/leaderboard/?skip=${skipParam}`);
    if (
      !response.ok ||
      response.headers.get("Content-Type") !== "application/json"
    ) {
      throw new Error("Failed to get leaderboard");
    }
    return Leaderboard.fromJson(await response.json());
  }, [skip]);

  const showTopThree =
    skip === 0 && data.value?.rows.some((x) => Number(x.coins) >= 300_000);

  return (
    <div
      className={classnames("flex", "items-center", "flex-col", "mt-2", "p-2")}
    >
      {data.loading === null && <HeartsSvg />}
      {data.value && (
        <>
          {showTopThree ? (
            <>
              <h3 className={classnames("mt-8", "text-center")}>
                Top 3 Tribes
              </h3>
              <TopThree rows={data.value.rows.slice(0, 3)} />
            </>
          ) : null}
          <h3 className={classnames("mt-8", "text-center")}>Leaderboard</h3>
          <LeaderboardTable leaderboard={data.value} />
        </>
      )}
    </div>
  );
}

export default LeaderboardView;
