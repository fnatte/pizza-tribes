import React from "react";
import { ReactComponent as SvgHouse } from "images/house.svg";
import { LevelInfoTable } from "./LevelInfoTable";
import { Building } from "../../../generated/building";

export function House() {
  return (
    <article className="prose prose-gray p-4">
      <h3>House</h3>
      <SvgHouse width={100} height={100} className="float-left mr-4 mb-4" />
      <p>
        Mice live in houses. Build more houses to get more mice in your tribe.
      </p>
      <LevelInfoTable className="my-8 clear-left" building={Building.HOUSE} />
    </article>
  );
}
