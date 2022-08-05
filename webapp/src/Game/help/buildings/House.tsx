import React from "react";
import { ReactComponent as SvgHouse } from "images/house.svg";

export function House() {
  return (
    <article className="prose prose-gray p-4">
      <h3>House</h3>
      <SvgHouse width={100} height={100} className="float-left mr-4" />
      <p>
        Mice live in houses. Build more houses to get more mice in your tribe.
      </p>
    </article>
  );
}
