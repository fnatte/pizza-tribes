import React from "react";
import { LevelInfoTable } from "./LevelInfoTable";
import { Building } from "../../../generated/building";
import { Link } from "react-router-dom";
import BuildingImage from "../../components/BuildingImage";

export function MarketingHq() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Marketing HQ</h3>
      <BuildingImage
        building={Building.HOUSE}
        width={100}
        height={100}
        className="float-left mr-4"
      />
      <p>
        The marketing HQ is used to increase publicity for your pizzas &mdash;
        when publicity is increase, so is demand. Mice that are educated as{" "}
        <Link to="../educations/publicist">publicists</Link> will work here.
      </p>
      <LevelInfoTable
        className="my-8 clear-left"
        building={Building.MARKETINGHQ}
      />
    </article>
  );
}
