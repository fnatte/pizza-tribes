import React, { useState } from "react";
import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { useMedia } from "react-use";
import classnames from "classnames";
import ConstructionQueue from "./ConstructionQueue";
import Population from "./Population";
import TownExpandMenu from "./TownExpandMenu";
import Town from "./town/Town";
import TravelQueue from "./TravelQueue";

function TownView() {
  const ref = useRef<SVGSVGElement>(null);
  const navigate = useNavigate();

  const onLotClick = (lotId: string) => {
    navigate(lotId);
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
        onLotClick(e.currentTarget.id.replace("lot", ""));
      }
    };
    lots.forEach((lot) => {
      lot.addEventListener("click", handler);
    });

    return () => {
      lots.forEach((lot) => lot.removeEventListener("click", handler));
    };
  }, [ref.current]);

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
      <h2 className="hidden lg:block">Town</h2>
      <div
        className={classnames(
          "mt-2",
          "relative",
          "lg:w-9/12",
          "md:w-11/12",
          "w-[98%]",
          "max-w-screen-lg"
        )}
      >
        <Town />
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
            "gap-5"
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

export default TownView;
