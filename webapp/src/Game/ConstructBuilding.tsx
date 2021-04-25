import React from "react";
import {useNavigate} from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import {useStore} from "../store";
import styles from "../styles";

const title = classnames("text-xl", "mb-2");
const label = classnames("text-sm", "mr-1");
const value = classnames("text-lg", "ml-1");

type Props = {
  lotId: string;
};

type BuildingInfo = {
  id: string;
  title: string;
  cost: number;
  buildTime: string;
};

const buildings: Array<BuildingInfo> = [
  {
    id: "kitchen",
    title: "Kitchen",
    cost: 300,
    buildTime: "2 work minutes",
  },
  {
    id: "shop",
    title: "Shop",
    cost: 100,
    buildTime: "2 work minutes",
  },
  {
    id: "house",
    title: "House",
    cost: 200,
    buildTime: "2 work minutes",
  },
];

const PlaceholderImage = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="300"
    height="150"
    viewBox="0 0 300 150"
  >
    <rect fill="#ddd" width="300" height="150" />
    <text
      fill="rgba(0,0,0,0.5)"
      fontFamily="sans-serif"
      fontSize="30"
      dy="10.5"
      fontWeight="bold"
      x="50%"
      y="50%"
      textAnchor="middle"
    >
      300×150
    </text>
  </svg>
);

const ConstructBuilding = ({ lotId }: Props) => {
  const constructBuilding = useStore((state) => state.constructBuilding);
  const navigate = useNavigate();

  const onSelectClick = (e: React.MouseEvent, building: string) => {
    e.preventDefault();
    constructBuilding(lotId, building);
    navigate("/town");
  };

  return (
    <div className={classnames("container", "mx-auto", "mt-4")}>
      <h2>Construct Building</h2>
      {buildings.map((building) => {
        return (
          <div className={classnames("flex", "mb-8")} key={building.id}>
            <PlaceholderImage />
            <div className={classnames("ml-4")} key={building.id}>
              <div className={title}>{building.title}</div>
              <div
                className={classnames("grid", "grid-cols-2", "items-center")}
              >
                <span className={label}>Cost:</span>
                <span className={value}>{building.cost} coins</span>
                <span className={label}>Build time:</span>
                <span className={value}>{building.buildTime}</span>
              </div>
              <div className={classnames("my-2")}>
                <button
                  className={classnames(styles.button)}
                  onClick={(e) => onSelectClick(e, building.id)}
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

