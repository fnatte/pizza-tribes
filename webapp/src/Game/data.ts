import {Building} from "../generated/building";
import {GameState_Population} from "../generated/gamestate";

export type RoleInfo = {
  id: keyof GameState_Population;
  title: string;
  titlePlural: string;
  employer?: Building;
};

export type BuildingInfo = {
  id: string;
  title: string;
  cost: number;
  buildTime: string;
  building: Building;
};


export const roles: Array<RoleInfo> = [
  {
    id: "uneducated",
    title: "Uneducated",
    titlePlural: "Uneducated",
  },
  {
    id: "chefs",
    title: "Chef",
    titlePlural: "Chefs",
    employer: Building.KITCHEN,
  },
  {
    id: "salesmice",
    title: "Salesmouse",
    titlePlural: "Salesmice",
    employer: Building.SHOP,
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

export const buildings: Array<BuildingInfo> = [
  {
    id: "kitchen",
    title: "Kitchen",
    cost: 300,
    buildTime: "2 work minutes",
    building: Building.KITCHEN,
  },
  {
    id: "shop",
    title: "Shop",
    cost: 100,
    buildTime: "2 work minutes",
    building: Building.SHOP,
  },
  {
    id: "house",
    title: "House",
    cost: 200,
    buildTime: "2 work minutes",
    building: Building.HOUSE,
  },
  {
    id: "school",
    title: "School",
    cost: 500,
    buildTime: "10 work minutes",
    building: Building.SCHOOL,
  },
];

export const buildingsByType = buildings.reduce<Record<number, BuildingInfo>>((res, building) => {
  res[building.building] = building;
  return res;
}, {});

