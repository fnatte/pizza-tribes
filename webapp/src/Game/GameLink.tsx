import { trimStart } from "lodash";
import React from "react";
import { Link, LinkProps, useParams } from "react-router-dom";

export function GameLink({ to, ...props }: LinkProps) {
  const { gameId } = useParams();
  console.log(to)
  return (
    <Link
      {...props}
      to={`/game/${gameId}/${typeof to === "string" ? trimStart(to, "/") : to}`}
    />
  );
}
