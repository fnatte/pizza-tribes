import React from "react";
import { ReactComponent as SvgResearchInstitute } from "images/research-institute.svg";

export function ResearchInstitute() {
  return (
    <article className="prose prose-gray p-4">
      <h3>ResearchInstitute</h3>
      <SvgResearchInstitute
        width={100}
        height={100}
        className="float-left mr-4"
      />
      <p>The research institute is used for... *drumroll* &mdash; research!</p>
    </article>
  );
}
