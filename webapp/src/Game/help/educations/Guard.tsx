import React from "react";
import { ReactComponent as SvgGuard } from "images/guard.svg";

export function Guard() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Guard</h3>
      <SvgGuard width={150} className="float-left mr-4" />
      <p>
        Guards protect your town against thieves, preventing them from stealing
        your coins.
      </p>
    </article>
  );
}
