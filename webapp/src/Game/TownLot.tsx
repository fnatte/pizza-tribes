import React from "react";
import {useParams} from "react-router-dom";
import { classnames } from "tailwindcss-classnames";

function TownLot() {
  const { id } = useParams();

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
      <h2>Lot {id}</h2>
    </div>
  );
}

export default TownLot;
