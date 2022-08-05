import React from "react";
import classnames from "classnames";
import { Link } from "react-router-dom";
import imgStats from "./demand-stats.png";

export function Demand() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Demand</h3>
      <p>
        Demand is a key component in selling pizzas. Your town can not sell more
        than than there is demand, no matter how many salesmice. You can
        increase the demand for your pizzas by{" "}
        <Link to="../../educations/publicist">publicists</Link> working in{" "}
        <Link to="../../buildings/marketing-hq">Marketing HQ</Link>.
      </p>
      <p>
        To find out your demand, see the <em>Stats page</em>:
        <img
          src={imgStats}
          alt="Screenshot of stats"
          className={classnames(
            "border-2",
            "border-green-300",
            "filter",
            "hue-rotate-15",
            "my-2"
          )}
        />
        On the stats page and as shown in the image above, the demand is
        specified both offpeak and during{" "}
        <Link to="../rush-hour">rush hour</Link>.
      </p>
    </article>
  );
}
