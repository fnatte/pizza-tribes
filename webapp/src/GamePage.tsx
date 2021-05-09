import React, { useEffect, useRef, useState } from "react";
import { Link, Navigate, Route, Routes, useNavigate } from "react-router-dom";
import { classnames, TArg } from "tailwindcss-classnames";
import Town from "./Game/Town";
import TownLot from "./Game/TownLot";
import { useStore } from "./store";
import styles from "./styles";
import { ReactComponent as HeartsSvg } from "../images/hearts.svg";
import { ConnectionState } from "./connect";
import StatsView from "./Game/StatsView";
import MapView from "./Game/map/MapView";
import WorldEntryView from "./Game/world/WorldEntryView";
import LeaderboardView from "./Game/LeaderboardView";
import { useClickAway, useMedia } from "react-use";
import ListReportsView from "./Game/reports/ListReportsView";
import ShowReportView from "./Game/reports/ShowReportView";

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

function BurgerMenuIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 100 80" width="25" height="25" fill="#333" {...props}>
      <rect width="100" height="20"></rect>
      <rect y="30" width="100" height="20"></rect>
      <rect y="60" width="100" height="20"></rect>
    </svg>
  );
}

const NotificationCount: React.VFC<{ className?: string; count: number }> = ({
  className,
  count,
}) => {
  return (
    <div
      className={classnames(
        "flex",
        "justify-center",
        "items-center",
        "rounded-full",
        "absolute",
        "top-0",
        "right-0",
        "bg-red-700",
        "w-7",
        "h-7",
        "text-white",
        className as TArg
      )}
    >
      {count}
    </div>
  );
};

const ButtonNotificationCount: React.FC<{ count: number }> = ({ count }) => {
  return (
    <NotificationCount
      count={count}
      className={classnames(
        "transform" as TArg,
        "translate-x-1",
        "-translate-y-3"
      )}
    />
  );
};

function Navigation() {
  const isMinLg = useMedia("(min-width: 1024px)", false);
  const menuRef = useRef(null);
  const logout = useStore((state) => state.logout);
  const navigate = useNavigate();
  const [menuExpanded, setMenuExpaded] = useState(false);

  const unreads = useStore((state) =>
    state.reports.reduce((count, report) => count + (report.unread ? 1 : 0), 0)
  );

  useClickAway(menuRef, () => {
    setMenuExpaded(false);
  });

  const onClickLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <nav className={classnames("flex", "justify-center", "items-center")}>
      <Link to="/map">
        <button className={classnames(styles.button, "mr-2")}>Map</button>
      </Link>
      <Link to="/town">
        <button className={classnames(styles.button, "mr-2")}>Town</button>
      </Link>
      <Link to="/stats">
        <button className={classnames(styles.button, "mr-2")}>Stats</button>
      </Link>
      {isMinLg ? (
        <>
          <Link to="/reports">
            <button className={classnames(styles.button, "mr-2", "relative")}>
              {unreads > 0 && <ButtonNotificationCount count={unreads} />}
              Reports
            </button>
          </Link>
          <Link to="/leaderboard">
            <button className={classnames(styles.button, "mr-2")}>
              Leaderboard
            </button>
          </Link>
          <button
            className={classnames(styles.button, "mr-2")}
            onClick={() => onClickLogout()}
          >
            Logout
          </button>{" "}
        </>
      ) : (
        <div className={classnames("relative", "ml-8")} ref={menuRef}>
          <div
            className={classnames("relative", "cursor-pointer")}
            onClick={() => setMenuExpaded((s) => !s)}
          >
            <BurgerMenuIcon />
            {unreads > 0 && !menuExpanded && (
              <NotificationCount
                count={unreads}
                className={classnames(
                  "transform" as TArg,
                  "translate-x-5",
                  "-translate-y-4"
                )}
              />
            )}
          </div>
          {menuExpanded && (
            <div
              className={classnames(
                "fixed",
                "flex",
                "flex-col",
                "bg-green-200",
                "transform" as TArg,
                "sm:-translate-x-1/2",
                "sm:right-auto",
                "translate-y-2",
                "p-2",
                "right-0"
              )}
            >
              <Link to="/reports" onClick={() => setMenuExpaded(false)}>
                <button
                  className={classnames(styles.button, "mr-2", "relative")}
                >
                  {unreads > 0 && <ButtonNotificationCount count={unreads} />}
                  Reports
                </button>
              </Link>
              <Link to="/leaderboard" onClick={() => setMenuExpaded(false)}>
                <button className={classnames(styles.button, "mr-2")}>
                  Leaderboard
                </button>
              </Link>
              <button
                className={classnames(styles.button, "mr-2")}
                onClick={() => {
                  setMenuExpaded(false);
                  onClickLogout();
                }}
              >
                Logout
              </button>{" "}
            </div>
          )}
        </div>
      )}
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

  if (connectionState?.error === "unauthorized") {
    console.log("redirect to /login");
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
        <Route path="stats" element={<StatsView />} />
        <Route path="world/entry" element={<WorldEntryView />} />
        <Route path="leaderboard" element={<LeaderboardView />} />
        <Route path="reports" element={<ListReportsView />} />
        <Route path="reports/:id" element={<ShowReportView />} />
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
