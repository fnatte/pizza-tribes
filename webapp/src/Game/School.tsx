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

type RoleInfo = {
  id: string;
  title: string;
  trainTime: string;
  education: Education;
};

const roles: Array<RoleInfo> = [
  {
    id: "chef",
    title: "Chef",
    trainTime: "2 minutes",
    education: Education.CHEF,
  },
  {
    id: "salesmouse",
    title: "Salesmouse",
    trainTime: "2 minutes",
    education: Education.SALESMOUSE,
  },
  {
    id: "guard",
    title: "Guard",
    trainTime: "2 minutes",
    education: Education.GUARD,
  },
  {
    id: "thief",
    title: "Thief",
    trainTime: "10 minutes",
    education: Education.THIEF,
  },
];

const rolesByEducation = roles.reduce<Record<number, RoleInfo>>((res, role) => {
  res[role.education] = role;
  return res;
}, {});

function School() {
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
                    {rolesByEducation[training.education].title}
                  </td>
                  <td className={classnames("p-2")}>
                    {`in ${
                      training.completeAt / BigInt(1e9) -
                      BigInt(now) / BigInt(1e3)
                    }s`}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}

      <h3>Train</h3>
      {roles.map((role) => {
        return (
          <div className={classnames("flex", "mb-8")} key={role.id}>
            <PlaceholderImage />
            <div className={classnames("ml-4")}>
              <div className={title}>{role.title}</div>
              <div
                className={classnames("grid", "grid-cols-2", "items-center")}
              >
                <span className={label}>Train time:</span>
                <span className={value}>{role.trainTime}</span>
              </div>
              <div className={classnames("my-2")}>
                <button
                  className={classnames(styles.button)}
                  onClick={(e) => onTrainClick(e, role.education)}
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
