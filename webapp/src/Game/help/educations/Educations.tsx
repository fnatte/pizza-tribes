import React from "react";
import { Link } from "react-router-dom";

export function Educations() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Educations</h3>
      <ul className="pl-4">
        <li>
          <Link to="educations/chef">Chef</Link>
        </li>
        <li>
          <Link to="educations/salesmouse">Salesmouse</Link>
        </li>
        <li>
          <Link to="educations/guard">Guard</Link>
        </li>
        <li>
          <Link to="educations/thief">Thief</Link>
        </li>
        <li>
          <Link to="educations/publicist">Publicist</Link>
        </li>
      </ul>
    </article>
  );
}
