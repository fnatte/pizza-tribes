import React from "react";
import { LevelInfoTable } from "./LevelInfoTable";
import { Building } from "../../../generated/building";
import BuildingImage from "../../components/BuildingImage";

export function House() {
  return (
    <article className="prose prose-gray p-4">
      <h3>House</h3>
      <BuildingImage
        building={Building.HOUSE}
        width={100}
        height={100}
        className="float-left mr-4"
      />
      <p>
        Mice live in houses. Build more houses to get more mice in your tribe.
      </p>
      <LevelInfoTable className="my-8 clear-left" building={Building.HOUSE} />
    </article>
  );
}
