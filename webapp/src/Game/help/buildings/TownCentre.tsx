import React from "react";
import { LevelInfoTable } from "./LevelInfoTable";
import { Building } from "../../../generated/building";
import BuildingImage from "../../components/BuildingImage";

export function TownCentre() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Town Centre</h3>
      <BuildingImage
        building={Building.TOWN_CENTRE}
        width={100}
        height={100}
        className="float-left mr-4"
      />
      <p>The town centre is your main building.</p>
      <LevelInfoTable
        className="my-8 clear-left"
        building={Building.TOWN_CENTRE}
      />
    </article>
  );
}
