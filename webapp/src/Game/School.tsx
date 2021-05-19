import { formatDistanceToNow } from "date-fns";
import JSBI from "jsbi";
import React, { useState } from "react";
import { useInterval } from "react-use";
import { classnames } from "tailwindcss-classnames";
import { Education } from "../generated/education";
import { useStore } from "../store";
import styles from "../styles";
import PlaceholderImage from "./PlaceholderImage";
import { ReactComponent as ChefSvg } from "../../images/chef.svg";
import { ReactComponent as SalesmouseSvg } from "../../images/salesmouse.svg";
import { ReactComponent as SecuritySvg } from "../../images/security.svg";
import { ReactComponent as ThiefSvg } from "../../images/thief.svg";
import { ReactComponent as SvgSchool } from "../../images/school.svg";
import { formatDurationShort } from "../utils";

const title = classnames("text-lg", "md:text-xl", "mb-2");
const label = classnames("text-xs", "md:text-sm");
const value = classnames("text-sm");
const descriptionStyle = classnames("text-sm", "text-gray-600");

const svgs: Record<number, React.VFC | undefined> = {
  [Education.CHEF]: ChefSvg,
  [Education.SALESMOUSE]: SalesmouseSvg,
  [Education.GUARD]: SecuritySvg,
  [Education.THIEF]: ThiefSvg,
};

const descriptions: Record<number, React.VFC | undefined> = {
  [Education.CHEF]: () => (
    <p className={descriptionStyle}>Chefs work in kitchens to make pizzas.</p>
  ),
  [Education.SALESMOUSE]: () => (
    <p className={descriptionStyle}>Salesmice work in shops to sell pizzas.</p>
  ),
  [Education.GUARD]: () => (
    <p className={descriptionStyle}>
      Security guards help protect your town against thieves.
    </p>
  ),
  [Education.THIEF]: () => (
    <p className={descriptionStyle}>
      Thieves can be sent to other towns to steal coins.
    </p>
  ),
};

const numberFormat = new Intl.NumberFormat();

function School() {
  const educations = useStore((state) => state.gameData?.educations) || [];
  const trainingQueue = useStore((state) => state.gameState.trainingQueue);
  const train = useStore((state) => state.train);

  const onTrainClick = (e: React.MouseEvent, education: Education) => {
    e.preventDefault();
    train(education, 1);
    window.scroll(0, 0);
  };

  const [now, setNow] = useState(Date.now());
  useInterval(() => {
    setNow(Date.now());
  }, 1000);

  return (
    <div className={classnames("max-w-full", "px-2")}>
      <h2>School</h2>
      <SvgSchool height={100} width={100} />
      <p className={classnames("my-4", "text-gray-700")}>
        Educate your mice so that they can start contributing to your town.
      </p>

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
                      Number(
                        JSBI.divide(
                          JSBI.BigInt(training.completeAt),
                          JSBI.BigInt(1e6)
                        )
                      ),
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
          const SvgImage = svgs[eduKey];
          const Description = descriptions[eduKey];
          return (
            <div className={classnames("flex", "mb-8")} key={education.title}>
              <div style={{ width: 110 }}>
                {SvgImage ? <SvgImage /> : <PlaceholderImage />}
              </div>
              <div className={classnames("ml-4")}>
                <div className={title}>{education.title}</div>
                {Description && <Description />}
                <table>
                  <tbody>
                    <tr>
                      <td className={classnames("px-2")}>
                        <span className={label}>Train time:</span>
                      </td>
                      <td className={classnames("px-2")}>
                        <span className={value}>
                          {formatDurationShort(education.trainTime)}
                        </span>
                      </td>
                    </tr>
                    {education.cost && (
                      <tr>
                        <td className={classnames("px-2")}>
                          <span className={label}>Cost:</span>
                        </td>
                        <td className={classnames("px-2")}>
                          <span className={value}>
                            {numberFormat.format(education.cost)} coins
                          </span>
                        </td>
                      </tr>
                    )}
                  </tbody>
                </table>
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
