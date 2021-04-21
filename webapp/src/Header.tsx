import React from "react";
import {classnames} from "tailwindcss-classnames";

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
      <h1 className={classnames("flex", "justify-center", "p-8", "text-4xl")}>
        Pizza Mouse
      </h1>
      <div className={classnames("text-2xl")}>🍕🍕🍕🍕</div>
    </div>
  );
}

export default Header;
