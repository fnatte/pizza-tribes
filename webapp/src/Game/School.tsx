import React from "react";
import {useNavigate} from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import {useStore} from "../store";
import styles from "../styles";
import PlaceholderImage from "./PlaceholderImage";

const title = classnames("text-xl", "mb-2");
const label = classnames("text-sm", "mr-1");
const value = classnames("text-lg", "ml-1");

type RoleInfo = {
  id: string;
  title: string;
  trainTime: string;
};

const roles: Array<RoleInfo> = [
  {
    id: "chef",
    title: "Chef",
    trainTime: "2 minutes",
  },
  {
    id: "salesmouse",
    title: "Salesmouse",
    trainTime: "2 minutes",
  },
  {
    id: "guard",
    title: "Guard",
    trainTime: "2 minutes",
  },
  {
    id: "thief",
    title: "Thief",
    trainTime: "10 minutes",
  },
];

function School() {
  const train = useStore((state) => state.train);
  const navigate = useNavigate();

  const onTrainClick = (e: React.MouseEvent, roleId: string) => {
    e.preventDefault();
    console.log("train ", roleId);
    train(roleId, 1);
    navigate("/town");
  };

  return (
    <div>
      <h2>School</h2>
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
                  onClick={(e) => onTrainClick(e, role.id)}
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
