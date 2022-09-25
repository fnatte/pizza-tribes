import React from "react";
import { Building } from "../../generated/building";
import shopImg from "images/v2/shop.png";
import kitchenImg from "images/v2/kitchen.png";
import researchInstituteImg from "images/v2/research-institute.png";
import schoolImg from "images/v2/school.png";
import townCentreImg from "images/v2/town-centre.png";
import houseImg from "images/v2/house.png";
import marketingHqImg from "images/v2/marketing-hq.png";

function getBuildingImgSrc(building: Building): string {
  switch (building) {
    case Building.SHOP:
      return shopImg;
    case Building.HOUSE:
      return houseImg;
    case Building.SCHOOL:
      return schoolImg;
    case Building.KITCHEN:
      return kitchenImg;
    case Building.TOWN_CENTRE:
      return townCentreImg;
    case Building.MARKETINGHQ:
      return marketingHqImg;
    case Building.RESEARCH_INSTITUTE:
      return researchInstituteImg;
  }
}

function BuildingImage({
  building,
  ...props
}: {
  building: Building;
} & JSX.IntrinsicElements["img"]) {
  return <img {...props} src={getBuildingImgSrc(building)} />;
}

export default BuildingImage;
