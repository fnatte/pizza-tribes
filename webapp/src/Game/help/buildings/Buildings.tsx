import React from "react";
import { Link } from "react-router-dom";

export function Buildings() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Buildings</h3>
      <ul className="pl-4">
        <li>
          <Link to="../buildings/town-centre">Town Centre</Link>
        </li>
        <li>
          <Link to="../buildings/kitchen">Kitchen</Link>
        </li>
        <li>
          <Link to="../buildings/shop">Shop</Link>
        </li>
        <li>
          <Link to="../buildings/house">House</Link>
        </li>
        <li>
          <Link to="../buildings/school">School</Link>
        </li>
        <li>
          <Link to="../buildings/marketing-hq">Marketing HQ</Link>
        </li>
        <li>
          <Link to="../buildings/research-institute">Research Institute</Link>
        </li>
      </ul>
    </article>
  );
}
