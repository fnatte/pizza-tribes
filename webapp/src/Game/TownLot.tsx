import React, { useCallback } from "react";
import { useParams } from "react-router-dom";
import { classnames, TArg } from "tailwindcss-classnames";
import { Building } from "../generated/building";
import { Lot, useStore } from "../store";
import ConstructBuilding from "./ConstructBuilding";
import School from "./School";
import { ReactComponent as SvgKitchen } from "../../images/kitchen.svg";
import { ReactComponent as SvgShop } from "../../images/shop.svg";
import { ReactComponent as SvgHouse } from "../../images/house.svg";
import JSBI from "jsbi";
import styles from "../styles";
import { formatDistanceToNow } from "date-fns";
import { countPopulation, formatDurationShort, formatNumber } from "../utils";
import { Education } from "../generated/education";

const title = classnames("text-lg", "md:text-xl", "mb-2");
const label = classnames("text-xs", "md:text-sm", "mr-1");
const value = classnames("text-sm", "md:text-lg", "ml-1");

const TapSection: React.VFC<{ lotId: string; lot: Lot }> = ({ lot, lotId }) => {
  const population = useStore((state) => state.gameState.population);

  // convert lot.tappedAt from ns ms
  const tappedAt = JSBI.toNumber(
    JSBI.divide(JSBI.BigInt(lot.tappedAt), JSBI.BigInt(1e6))
  );
  const nextTapAt = tappedAt + 60_000 * 60;
  const canTap = nextTapAt < Date.now();

  const tap = useStore((state) => state.tap);

  let tapResource;
  let tapGains;
  switch (lot.building) {
    case Building.KITCHEN:
      tapResource = "pizzas";
      tapGains = 80 * countPopulation(population);
      break;
    case Building.SHOP:
      tapResource = "coins";
      tapGains = 35 * countPopulation(population);
      break;
    default:
      return null;
  }

  const onClick = () => {
    tap(lotId);
  };

  return (
    <section className={classnames("m-4", "p-4", "bg-green-200")}>
      <button className={styles.button} disabled={!canTap} onClick={onClick}>
        Tap
      </button>{" "}
      <span>
        (+{tapGains} {tapResource})
      </span>
      {!canTap && (
        <div>
          Next tap in{" "}
          {formatDistanceToNow(new Date(nextTapAt), {
            includeSeconds: true,
            addSuffix: true,
          })}
        </div>
      )}
    </section>
  );
};

const UpgradeSection: React.VFC<{ lotId: string; lot: Lot }> = ({
  lot,
  lotId,
}) => {
  const coins = useStore((state) => state.gameState.resources.coins);
  const constructionQueue = useStore(
    (state) => state.gameState.constructionQueue
  );
  const gameData = useStore((state) => state.gameData);
  const buildingInfo = gameData?.buildings[lot.building];
  const upgradeBuilding = useStore((state) => state.upgradeBuilding);

  if (buildingInfo == null) {
    return null;
  }

  if (lot.level + 1 >= buildingInfo.levelInfos.length) {
    return (
      <section className={classnames("m-4", "p-4", "bg-green-200")}>
        <span>Already at max level</span>
      </section>
    );
  }

  if (constructionQueue.some((x) => x.lotId === lotId)) {
    return (
      <section className={classnames("m-4", "p-4", "bg-green-200")}>
        <span>This building is being upgraded.</span>
      </section>
    );
  }

  const onClick = () => {
    upgradeBuilding(lotId);
  };

  const currentLevelInfo = buildingInfo.levelInfos[lot.level];
  const nextLevelInfo = buildingInfo.levelInfos[lot.level + 1];

  const increasesWorkforce =
    (nextLevelInfo.employer?.maxWorkforce ?? 0) -
    (currentLevelInfo.employer?.maxWorkforce ?? 0);
  const increasedPopulation =
    (nextLevelInfo.residence?.beds ?? 0) -
    (currentLevelInfo.residence?.beds ?? 0);

  const employsSalesmice =
    gameData?.educations[Education.SALESMOUSE].employer === lot.building;
  const employsChefs =
    gameData?.educations[Education.CHEF].employer === lot.building;

  const { cost, constructionTime } = buildingInfo.levelInfos[lot.level + 1];

  const canAfford = coins >= cost;

  return (
    <section className={classnames("m-4", "p-4", "bg-green-200")}>
      <table>
        <tbody>
          <tr>
            <td className={classnames(label as TArg, "pr-2")}>Cost:</td>
            <td className={classnames(value as TArg, "pr-2")}>
              {formatNumber(cost)} coins
            </td>
          </tr>
          <tr>
            <td className={classnames(label as TArg, "pr-2")}>Build time:</td>
            <td className={classnames(value as TArg, "pr-2")}>
              {formatDurationShort(constructionTime)}
            </td>
          </tr>
          {increasedPopulation > 0 && (
            <tr>
              <td className={classnames(label as TArg, "pr-2")}>Population:</td>
              <td className={classnames(value as TArg, "pr-2")}>
                +{formatNumber(increasedPopulation)}
              </td>
            </tr>
          )}
          {increasesWorkforce > 0 && (
            <tr>
              <td className={classnames(label as TArg, "pr-2")}>
                {employsChefs && "Chef positions:"}
                {employsSalesmice && "Salesmouse positions:"}
              </td>
              <td className={classnames(value as TArg, "pr-2")}>
                +{formatNumber(increasesWorkforce)}
              </td>
            </tr>
          )}
        </tbody>
      </table>
      <hr className={classnames("border-t-2", "border-green-300", "my-2")} />
      {!canAfford && (
        <div className={classnames("m-2", "text-sm", "text-red-800")}>
          Not enough coins
        </div>
      )}
      <button className={styles.button} disabled={!canAfford} onClick={onClick}>
        Upgrade to level {lot.level + 2}
      </button>
    </section>
  );
};

function TownLot() {
  const { id } = useParams();

  const lot = useStore(useCallback((state) => state.gameState.lots[id], [id]));
  const stats = useStore((state) => state.gameStats);
  const population = useStore((state) => state.gameState.population);
  const gameData = useStore((state) => state.gameData);

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
          <h2>Kitchen (level {lot.level + 1})</h2>
          <SvgKitchen width={100} height={100} />
          <p className={classnames("my-4", "text-gray-700")}>
            Wow! It's hot in here. This is were you chefs are making pizza.
          </p>
          <TapSection lot={lot} lotId={id} />
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
            If you upgrade or build more kitchens you can have even more
            employed chefs!
          </p>
          <UpgradeSection lot={lot} lotId={id} />
        </>
      )}
      {lot?.building === Building.HOUSE && (
        <>
          <h2>House (level {lot.level + 1})</h2>
          <SvgHouse height={50} width={50} />
          <p className={classnames("my-4", "text-gray-700")}>
            Up to{" "}
            {gameData?.buildings[Building.HOUSE].levelInfos[lot.level].residence
              ?.beds ?? 0}{" "}
            mice can live in this small house.
          </p>
          <p className={classnames("my-4", "text-gray-700")}>
            If you upgrade or build more houses your population will grow.
          </p>
          <UpgradeSection lot={lot} lotId={id} />
        </>
      )}
      {lot?.building === Building.SHOP && (
        <>
          <h2>Shop (level {lot.level + 1})</h2>
          <SvgShop height={100} width={100} />
          <p className={classnames("my-4", "text-gray-700")}>
            This is were your salesmice work to sell pizzas.
          </p>
          <TapSection lot={lot} lotId={id} />
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
            If you upgrade or build more shops you can have even more employed
            chefs!
          </p>
          <UpgradeSection lot={lot} lotId={id} />
        </>
      )}
      {lot?.building === Building.SCHOOL && <School />}
    </div>
  );
}

export default TownLot;
