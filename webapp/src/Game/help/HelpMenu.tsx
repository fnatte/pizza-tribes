import React, { useEffect, useState } from "react";

import { Link, LinkProps } from "react-router-dom";
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

  const linkProps: Partial<LinkProps> = {
    onClick: () => setIsExpanded(false),
  };

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
                <Link {...linkProps} to="getting-started">
                  Getting Started
                </Link>
                <ul className="pl-4">
                  <li>
                    <Link {...linkProps} to="getting-started#introduction">
                      Introduction
                    </Link>
                  </li>
                  <li>
                    <Link
                      {...linkProps}
                      to="getting-started#constructing-buildings"
                    >
                      Constructing Buildings
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="getting-started#educating-mice">
                      Educating Mice
                    </Link>
                  </li>
                </ul>
              </li>

              <li className="mb-2">
                <Link {...linkProps} to="buildings">
                  Buildings
                </Link>
                <ul className="pl-4">
                  <li>
                    <Link
                      {...linkProps}
                      to="buildings/town-centre"
                      onClick={() => setIsExpanded(false)}
                    >
                      Town Centre
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="buildings/kitchen">
                      Kitchen
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="buildings/shop">
                      Shop
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="buildings/house">
                      House
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="buildings/school">
                      School
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="buildings/marketing-hq">
                      Marketing HQ
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="buildings/research-institute">
                      Research Institute
                    </Link>
                  </li>
                </ul>
              </li>

              <li className="mb-2">
                <Link {...linkProps} to="educations">
                  Educations
                </Link>
                <ul className="pl-4">
                  <li>
                    <Link {...linkProps} to="educations/chef">
                      Chef
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="educations/salesmouse">
                      Salesmouse
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="educations/guard">
                      Guard
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="educations/thief">
                      Thief
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="educations/publicist">
                      Publicist
                    </Link>
                  </li>
                </ul>
              </li>

              <li className="mb-2">
                <Link {...linkProps} to="concepts">
                  Concepts
                </Link>
                <ul className="pl-4">
                  <li>
                    <Link {...linkProps} to="concepts/demand">
                      Demand
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="concepts/quality">
                      Quality
                    </Link>
                  </li>
                  <li>
                    <Link {...linkProps} to="concepts/rush-hour">
                      Rush Hour
                    </Link>
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
