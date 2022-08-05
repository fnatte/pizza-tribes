import React from "react";
import { ReactComponent as SvgTownCentre } from "images/town-centre.svg";

export function TownCentre() {
  return (
    <article className="prose prose-gray p-4">
      <h3>TownCentre</h3>
      <SvgTownCentre width={100} height={100} className="float-left mr-4" />
      <p>
        The town centre is your main building.
      </p>
    </article>
  );
}
