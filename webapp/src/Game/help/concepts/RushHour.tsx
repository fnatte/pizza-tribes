import React from "react";
import classnames from "classnames";
import imgRushHourIndicator from "./rush-hour-indicator.png";
import { Link } from "react-router-dom";

export function RushHour() {
  return (
    <article className="prose prose-gray p-4">
      <h3>Rush Hour</h3>
      <p>
        During every mouse day there are 5 rush hours:
        <ul>
          <li>Between 11 and 13</li>
          <li>Between 18 and 21</li>
        </ul>
        Note that this is not in your regular human time, but in mouse time.
      </p>
      <p>
        During rush hour the <Link to="../demand">demand for pizza</Link> is
        higher, allowing your town to sell even more pizzas if your salesmice
        can work fast enough.
      </p>
      <p>
        During rush hour, an indicator is shown by the side of the current mouse
        time as shown in the image below:
        <img
          src={imgRushHourIndicator}
          alt="Screenshot of rush hour"
          className={classnames(
            "border-2",
            "border-green-300",
            "filter",
            "hue-rotate-15",
            "my-2"
          )}
        />
      </p>
    </article>
  );
}
