import { formatISO9075 } from "date-fns";
import React from "react";
import { useNavigate } from "react-router";
import { classnames } from "tailwindcss-classnames";
import { WorldState } from "./generated/world";
import Header from "./Header";
import { useStore } from "./store";
import styles from "./styles";

export const WorldStarting: React.FC<{ state: WorldState }> = ({ state }) => {
  const logout = useStore((state) => state.logout);
  const navigate = useNavigate();
  const onClickLogout = () => {
    logout();
    navigate("/login");
  };

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
      <div className={classnames("p-2", "prose" as any)}>
        <h2>The round has not yet started.</h2>
        <p>Ready for a round of Pizza Tribes? We are almost ready to go!</p>
        <p>
          The game will start at{" "}
          {formatISO9075(new Date(Number(state.startTime) * 1e3))}.
        </p>
        <div className={classnames("mt-8")}>
          <button
            className={classnames(styles.primaryButton, "mr-2")}
            onClick={() => onClickLogout()}
          >
            Logout
          </button>{" "}
        </div>
      </div>
    </div>
  );
};
