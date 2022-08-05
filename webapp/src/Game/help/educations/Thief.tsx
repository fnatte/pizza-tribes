import React from "react";
import { ReactComponent as SvgThief } from "images/thief.svg";
import { Link } from "react-router-dom";

export function Thief() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Thief</h3>
      <SvgThief width={150} className="float-left mr-4" />
      <p>Thieves can be sent to other towns to steal their coins!</p>
      <p>
        When thieves are sent on a heist they might be caught by{" "}
        <Link to="../educations/guard">guards</Link>. Any caught thief will be
        demoted from their thief education. You can increase your chances of a
        succesfull heist by outnumbering the guards.
      </p>
    </article>
  );
}
