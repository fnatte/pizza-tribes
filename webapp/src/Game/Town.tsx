import React from "react";
import { useEffect, useRef } from "react";
import {useNavigate} from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import {useStore} from "../store";
import classes from "./town.module.css";
import TownSvg from "./TownSvg";

function Town() {
  const ref = useRef<SVGSVGElement>(null);
  const navigate = useNavigate();
  const lots = useStore(state => state.gameState.lots);

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
        <TownSvg ref={ref} className={classes.svg} lots={lots} />
    </div>
  );
}

export default Town;
