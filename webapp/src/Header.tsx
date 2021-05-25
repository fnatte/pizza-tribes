import React from "react";
import { classnames } from "tailwindcss-classnames";
import { ReactComponent as SvgLogoWide } from "../images/logo-wide.svg";

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
      <h1 className={classnames("flex", "justify-center")}>
        <SvgLogoWide />
      </h1>
    </div>
  );
}

export default Header;
