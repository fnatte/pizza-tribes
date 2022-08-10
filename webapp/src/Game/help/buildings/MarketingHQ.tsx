import React from "react";
import { ReactComponent as SvgMarketingHq } from "images/marketing-hq.svg";
import { LevelInfoTable } from "./LevelInfoTable";
import { Building } from "../../../generated/building";

export function MarketingHq() {
  return (
    <article className="prose prose-gray p-4">
      <h3>MarketingHq</h3>
      <SvgMarketingHq width={100} height={100} className="float-left mr-4" />
      <p>
        The marketing HQ is used to increase publicity for your pizzas &mdash;
        when publicity is increase, so is demand. Mice that are educated as{" "}
        <em>Publicists</em> will work here.
      </p>
      <LevelInfoTable className="my-8 clear-left" building={Building.MARKETINGHQ} />
    </article>
  );
}
