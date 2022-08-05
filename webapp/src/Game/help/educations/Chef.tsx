import React from "react";
import { ReactComponent as SvgChef } from "images/chef.svg";
import { Link } from "react-router-dom";

export function Chef() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Chef</h3>
      <SvgChef width={150} className="float-left mr-4" />
      <p>
        Chefs work in <Link to="../buildings/kitchen">kitchens</Link> to make
        pizzas. A working chef in a simple kitchen makes 1 pizza every 5
        seconds. However, with a better kitchen it a chef can be even more
        efficient.
      </p>
    </article>
  );
}
