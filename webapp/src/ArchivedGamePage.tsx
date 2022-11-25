import React from "react";
import classnames from "classnames";
import { Link, useNavigate, useParams } from "react-router-dom";
import Header from "./Header";
import { useCentralLeaderboard } from "./queries/useCentralLeaderboard";
import { useCentralLeaderboardMeRank } from "./queries/useCentralLeaderboardMeRank";
import { LeaderboardTable, TopThree } from "./Game/LeaderboardView";
import { Hearts } from "./icons";
import styles from "./styles";

function WinnerDisplay({ username }: { username: string }) {
  return (
    <div
      className={classnames(
        "container",
        "mx-auto",
        "max-w-md",
        "bg-green-200",
        "flex",
        "flex-col",
        "items-center",
        "mt-4"
      )}
    >
      <h3 className={classnames("mt-2", "mb-0")}>Winner</h3>
      <h2 className={classnames("flex", "items-center", "h-16", "mt-0")}>
        {username}
      </h2>
    </div>
  );
}

function ArchivedGamePage() {
  const navigate = useNavigate();
  const { gameId } = useParams();
  const { data: myRank } = useCentralLeaderboardMeRank(gameId ?? "", {
    enabled: gameId !== undefined,
  });
  const { data, isLoading } = useCentralLeaderboard(gameId ?? "", myRank ?? 0, {
    enabled: gameId !== undefined && myRank !== undefined,
  });

  return (
    <div
      className={classnames(
        "flex",
        "flex-col",
        "justify-center",
        "items-center"
      )}
    >
      <Link to="/">
        <Header />
      </Link>
      <div className={classnames("mt-4")}>
        <button
          className={classnames(styles.primaryButton, "mr-2")}
          onClick={() => navigate(-1)}
        >
          Back
        </button>
      </div>
      {isLoading && <Hearts className="fill-green-600" />}
      {data && (
        <>
          <WinnerDisplay username={data.rows[0].username} />
          <h3 className={classnames("mt-8", "text-center")}>Top 3 Tribes</h3>
          <TopThree rows={data.rows.slice(0, 3)} />
          <h3 className={classnames("mt-8", "text-center")}>Leaderboard</h3>
          <LeaderboardTable leaderboard={data} />
        </>
      )}
    </div>
  );
}

export default ArchivedGamePage;
