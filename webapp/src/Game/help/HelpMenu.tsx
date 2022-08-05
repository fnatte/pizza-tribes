import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import classnames from "classnames";

export function HelpMenu({
  className,
  alwaysExpanded,
}: {
  className: string;
  alwaysExpanded: boolean;
}) {
  const [isExpanded, setIsExpanded] = useState(false);

  useEffect(() => {
    setIsExpanded(alwaysExpanded);
  }, [alwaysExpanded]);

  return (
    <div className={classnames("bg-green-50 py-2 px-8", className)}>
      {!alwaysExpanded && (
        <button
          className="block mx-auto"
          aria-expanded={isExpanded}
          onClick={() => setIsExpanded((x) => !x)}
        >
          {isExpanded ? "➖" : "➕"} Help Menu
        </button>
      )}
      {(isExpanded || alwaysExpanded) && (
        <>
          {!alwaysExpanded && <hr className="my-2" />}
          <nav>
            <ul className="leading-loose">
              <li className="mb-2">
                <Link to="getting-started">Getting Started</Link>
                <ul className="pl-4">
                  <li>
                    <Link to="getting-started#introduction">Introduction</Link>
                  </li>
                  <li>
                    <Link to="getting-started#constructing-buildings">
                      Constructing Buildings
                    </Link>
                  </li>
                  <li>
                    <Link to="getting-started#educating-mice">
                      Educating Mice
                    </Link>
                  </li>
                </ul>
              </li>

              <li className="mb-2">
                <Link to="buildings">Buildings</Link>
                <ul className="pl-4">
                  <li>
                    <Link to="buildings/town-centre">Town Centre</Link>
                  </li>
                  <li>
                    <Link to="buildings/kitchen">Kitchen</Link>
                  </li>
                  <li>
                    <Link to="buildings/shop">Shop</Link>
                  </li>
                  <li>
                    <Link to="buildings/house">House</Link>
                  </li>
                  <li>
                    <Link to="buildings/school">School</Link>
                  </li>
                  <li>
                    <Link to="buildings/marketing-hq">Marketing HQ</Link>
                  </li>
                  <li>
                    <Link to="buildings/research-institute">
                      Research Institute
                    </Link>
                  </li>
                </ul>
              </li>

              <li className="mb-2">
                <Link to="educations">Educations</Link>
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
              </li>

              <li className="mb-2">
                <Link to="concepts">Concepts</Link>
                <ul className="pl-4">
                  <li>
                    <Link to="concepts/demand">Demand</Link>
                  </li>
                  <li>
                    <Link to="concepts/rush-hour">Rush Hour</Link>
                  </li>
                </ul>
              </li>
            </ul>
          </nav>
        </>
      )}
    </div>
  );
}
