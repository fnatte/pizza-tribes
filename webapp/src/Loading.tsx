import React from 'react';
import classnames from "classnames";
import { Hearts } from "./icons";

export function Loading() {
  return (
    <div
      className={classnames(
        "fixed",
        "left-1/2",
        "top-1/2",
        "-translate-y-1/2",
        "-translate-x-1/2"
      )}
    >
      <Hearts />
    </div>
  );
}
