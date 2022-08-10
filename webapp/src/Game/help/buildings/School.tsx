import React from "react";
import { ReactComponent as SvgSchool } from "images/school.svg";
import { LevelInfoTable } from "./LevelInfoTable";
import { Building } from "../../../generated/building";

export function School() {
  return (
    <article className="prose prose-gray p-4">
      <h3>School</h3>
      <SvgSchool width={100} height={100} className="float-left mr-4" />
      <p>
        The school is used to educate mice. When a house is built, uneducated
        mice will move into that house. You should educate the mices in your
        tribe so that they can work.
      </p>
      <LevelInfoTable className="my-8 clear-left" building={Building.SCHOOL} />
    </article>
  );
}
