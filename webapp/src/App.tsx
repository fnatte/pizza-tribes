import React from "react";
import { classnames } from "tailwindcss-classnames";
import LoginForm from "./LoginForm";
import { useStore } from "./store";

function Welcome() {
  return (
    <div
      className={classnames(
        "flex",
        "justify-center",
        "flex-col",
        "items-center"
      )}
    >
      <h1 className={classnames("flex", "justify-center", "p-8", "text-4xl")}>
        Pizza Mouse
      </h1>
      <div className={classnames("text-2xl")}>🍕🍕🍕🍕</div>
    </div>
  );
}

function Game() {
  const logout = useStore(state => state.logout);
  const tap = useStore(state => state.tap);
  const { pizzas, coins } = useStore(state => state.gameState.resources);

  return (
    <div>
      <button
        className={classnames(
          "my-2",
          "py-2",
          "px-8",
          "text-white",
          "bg-green-600",
        )}
        onClick={() => logout()}
      >
        Logout
      </button>
      <button
        className={classnames(
          "my-2",
          "py-2",
          "px-8",
          "text-white",
          "bg-green-600",
        )}
        onClick={() => tap()}
      >
        Tap
      </button>
      <div>
        <span>Coins: {coins}</span>
        <span> | </span>
        <span>Pizzas: {pizzas}</span>
      </div>
    </div>
  );
}

function App() {
  const user = useStore((state) => state.user);
  const start = useStore((state) => state.start);

  return (
    <div>
      {user ? (
        <Game />
      ) : (
        <>
          <Welcome />
          <LoginForm onLogin={() => start()} />
        </>
      )}
    </div>
  );
}

export default App;
