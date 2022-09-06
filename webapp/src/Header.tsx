import React from "react";
import classnames from "classnames";
import { LogoWithText } from "./LogoWithText";

function Header() {
  return (
    <div
      className={classnames(
        "flex",
        "justify-center",
        "flex-col",
        "items-center"
      )}
    >
      <h1 className="flex justify-center mt-2">
        <LogoWithText />
      </h1>
    </div>
  );
}

export default Header;
