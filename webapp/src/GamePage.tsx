import React, { useEffect, useState } from "react";
import { Link, Navigate, Route, Routes } from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import Town from "./Game/Town";
import TownLot from "./Game/TownLot";
import { useStore } from "./store";
import styles from "./styles";

type ClockState = {
  formatted: string;
  isRushHour: boolean;
};

const clockStateNow = () => {
  const mt = Date.now() * 24;
  const t = new Date(mt);
  const h = t.getUTCHours();
  const m = t.getUTCMinutes();
  return {
    formatted: `${("0" + h).substr(-2)}:${("0" + m).substr(-2)}`,
    isRushHour: (h >= 11 && h < 13) || (h >= 18 && h < 21),
  };
};

const useMouseClock = () => {
  const [clock, setClock] = useState<ClockState>(clockStateNow());

  useEffect(() => {
    const handle = window.setInterval(() => {
      setClock(clockStateNow());
    }, 1000);
    return () => window.clearInterval(handle);
  }, [setClock]);

  return clock;
};

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

function ClockEmoji() {
  return <span>🕓</span>;
}

function SparkleEmoji() {
  return <span>✨</span>;
}

function ResourceBar() {
  const { pizzas, coins } = useStore((state) => state.gameState.resources);
  const clock = useMouseClock();

  return (
    <div className={classnames("flex", "justify-center", "flex-wrap", "text-xl", "sm:text-2xl", "mt-2")}>
      <span className={classnames("px-6", "mb-2")}>
        <CoinEmoji /> {coins.toString()}
      </span>
      <span className={classnames("px-6", "mb-2")}>
        <PizzaEmoji /> {pizzas.toString()}
      </span>
      <span className={classnames("px-6", "mb-2")}>
        <span className={classnames("px-2")}>
          <ClockEmoji />{" "}
          <span style={{ minWidth: 60, display: "inline-block" }}>
            {clock.formatted}
          </span>
        </span>
        {clock.isRushHour && (
          <span className={classnames("px-2")}>
            <SparkleEmoji /> Rush Hour!
          </span>
        )}
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
