import { useCallback } from "react";
import { useNavigate, useParams } from "react-router";

export type Page =
  | "town"
  | "town-lot"
  | "mouse"
  | "mouse-appearance"
  | "world"
  | "world-entry";
type PageParams<T> = T extends "town-lot"
  ? [string]
  : T extends "mouse"
  ? [string]
  : T extends "mouse-appearance"
  ? [string]
  : T extends "world-entry"
  ? [number, number]
  : [];

export function useGameNavigate() {
  const navigate = useNavigate();
  const { gameId } = useParams();

  const gameNavigate = useCallback(
    <T extends Page, P extends PageParams<T>>(page: T, ...params: P) => {
      const prefix = `/game/${gameId}`;
      switch (page) {
        case "town":
          navigate(`${prefix}/town`);
          return;
        case "world":
          navigate(`${prefix}/world`);
          return;
        case "world-entry":
          navigate(`${prefix}/world/entry?x=${params[0]}&y=${params[1]}`);
          return;
        case "town-lot":
          navigate(`${prefix}/town/${params[0]}`);
          return;
        case "mouse":
          navigate(`${prefix}/mouse/${params[0]}`);
          return;
        case "mouse-appearance":
          navigate(`${prefix}/mouse/${params[0]}/appearance`);
          return;
      }
    },
    [navigate]
  );

  return gameNavigate;
}
