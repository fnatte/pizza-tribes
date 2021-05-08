import React from "react";
import { useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { classnames, TArg } from "tailwindcss-classnames";
import { useStore } from "../store";
import ConstructionQueue from "./ConstructionQueue";
import Population from "./Population";
import classes from "./town.module.css";
import TownSvg from "./TownSvg";
import TravelQueue from "./TravelQueue";

function Town() {
  const ref = useRef<SVGSVGElement>(null);
  const navigate = useNavigate();
  const lots = useStore((state) => state.gameState.lots);

  const onLotClick = (lotId: string) => {
    navigate(`/town/${lotId.replace("lot", "")}`);
  };

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
          className={classnames("w-full", classes.svg as TArg)}
          lots={lots}
          height={undefined}
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
            "pointer-events-none"
          )}
        >
          <ConstructionQueue />
          <TravelQueue />
          <Population />
        </div>
      </div>
    </div>
  );
}

export default Town;
