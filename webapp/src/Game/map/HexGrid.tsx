import React from "react";
import classNames from "classnames";
import { ReactComponent as TownSvg } from "../../../images/town.svg";

import style from "./hexgrid.module.css";
import { WorldEntry } from "../../generated/world";

type Props = {
  x: number;
  y: number;
  size: number;
  data: Map<string, WorldEntry|null>;
  onNavigate: (x: number, y: number) => void;
  onClick: (x: number, y: number) => void;
};

const range = (from: number, to: number) => [
  ...[...Array(to - from).keys()].map((i) => from + i),
];

const getMapKey = (x: number, y: number) => `${x}:${y}`;

const Arrow: React.FC<{
  dir: "top" | "left" | "bottom" | "right";
  onNavigate: (x: number, y: number) => void;
}> = ({ dir, onNavigate }) => {
  const onClick = () => {
    switch (dir) {
      case "top":
        onNavigate(0, -1);
        break;
      case "right":
        onNavigate(1, 0);
        break;
      case "bottom":
        onNavigate(0, 1);
        break;
      case "left":
        onNavigate(-1, 0);
        break;
    }
  };
  return (
    <div
      className={`${style.arrow} ${style["arrow" + dir]}`}
      onClick={onClick}
    />
  );
};

const HexGrid: React.FC<Props> = ({
  data,
  x,
  y,
  size,
  onNavigate,
  onClick,
}) => {
  const sh2floor = Math.floor(size / 2);
  const sh2round = Math.round(size / 2);
  return (
    <div className={`${style.root} ${style[`size-${size}`]}`}>
      {range(y - sh2floor, y + sh2round).map((y) => (
        <div className={style.row} key={y}>
          {range(x - sh2floor, x + sh2round).map((x) => {
            const entry = data.get(getMapKey(x, y));
            const objectType = entry?.object?.oneofKind;
            return (
              <div
                className={classNames(
                  style.tile,
                  style.grass,
                  objectType && style[objectType]
                )}
                key={x}
                onClick={() => onClick(x, y)}
              >
                {objectType === "town" && (
                  <TownSvg style={{ width: "100%", height: "100%" }} />
                )}
              </div>
            );
          })}
        </div>
      ))}
      <Arrow dir="top" onNavigate={onNavigate} />
      <Arrow dir="right" onNavigate={onNavigate} />
      <Arrow dir="bottom" onNavigate={onNavigate} />
      <Arrow dir="left" onNavigate={onNavigate} />
    </div>
  );
};

export default HexGrid;
