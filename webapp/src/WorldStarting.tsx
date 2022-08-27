import { formatISO9075 } from "date-fns";
import React from "react";
import { useNavigate } from "react-router";
import classnames from "classnames";
import { WorldState } from "./generated/world";
import Header from "./Header";
import { useStore } from "./store";
import styles from "./styles";
import { Link } from "react-router-dom";

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
      <div
        className={classnames(
          "container",
          "mx-auto",
          "mt-4",
          "p-4",
          "max-w-md",
          "bg-green-50",
          "prose"
        )}
      >
        <h2>The round has not yet started.</h2>
        <p>Ready for a round of Pizza Tribes? We are almost ready to go!</p>
        <p>
          The game will start at{" "}
          <span className="font-bold">
            {formatISO9075(new Date(Number(state.startTime) * 1e3))}
          </span>
          .
        </p>
      </div>
      <div className={classnames("mt-8")}>
        <Link to="/games">
          <button className={classnames(styles.primaryButton, "mr-2")}>
            Show Games
          </button>
        </Link>
        <button
          className={classnames(styles.primaryButton, "mr-2")}
          onClick={() => onClickLogout()}
        >
          Logout
        </button>
      </div>
    </div>
  );
};
