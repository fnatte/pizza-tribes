import React, { useCallback, useEffect, useRef, useState } from "react";
import { Link, Navigate, Route, Routes, useNavigate } from "react-router-dom";
import classnames from "classnames";
import TownView from "./Game/TownView";
import TownLot from "./Game/TownLot";
import { useStore } from "./store";
import styles from "./styles";
import { ConnectionState } from "./connect";
import StatsView from "./Game/StatsView";
import MapView from "./Game/map/MapView";
import WorldEntryView from "./Game/world/WorldEntryView";
import LeaderboardView from "./Game/LeaderboardView";
import { useClickAway, useDebounce, useInterval, useMedia } from "react-use";
import ListReportsView from "./Game/reports/ListReportsView";
import ShowReportView from "./Game/reports/ShowReportView";
import { formatNumber, getTapInfo } from "./utils";
import HelpView from "./Game/help/HelpView";
import { useWorldState } from "./queries/useWorldState";
import { WorldStarting } from "./WorldStarting";
import { WorldEnded } from "./WorldEnded";
import { Clock, Coin, Pizza, Sparkles } from "./icons";
import MouseView from "./Game/MouseView";
import QuestsView from "./Game/quests/QuestsView";
import { useActivity } from "./useActivity";
import MouseAppearanceView from "./Game/MouseAppearanceView";
import { Loading } from "./Loading";

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

const BadgeCount: React.VFC<{ className?: string; count: number }> = ({
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
        className
      )}
    >
      {count}
    </div>
  );
};

const ButtonBadgeCount: React.FC<{ count: number }> = ({ count }) => {
  return (
    <BadgeCount
      count={count}
      className={classnames("translate-x-1", "-translate-y-3")}
    />
  );
};

function Navigation() {
  const isMinMd = useMedia("(min-width: 767px)", false);
  const isMinLg = useMedia("(min-width: 1024px)", false);
  const isMinXl = useMedia("(min-width: 1280px)", false);
  const menuRef = useRef(null);
  const logout = useStore((state) => state.logout);
  const discoveries = useStore((state) => state.gameState.discoveries);
  const navigate = useNavigate();
  const [menuExpanded, setMenuExpaded] = useState(false);

  const unreads = useStore((state) =>
    state.reports.reduce((count, report) => count + (report.unread ? 1 : 0), 0)
  );
  const now = new Date();

  const tapBadgeCount = useStore((store) =>
    Object.values(store.gameState.lots).reduce(
      (count, lot) =>
        count + (getTapInfo(lot, discoveries, now).tapsRemaining > 0 ? 1 : 0),
      0
    )
  );

  const questBadgeCount = useStore((state) =>
    Object.values(state.gameState.quests).reduce(
      (count, questState) =>
        count +
        ((questState.completed && !questState.claimedReward) ||
        !questState.opened
          ? 1
          : 0),
      0
    )
  );

  useClickAway(menuRef, () => {
    setMenuExpaded(false);
  });

  const onClickLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <nav
      className={classnames(
        "flex",
        "justify-center",
        "items-center",
        "mt-2",
        "xs:mt-0"
      )}
      data-cy="main-nav"
    >
      <Link to="map">
        <button className={classnames(styles.primaryButton, "mr-2")}>
          Map
        </button>
      </Link>
      <Link to="town">
        <button
          className={classnames(styles.primaryButton, "mr-2", "relative")}
        >
          {tapBadgeCount > 0 && <ButtonBadgeCount count={tapBadgeCount} />}
          Town
        </button>
      </Link>
      <Link to="quests">
        <button
          className={classnames(styles.primaryButton, "mr-2", "relative")}
        >
          {questBadgeCount > 0 && <ButtonBadgeCount count={questBadgeCount} />}
          Quests
        </button>
      </Link>
      {isMinMd && (
        <Link to="stats">
          <button className={classnames(styles.primaryButton, "mr-2")}>
            Stats
          </button>
        </Link>
      )}
      {isMinLg && (
        <Link to="reports">
          <button
            className={classnames(styles.primaryButton, "mr-2", "relative")}
          >
            {unreads > 0 && <ButtonBadgeCount count={unreads} />}
            Reports
          </button>
        </Link>
      )}
      {isMinXl && (
        <>
          <Link to="leaderboard">
            <button className={classnames(styles.primaryButton, "mr-2")}>
              Leaderboard
            </button>
          </Link>
          <Link to="help">
            <button className={classnames(styles.primaryButton, "mr-2")}>
              Help
            </button>
          </Link>
          <button
            className={classnames(styles.primaryButton, "mr-2")}
            onClick={() => onClickLogout()}
          >
            Logout
          </button>{" "}
        </>
      )}
      {!isMinXl && (
        <div className={classnames("relative", "ml-8")} ref={menuRef}>
          <button
            className={classnames("relative", "cursor-pointer")}
            onClick={() => setMenuExpaded((s) => !s)}
            aria-expanded={menuExpanded}
            data-cy="menu-expand-button"
          >
            <BurgerMenuIcon />
            {unreads > 0 && !menuExpanded && (
              <BadgeCount
                count={unreads}
                className={classnames("translate-x-5", "-translate-y-4")}
              />
            )}
          </button>
          {menuExpanded && (
            <div
              className={classnames(
                "fixed",
                "flex",
                "flex-col",
                "bg-green-200",
                "z-10",
                "sm:-translate-x-1/2",
                "sm:right-auto",
                "translate-y-2",
                "p-2",
                "right-0"
              )}
            >
              {!isMinMd && (
                <Link to="stats" onClick={() => setMenuExpaded(false)}>
                  <button
                    className={classnames(
                      styles.primaryButton,
                      "mr-2",
                      "w-full"
                    )}
                  >
                    Stats
                  </button>
                </Link>
              )}
              {!isMinLg && (
                <Link to="reports" onClick={() => setMenuExpaded(false)}>
                  <button
                    className={classnames(
                      styles.primaryButton,
                      "mr-2",
                      "w-full",
                      "relative"
                    )}
                  >
                    {unreads > 0 && <ButtonBadgeCount count={unreads} />}
                    Reports
                  </button>
                </Link>
              )}
              <Link to="leaderboard" onClick={() => setMenuExpaded(false)}>
                <button
                  className={classnames(styles.primaryButton, "mr-2", "w-full")}
                >
                  Leaderboard
                </button>
              </Link>
              <Link to="help" onClick={() => setMenuExpaded(false)}>
                <button
                  className={classnames(styles.primaryButton, "mr-2", "w-full")}
                >
                  Help
                </button>
              </Link>
              <button
                className={classnames(styles.primaryButton, "mr-2", "w-full")}
                onClick={() => {
                  setMenuExpaded(false);
                  onClickLogout();
                }}
              >
                Logout
              </button>
            </div>
          )}
        </div>
      )}
    </nav>
  );
}

function GameTitle() {
  return (
    <div
      className={classnames(
        "justify-center",
        "text-xl",
        "mt-2",
        "hidden",
        "xs:flex"
      )}
    >
      Pizza Tribes
    </div>
  );
}

function ResourceBar() {
  const { pizzas, coins } = useStore(
    (state) => state.gameState.resources ?? { coins: 0, pizzas: 0 }
  );
  const pizzaPrice = useStore((state) =>
    Math.max(1, Math.min(15, state.gameState.pizzaPrice))
  );
  const stats = useStore((state) => state.gameStats);
  const clock = useMouseClock();

  const [displayPizzas, setDisplayPizzas] = useState(pizzas);
  const [displayCoins, setDisplayCoins] = useState(coins);

  useInterval(() => {
    if (stats === null) {
      return;
    }
    const dt = 0.1;

    const demand =
      (clock.isRushHour ? stats.demandRushHour : stats.demandOffpeak) * dt;
    const pizzasProduced = stats.pizzasProducedPerSecond * dt;
    const pizzasAvailable = pizzas + pizzasProduced;
    const maxSellsByMice = stats.maxSellsByMicePerSecond * dt;
    const pizzasSold = Math.min(demand, maxSellsByMice, pizzasAvailable);

    setDisplayCoins((c) => c + pizzasSold * pizzaPrice);
    setDisplayPizzas((p) => Math.max(p + pizzasProduced - pizzasSold, 0));
  }, 100);

  useEffect(() => setDisplayCoins(coins), [coins]);
  useEffect(() => setDisplayPizzas(pizzas), [pizzas]);

  const pizzasDisplayText = formatNumber(Math.floor(displayPizzas));
  const pizzasWidth = getTextWidthClass(pizzasDisplayText);

  const coinsDisplayText = formatNumber(Math.floor(displayCoins));
  const coinsWidth = getTextWidthClass(coinsDisplayText);

  return (
    <div
      className={classnames(
        "flex",
        "justify-center",
        "flex-wrap",
        "text-xl",
        "sm:text-2xl",
        "mt-2",
        "gap-2"
      )}
    >
      <div className="flex gap-1">
        <Coin className={"h-[1.25em] w-[1.25em]"} />
        <div
          className={classnames(
            coinsWidth,
            "h-[1.25em]",
            "overflow-hidden",
            "[contain:strict]"
          )}
          data-cy="resource-bar-coins"
        >
          {coinsDisplayText}
        </div>
      </div>
      <div className="flex gap-1">
        <Pizza className={"h-[1.25em] w-[1.25em]"} />
        <div
          className={classnames(
            pizzasWidth,
            "h-[1.25em]",
            "overflow-hidden",
            "[contain:strict]"
          )}
          data-cy="resource-bar-pizzas"
        >
          {pizzasDisplayText}
        </div>
      </div>
      <div className="flex gap-2">
        <div className="flex gap-1">
          <Clock className="h-[1.25em] w-[1.25em]" />
          <div className="w-16 h-[1.25em] overflow-hidden [contain:strict]">
            {clock.formatted}
          </div>
        </div>
        {clock.isRushHour && (
          <div className="flex gap-1">
            <Sparkles className="h-[1.25em] w-[1.25em]" />
            <div>Rush Hour!</div>
          </div>
        )}
      </div>
    </div>
  );
}

function getTextWidthClass(text: string) {
  if (text.length < 3) {
    return "w-[3ch]";
  } else if (text.length < 5) {
    return "w-[5ch]";
  } else if (text.length < 8) {
    return "w-[8ch]";
  } else {
    return "w-[10ch]";
  }
}

function Separator() {
  return (
    <hr
      className={classnames(
        "border-t-2",
        "border-gray-300",
        "my-1",
        "md:my-4",
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

function GamePage(): JSX.Element {
  const connectionState = useStore((state) => state.connectionState);
  const user = useStore((state) => state.user);
  const gameData = useStore((state) => state.gameData);
  const gameDataLoading = useStore((state) => state.gameDataLoading);
  const fetchGameData = useStore((state) => state.fetchGameData);
  const reportActivity = useStore((state) => state.reportActivity);
  const start = useStore((state) => state.start);

  const needToConnect =
    !connectionState ||
    (!connectionState.connected &&
      !connectionState.connecting &&
      !connectionState.error);

  useEffect(() => {
    if (needToConnect) {
      start();
    }
  }, [needToConnect]);

  const { isLoading, data, error, refetch } = useWorldState();

  useEffect(() => {
    if (user && error) {
      refetch();
    }
  }, [refetch, user, error]);

  const onReportActivity = useCallback(() => reportActivity(), []);
  useActivity(onReportActivity, {
    enabled: connectionState?.connected ?? false,
  });

  // Delay showing connection popup for 2s, but remove it instantly if connection is restored
  const [showConnectionPopup, setShowConnectionPopup] = useState(false);
  useDebounce(
    () => {
      setShowConnectionPopup(connectionState?.connecting ?? false);
    },
    2000,
    [connectionState?.connecting]
  );
  useEffect(() => {
    if (!connectionState?.connecting) {
      setShowConnectionPopup(false);
    }
  }, [connectionState?.connecting]);

  useEffect(() => {
    if (gameData === null && !gameDataLoading) {
      fetchGameData();
    }
  }, [gameData, gameDataLoading]);

  if (connectionState?.error === "unauthorized") {
    return <Navigate to="/login" replace />;
  }

  if (gameData === null) {
    return <Loading />;
  }

  if (user === null) {
    return <Loading />;
  }

  if (isLoading || (error && user) || connectionState === null) {
    return <Loading />;
  }

  if (data && data.type.oneofKind === "starting") {
    return <WorldStarting state={data} />;
  }

  if (data && data.type.oneofKind === "ended") {
    return <WorldEnded state={data} />;
  }

  return (
    <div>
      <GameTitle />
      <Navigation />
      <ResourceBar />
      <Separator />
      <Routes>
        <Route path="map" element={<MapView />} />
        <Route path="town/:id/*" element={<TownLot />} />
        <Route path="town" element={<TownView />} />
        <Route path="quests" element={<QuestsView />} />
        <Route path="stats" element={<StatsView />} />
        <Route path="world/entry" element={<WorldEntryView />} />
        <Route path="leaderboard" element={<LeaderboardView />} />
        <Route path="reports" element={<ListReportsView />} />
        <Route path="reports/:id" element={<ShowReportView />} />
        <Route path="mouse/:id/appearance" element={<MouseAppearanceView />} />
        <Route path="mouse/:id" element={<MouseView />} />
        <Route path="help/*" element={<HelpView />} />
        <Route path="/" element={<Navigate to="town" replace />} />
      </Routes>
      {showConnectionPopup && (
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
