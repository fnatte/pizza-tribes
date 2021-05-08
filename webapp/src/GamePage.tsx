import React, { useEffect, useState } from "react";
import { Link, Navigate, Route, Routes, useNavigate } from "react-router-dom";
import { classnames, TArg } from "tailwindcss-classnames";
import Town from "./Game/Town";
import TownLot from "./Game/TownLot";
import { useStore } from "./store";
import styles from "./styles";
import { ReactComponent as HeartsSvg } from "../images/hearts.svg";
import { ConnectionState } from "./connect";
import StatsContent from "./Game/StatsContent";
import MapView from "./Game/map/MapView";
import WorldEntryView from "./Game/world/WorldEntryView";

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
  const navigate = useNavigate();
  const onClickLogout = () => {
    logout();
    navigate("/login");
  }

  return (
    <nav className={classnames("flex", "justify-center")}>
      <Link to="/map">
        <button className={classnames(styles.button, "mr-2")}>Map</button>
      </Link>
      <Link to="/town">
        <button className={classnames(styles.button, "mr-2")}>Town</button>
      </Link>
      <Link to="/stats">
        <button className={classnames(styles.button, "mr-2")}>Stats</button>
      </Link>
      <button
        className={classnames(styles.button, "mr-2")}
        onClick={() => onClickLogout()}
      >
        Logout
      </button>
    </nav>
  );
}

function GameTitle() {
  return (
    <div className={classnames("flex", "justify-center", "text-xl", "mt-2")}>
      Pizza Tribes
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
    <div
      className={classnames(
        "flex",
        "justify-center",
        "flex-wrap",
        "text-xl",
        "sm:text-2xl",
        "mt-2"
      )}
    >
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

const ConnectionPopup: React.VFC<{ connectionState: ConnectionState }> = ({
  connectionState,
}) => {
  const [wasConnected, setWasConnected] = useState(connectionState.connected);
  useEffect(() => {
    if (!wasConnected) {
      setWasConnected(true);
    }
  }, [connectionState.connected]);

  return (
    <div
      className={classnames(
        "fixed",
        "top-1/2",
        "left-1/2",
        "transform" as TArg,
        "-translate-x-1/2",
        "-translate-y-1/2",
        "p-4",
        "bg-white"
      )}
    >
      {wasConnected ? (
        <>
          <h2>Connection lost</h2>
          <p>Trying to reconnect...</p>
          <p>Reconnect attempts: {connectionState.reconnectAttempts}</p>
        </>
      ) : (
        <h2>Connecting</h2>
      )}
    </div>
  );
};

function Loading() {
  return (
      <div
        className={classnames(
          "fixed",
          "left-1/2",
          "top-1/2",
          "transform" as TArg,
          "-translate-y-1/2",
          "-translate-x-1/2"
        )}
      >
        <HeartsSvg />
      </div>
  );
}

function GamePage(): JSX.Element {
  const connectionState = useStore((state) => state.connectionState);
  const user = useStore((state) => state.user);
  const gameData = useStore((state) => state.gameData);
  const gameDataLoading = useStore((state) => state.gameDataLoading);
  const fetchGameData = useStore((state) => state.fetchGameData);
  const start = useStore((state) => state.start);

  useEffect(() => {
    start();
  }, []);

  if (connectionState?.error === 'unauthorized') {
    console.log('redirect to /login');
    return <Navigate to="/login" replace />;
  }

  if (gameData === null) {
    if (!gameDataLoading) {
      fetchGameData();
    }
    return <Loading />;
  }

  if (user === null) {
    return <Loading />;
  }

  return (
    <div>
      <GameTitle />
      <Navigation />
      <ResourceBar />
      <Separator />
      <Routes>
        <Route path="map" element={<MapView />} />
        <Route path="town/:id" element={<TownLot />} />
        <Route path="town" element={<Town />} />
        <Route path="stats" element={<StatsContent />} />
        <Route path="world/entry" element={<WorldEntryView />} />
        <Route path="/" element={<Navigate to="/town" replace />} />
      </Routes>
      {connectionState?.connecting && (
        <div
          className={classnames(
            "fixed",
            "top-0",
            "left-0",
            "bg-gray-600",
            "bg-opacity-50",
            "w-full",
            "h-full"
          )}
        >
          <ConnectionPopup connectionState={connectionState} />
        </div>
      )}
    </div>
  );
}

export default GamePage;
