import React from "react";
import { ReactComponent as SvgChef } from "images/chef.svg";

export function Chef() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Chef</h3>
      <SvgChef width={150} className="float-left mr-4" />
      <p>
        Chefs work in kitchens to make pizzas.
      </p>
    </article>
  );
}
