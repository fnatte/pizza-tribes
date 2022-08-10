import React from "react";
import { ReactComponent as SvgSalesmouse } from "images/salesmouse.svg";

export function Salesmouse() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Salesmouse</h3>
      <SvgSalesmouse width={150} className="float-left mr-4" />
      <p>
        Salesmice work in shops to sell pizzas.
      </p>
    </article>
  );
}
