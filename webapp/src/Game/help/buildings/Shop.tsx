import React from "react";
import { LevelInfoTable } from "./LevelInfoTable";
import { Building } from "../../../generated/building";
import { Link } from "react-router-dom";
import BuildingImage from "../../components/BuildingImage";

export function Shop() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Shop</h3>
      <BuildingImage
        building={Building.SHOP}
        width={100}
        height={100}
        className="float-left mr-4"
      />
      <p>
        The shop is used to sell pizzas &mdash; and as such earning you coins.
        Mice that are educated as{" "}
        <Link to="../educations/salesmouse">salesmice</Link> will work here.
        Without salesmice there no pizzas will be sold.
      </p>
      <LevelInfoTable className="my-8 clear-left" building={Building.SHOP} />
    </article>
  );
}
