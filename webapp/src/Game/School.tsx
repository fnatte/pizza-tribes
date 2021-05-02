import {
  formatDistanceToNow,
  formatDuration,
  intervalToDuration,
} from "date-fns";
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useInterval } from "react-use";
import { classnames } from "tailwindcss-classnames";
import { Education } from "../generated/education";
import { useStore } from "../store";
import styles from "../styles";
import PlaceholderImage from "./PlaceholderImage";

const title = classnames("text-xl", "mb-2");
const label = classnames("text-sm", "mr-1");
const value = classnames("text-lg", "ml-1");

function School() {
  const educations = useStore((state) => state.gameData?.educations) || [];
  const trainingQueue = useStore((state) => state.gameState.trainingQueue);
  const train = useStore((state) => state.train);
  const navigate = useNavigate();

  const onTrainClick = (e: React.MouseEvent, education: Education) => {
    e.preventDefault();
    train(education, 1);
    navigate("/town");
  };

  const [now, setNow] = useState(Date.now());
  useInterval(() => {
    setNow(Date.now());
  }, 1000);

  return (
    <div>
      <h2>School</h2>

      {trainingQueue.length > 0 && (
        <>
          <h3>In Training</h3>
          <table>
            <tbody>
              {trainingQueue.map((training) => (
                <tr
                  key={
                    training.completeAt.toString() +
                    training.education +
                    training.amount
                  }
                >
                  <td className={classnames("p-2")}>{training.amount}</td>
                  <td className={classnames("p-2")}>
                    {educations[training.education].title}
                  </td>
                  <td className={classnames("p-2")}>
                    {formatDistanceToNow(
                      Number(training.completeAt / BigInt(1e6)),
                      {
                        includeSeconds: true,
                        addSuffix: true,
                      }
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}

      <h3>Train</h3>
      {Object.keys(educations)
        .map(Number)
        .map((eduKey) => {
          const education = educations[eduKey];
          return (
            <div className={classnames("flex", "mb-8")} key={education.title}>
              <PlaceholderImage />
              <div className={classnames("ml-4")}>
                <div className={title}>{education.title}</div>
                <div
                  className={classnames("grid", "grid-cols-2", "items-center")}
                >
                  <span className={label}>Train time:</span>
                  <span className={value}>
                    {formatDuration(
                      intervalToDuration({
                        start: 0,
                        end: education.trainTime * 1000,
                      }),
                      { delimiter: ", " }
                    )}
                  </span>
                </div>
                <div className={classnames("my-2")}>
                  <button
                    className={classnames(styles.button)}
                    onClick={(e) => onTrainClick(e, eduKey)}
                  >
                    Train
                  </button>
                </div>
              </div>
            </div>
          );
        })}
    </div>
  );
}

export default School;
