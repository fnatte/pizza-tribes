import React from "react";
import { ReactComponent as SvgThief } from "images/thief.svg";

export function Thief() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Thief</h3>
      <SvgThief width={150} className="float-left mr-4" />
      <p>
        Thieves can be sent to other towns to steal their coins!
      </p>
    </article>
  );
}
