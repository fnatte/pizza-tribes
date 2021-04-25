import {GameState_Population} from "../generated/gamestate";

export type RoleInfo = {
  id: keyof GameState_Population;
  title: string;
  titlePlural: string;
};

export const roles: Array<RoleInfo> = [
  {
    id: "unemployed",
    title: "Unemployed",
    titlePlural: "Unemployees",
  },
  {
    id: "chefs",
    title: "Chef",
    titlePlural: "Chefs",
  },
  {
    id: "salesmice",
    title: "Salesmouse",
    titlePlural: "Salesmice",
  },
  {
    id: "guards",
    title: "Guard",
    titlePlural: "Guards",
  },
  {
    id: "thieves",
    title: "Thief",
    titlePlural: "Thieves",
  },
];

