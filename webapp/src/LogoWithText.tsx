import React from "react";
import { Logo } from "./Logo";

export function LogoWithText() {
  return (
    <div className="flex items-center">
      <Logo className="w-40" />
      <div className="text-center mt-4">
        <div className="text-4xl">Pizza</div>
        <div className="text-5xl -mt-2">Tribes</div>
      </div>
    </div>
  );
}
