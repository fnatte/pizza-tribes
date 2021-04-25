import React, {} from "react";
import { Link, Navigate, Route, Routes } from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import Town from "./Game/Town";
import TownLot from "./Game/TownLot";
import { useStore } from "./store";
import styles from "./styles";

function Navigation() {
  const logout = useStore((state) => state.logout);

  return (
    <nav className={classnames("flex", "justify-center")}>
      <Link to="/map">
        <button className={classnames(styles.button, "mr-2")}>Map</button>
      </Link>
      <Link to="/town">
        <button className={classnames(styles.button, "mr-2")}>Town</button>
      </Link>
      <button
        className={classnames(styles.button, "mr-2")}
        onClick={() => logout()}
      >
        Logout
      </button>
    </nav>
  );
}

function GameTitle() {
  return (
    <div className={classnames("flex", "justify-center", "text-xl", "mt-2")}>
      Pizza Mouse
    </div>
  );
}

function CoinEmoji() {
  return <span>🪙</span>;
}

function PizzaEmoji() {
  return <span>🍕</span>;
}

function ResourceBar() {
  const { pizzas, coins } = useStore((state) => state.gameState.resources);

  return (
    <div className={classnames("flex", "justify-center", "text-2xl", "mt-2")}>
      <span className={classnames("px-6")}>
        <CoinEmoji /> {coins.toString()}
      </span>
      <span className={classnames("px-6")}>
        <PizzaEmoji /> {pizzas.toString()}
      </span>
    </div>
  );
}

function Separator() {
  return (
    <hr
      className={classnames(
        "border-t-2",
        "border-gray-300",
        "my-4",
        "w-10/12",
        "mx-auto"
      )}
    />
  );
}

function Map() {
  return (
    <div className={classnames("flex", "justify-center", "mt-2")}>
      <h2>Map</h2>
    </div>
  );
}

function GamePage() {
  const user = useStore((state) => state.user);

  if (user === null) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div>
      <GameTitle />
      <Navigation />
      <ResourceBar />
      <Separator />
      <Routes>
        <Route path="map" element={<Map />} />
        <Route path="town/:id" element={<TownLot />} />
        <Route path="town" element={<Town />} />
        <Route path="/" element={<Navigate to="/town" replace />} />
      </Routes>
    </div>
  );
}

export default GamePage;
