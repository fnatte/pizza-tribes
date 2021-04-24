import React from "react";
import { useEffect, useRef } from "react";
import {useNavigate} from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import { ReactComponent as TownSvg } from "../../images/town.svg";
import classes from "./town.module.css";

function Town() {
  const ref = useRef<SVGImageElement>(null);
  const navigate = useNavigate();

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
      <TownSvg ref={ref} className={classes.svg} />
    </div>
  );
}

export default Town;
