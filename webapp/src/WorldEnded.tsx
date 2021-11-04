import React, { useState } from "react";
import { useNavigate } from "react-router";
import { classnames } from "tailwindcss-classnames";
import { WorldState } from "./generated/world";
import Header from "./Header";
import { useUser } from "./queries/useUser";
import { useStore } from "./store";
import styles from "./styles";
import { ReactComponent as HeartsSvg } from "../images/hearts.svg";
import { useTimeoutFn } from "react-use";
import LeaderboardView from "./Game/LeaderboardView";

const InlineLoader: React.FC = () => {
  return (
    <span>
      <HeartsSvg />
    </span>
  );
};

export const WorldEnded: React.FC<{ state: WorldState }> = ({ state }) => {
  const logout = useStore((state) => state.logout);
  const navigate = useNavigate();
  const onClickLogout = () => {
    logout();
    navigate("/login");
  };

  if (state.type.oneofKind !== "ended") {
    return null;
  }

  const [enabled, setEnabled] = useState(false);
  useTimeoutFn(() => setEnabled(true), 3000);

  const { data } = useUser(state.type.ended.winnerUserId, { enabled });

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
      <Header />
      <div
        className={classnames(
          "container",
          "mx-auto",
          "max-w-md",
          "bg-green-200",
          "flex",
          "flex-col",
          "items-center"
        )}
      >
        <h3 className={classnames("mt-2", "mb-0")}>And the winner is...</h3>
        <h2 className={classnames("flex", "items-center", "h-16", "mt-0")}>
          {data ? data.username : <InlineLoader />}
        </h2>
      </div>
      {data ? (
        <article
          className={classnames(
            "text-black",
            "p-4",
            "container",
            "mx-auto",
            "mt-4",
            "px-4",
            "max-w-md",
            "bg-green-50",
            "prose" as any
          )}
        >
          <p>Hi,</p>
          <p>
            This game round has ended. Congratulations to {data.username} for
            reaching 10,000,000 coins fastest.
          </p>
          <p>Thanks for playing, hope to see you next round.</p>
          <p>
            Yours truly,
            <br />
            Jerry
          </p>
        </article>
      ) : null}
      <LeaderboardView />
      <div className={classnames("mt-10")}>
        <button
          className={classnames(styles.primaryButton, "mr-2")}
          onClick={() => onClickLogout()}
        >
          Logout
        </button>{" "}
      </div>
    </div>
  );
};
