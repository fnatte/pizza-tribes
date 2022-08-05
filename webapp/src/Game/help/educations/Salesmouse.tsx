import React from "react";
import { ReactComponent as SvgSalesmouse } from "images/salesmouse.svg";
import { Link } from "react-router-dom";

export function Salesmouse() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Salesmouse</h3>
      <SvgSalesmouse width={150} className="float-left mr-4" />
      <p>
        Salesmice work in <Link to="../buildings/shop">shops</Link> to sell
        pizzas. A salesmouse (without additional aid) can sell 1 pizza every 2
        seconds. However, your salesmice cannot sell more pizzas than the{" "}
        <Link to="../concepts/demand">demand</Link>.
      </p>
    </article>
  );
}
