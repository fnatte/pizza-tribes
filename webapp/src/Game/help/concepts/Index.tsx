import React from "react";
import { Link } from "react-router-dom";

export function Concepts() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Concepts</h3>
      <ul className="pl-4">
        <li>
          <Link to="demand">Demand</Link>
        </li>
        <li>
          <Link to="rush-hour">Rush Hour</Link>
        </li>
      </ul>
    </article>
  );
}
