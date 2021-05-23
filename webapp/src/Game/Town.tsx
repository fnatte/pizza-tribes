import React, { useState } from "react";
import { useEffect, useRef } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { useLocalStorage, useMedia } from "react-use";
import { classnames, TArg } from "tailwindcss-classnames";
import { useStore } from "../store";
import ConstructionQueue from "./ConstructionQueue";
import Population from "./Population";
import classes from "./town.module.css";
import TownExpandMenu from "./TownExpandMenu";
import TownSvg from "./TownSvg";
import TravelQueue from "./TravelQueue";

function Town() {
  const ref = useRef<SVGSVGElement>(null);
  const navigate = useNavigate();
  const lots = useStore((state) => state.gameState.lots);
  const constructionQueue = useStore((state) => state.gameState.constructionQueue);

  const onLotClick = (lotId: string) => {
    navigate(`/town/${lotId.replace("lot", "")}`);
  };

  const isMinLg = useMedia("(min-width: 1024px)", false);
  const isMinSm = useMedia("(min-width: 640px)", false);

  const [minimizedTravelQueue, setMinimizedTravelQueue] = useState(isMinLg);
  const onToggleClickTravelQueue = () => {
    setMinimizedTravelQueue((value) => !value);
    if (!isMinLg) {
      setMinimizedConstructionQueue(true);
      setMinimizedPopulation(true);
    }
  };
  const [minimizedConstructionQueue, setMinimizedConstructionQueue] = useState(
    isMinLg
  );
  const onToggleClickConstructionQueue = () => {
    setMinimizedConstructionQueue((value) => !value);
    if (!isMinLg) {
      setMinimizedTravelQueue(true);
      setMinimizedPopulation(true);
    }
  };
  const [minimizedPopulation, setMinimizedPopulation] = useState(isMinLg);
  const onToggleClickMinimizedPopulation = () => {
    setMinimizedPopulation((value) => !value);
    if (!isMinLg) {
      setMinimizedTravelQueue(true);
      setMinimizedConstructionQueue(true);
    }
  };

  useEffect(() => {
    setMinimizedTravelQueue(!isMinLg);
    setMinimizedConstructionQueue(!isMinLg);
    setMinimizedPopulation(!isMinLg);
  }, [
    isMinLg,
    setMinimizedTravelQueue,
    setMinimizedConstructionQueue,
    setMinimizedPopulation,
  ]);

  useEffect(() => {
    if (!ref.current) {
      return;
    }

    const lots = Array.from(ref.current.querySelectorAll("[data-type=lot]"));
    const handler = (e: Event) => {
      e.preventDefault();
      if (e.currentTarget instanceof SVGElement) {
        onLotClick(e.currentTarget.id);
      }
    };
    lots.forEach((lot) => {
      lot.addEventListener("click", handler);
    });

    return () => {
      lots.forEach((lot) => lot.removeEventListener("click", handler));
    };
  }, [ref.current]);

  const [hasSeenHelpPage] = useLocalStorage(
    "hasSeenHelpPage",
    false
  );

  if (hasSeenHelpPage === false) {
    return <Navigate to="/help" replace />;
  }

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
      <h2>Town</h2>
      <div
        className={classnames(
          "relative",
          "lg:w-9/12",
          "w-11/12",
          "max-w-screen-lg"
        )}
      >
        <TownSvg
          ref={ref}
          className={classnames("w-full", "h-auto", classes.svg as TArg)}
          lots={lots}
          constructionQueue={constructionQueue}
        />
        <div
          className={classnames(
            "absolute",
            "top-0",
            "left-0",
            "w-full",
            "flex",
            "justify-between",
            "items-start",
            "pointer-events-none",
            "gap-5",
          )}
        >
          {isMinSm ? (
            <>
              <ConstructionQueue
                minimized={minimizedConstructionQueue}
                onToggleClick={onToggleClickConstructionQueue}
              />
              <TravelQueue
                minimized={minimizedTravelQueue}
                onToggleClick={onToggleClickTravelQueue}
              />
              <Population
                minimized={minimizedPopulation}
                onToggleClick={onToggleClickMinimizedPopulation}
              />
            </>
          ) : (
            <div className={classnames("absolute", "right-0")}>
              <TownExpandMenu />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default Town;
