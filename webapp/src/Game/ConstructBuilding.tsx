import React from "react";
import {useNavigate} from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import {Building} from "../generated/building";
import {useStore} from "../store";
import styles from "../styles";
import PlaceholderImage from "./PlaceholderImage";

const title = classnames("text-xl", "mb-2");
const label = classnames("text-sm", "mr-1");
const value = classnames("text-lg", "ml-1");

type Props = {
  lotId: string;
};

const ConstructBuilding = ({ lotId }: Props) => {
  const constructBuilding = useStore((state) => state.constructBuilding);
  const buildings = useStore((state) => state.gameData?.buildings) ?? [];
  const navigate = useNavigate();

  const onSelectClick = (e: React.MouseEvent, building: Building) => {
    e.preventDefault();
    constructBuilding(lotId, building);
    navigate("/town");
  };

  return (
    <div className={classnames("container", "mx-auto", "mt-4")}>
      <h2>Construct Building</h2>
      {Object.keys(buildings).map(Number).map((id) => {
        return (
          <div className={classnames("flex", "mb-8")} key={id}>
            <PlaceholderImage />
            <div className={classnames("ml-4")}>
              <div className={title}>{buildings[id].title}</div>
              <div
                className={classnames("grid", "grid-cols-2", "items-center")}
              >
                <span className={label}>Cost:</span>
                <span className={value}>{buildings[id].cost} coins</span>
                <span className={label}>Build time:</span>
                <span className={value}>{buildings[id].constructionTime}s</span>
              </div>
              <div className={classnames("my-2")}>
                <button
                  className={classnames(styles.button)}
                  onClick={(e) => onSelectClick(e, id)}
                >
                  Place Building
                </button>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default ConstructBuilding;

