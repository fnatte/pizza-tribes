import React from "react";
import { ReactComponent as SvgKitchen } from "images/kitchen.svg";
import { LevelInfoTable } from "./LevelInfoTable";
import { Building } from "../../../generated/building";
import { Link } from "react-router-dom";

export function Kitchen() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Kitchen</h3>
      <SvgKitchen width={100} height={100} className="float-left mr-4" />
      <p>
        The kitchen is used to make pizzas. Mice that are educated as{" "}
        <Link to="../educations/chef">Chefs</Link> will work here. Without chefs
        there will be no pizzas.
      </p>
      <LevelInfoTable className="my-8 clear-left" building={Building.KITCHEN} />
    </article>
  );
}
