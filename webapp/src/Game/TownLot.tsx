import React, { useCallback } from "react";
import { useParams } from "react-router-dom";
import { classnames } from "tailwindcss-classnames";
import { Building } from "../generated/building";
import { useStore } from "../store";
import ConstructBuilding from "./ConstructBuilding";
import School from "./School";
import { ReactComponent as SvgKitchen } from "../../images/kitchen.svg";
import { ReactComponent as SvgShop } from "../../images/shop.svg";
import { ReactComponent as SvgHouse } from "../../images/house.svg";

function TownLot() {
  const { id } = useParams();

  const lot = useStore(useCallback((state) => state.gameState.lots[id], [id]));
  const stats = useStore((state) => state.gameStats);
  const population = useStore((state) => state.gameState.population);

  return (
    <div
      className={classnames(
        "flex",
        "flex-col",
        "items-center",
        "justify-center",
        "mt-2",
        "p-2"
      )}
    >
      {!lot && <ConstructBuilding lotId={id} />}
      {lot?.building === Building.KITCHEN && (
        <>
          <h2>Kitchen</h2>
          <SvgKitchen width={100} height={100} />
          <p className={classnames("my-4", "text-gray-700")}>
            Wow! It's hot in here. This is were you chefs are making pizza.
          </p>
          <p className={classnames("my-4", "text-gray-700")}>
            There are currently{" "}
            <span className={classnames("font-bold", "text-gray-900")}>
              {stats?.employedChefs} employed chefs
            </span>{" "}
            in your town out of your{" "}
            <span className={classnames("font-bold", "text-gray-900")}>
              {population.chefs} educated chefs
            </span>
            .
          </p>
          <p className={classnames("my-4", "text-gray-700")}>
            If you build more kitchens you can have even more employed chefs!
          </p>
        </>
      )}
      {lot?.building === Building.HOUSE && (
        <>
          <h2>House</h2>
          <SvgHouse height={50} width={50} />
          <p className={classnames("my-4", "text-gray-700")}>
            Up to 10 mice can live in this small house.
          </p>
          <p className={classnames("my-4", "text-gray-700")}>
            If you build more houses your population will grow.
          </p>
        </>
      )}
      {lot?.building === Building.SHOP && (
        <>
          <h2>Shop</h2>
          <SvgShop height={100} width={100} />
          <p className={classnames("my-4", "text-gray-700")}>
            This is were your salesmice work to sell pizzas.
          </p>
          <p className={classnames("my-4", "text-gray-700")}>
            There are currently{" "}
            <span className={classnames("font-bold", "text-gray-900")}>
              {stats?.employedSalesmice} employed salesmice
            </span>{" "}
            in your town out of your{" "}
            <span className={classnames("font-bold", "text-gray-900")}>
              {population.salesmice} educated salesmice
            </span>
            .
          </p>
          <p className={classnames("my-4", "text-gray-700")}>
            If you build more shops you can have even more employed chefs!
          </p>
        </>
      )}
      {lot?.building === Building.SCHOOL && (
          <School />
      )}
    </div>
  );
}

export default TownLot;
