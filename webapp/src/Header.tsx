import React from "react";
import classnames from "classnames";
import { LogoWithText } from "./LogoWithText";

function Header({ className }: { className?: string }) {
  return (
    <div
      className={classnames(
        "flex",
        "justify-center",
        "flex-col",
        "items-center",
        className
      )}
    >
      <h1 className="flex justify-center mt-2">
        <LogoWithText />
      </h1>
    </div>
  );
}

export default Header;
