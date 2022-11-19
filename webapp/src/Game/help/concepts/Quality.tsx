import React from "react";
import { Link } from "react-router-dom";

export function Quality() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Quality</h3>
      <p>
        Quality is a key factor for <Link to="./demand">demand</Link>. If your
        quality is improved, your demand will rise and you can sell more pizzas
        (or increase the price). Quality is improved from{" "}
        <Link to="../../buildings/research-institute">research</Link>.
      </p>
    </article>
  );
}
