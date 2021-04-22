import React from "react";
import {Navigate} from "react-router-dom";
import {classnames} from "tailwindcss-classnames";
import {useStore} from "./store";
import styles from "./styles";

function GamePage() {
  const user = useStore((state) => state.user);
  const logout = useStore((state) => state.logout);
  const tap = useStore((state) => state.tap);
  const { pizzas, coins } = useStore((state) => state.gameState.resources);

  if (user === null) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div>
      <button className={classnames(styles.button)} onClick={() => logout()}>
        Logout
      </button>
      <button className={classnames(styles.button)} onClick={() => tap()}>
        Tap
      </button>
      <div>
        <span>Coins: {coins.toString()}</span>
        <span> | </span>
        <span>Pizzas: {pizzas.toString()}</span>
      </div>
    </div>
  );
}

export default GamePage;
