import React from "react";
import { ReactComponent as SvgPublicist } from "images/publicist.svg";
import { Link } from "react-router-dom";

export function Publicist() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Publicist</h3>
      <SvgPublicist width={150} className="float-left mr-4" />
      <p>
        Publicists work at{" "}
        <Link to="../buildings/marketing-hq">Marketing HQ</Link> to increase
        demand for your pizzas.
      </p>
    </article>
  );
}
