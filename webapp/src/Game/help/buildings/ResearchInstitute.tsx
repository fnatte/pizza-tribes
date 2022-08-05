import React from "react";
import { ReactComponent as SvgResearchInstitute } from "images/research-institute.svg";
import { LevelInfoTable } from "./LevelInfoTable";
import { Building } from "../../../generated/building";

export function ResearchInstitute() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Research Institute</h3>
      <SvgResearchInstitute
        width={100}
        height={100}
        className="float-left mr-4 mb-4"
      />
      <p>The research institute is used for... *drumroll* &mdash; research!</p>
      <LevelInfoTable
        className="my-8 clear-left"
        building={Building.RESEARCH_INSTITUTE}
      />
    </article>
  );
}
