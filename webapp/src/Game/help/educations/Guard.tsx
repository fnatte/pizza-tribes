import React from "react";
import { ReactComponent as SvgGuard } from "images/guard.svg";
import { Link } from "react-router-dom";

export function Guard() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Guard</h3>
      <SvgGuard width={150} className="float-left mr-4" />
      <p>
        Guards protect your town against{" "}
        <Link to="../educations/thief">thieves</Link>, preventing them from
        stealing your coins.
      </p>
      <p>
        If there is a heist on your town your guard have the chance of catching
        theives and demote them of their thief education. However, if a guard is
        found sleeping during a heist they will also be degraded from their
        guard education. You can increase your chances of a succesfull defence
        by outnumbering the theives.
      </p>
    </article>
  );
}
