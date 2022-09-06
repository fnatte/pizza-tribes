import React from "react";
import logoImg from "images/logos/logo-1024.png";

export function Logo({ className }: { className?: string }) {
  return <img className={className} src={logoImg} />;
}
