import React, { useCallback, useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import classnames from "classnames";
import { Building } from "../generated/building";
import { Lot, useStore } from "../store";
import ConstructBuilding from "./ConstructBuilding";
import School from "./School";
import { ReactComponent as SvgKitchen } from "../../images/kitchen.svg";
import { ReactComponent as SvgShop } from "../../images/shop.svg";
import { ReactComponent as SvgHouse } from "../../images/house.svg";
import { ReactComponent as SvgMarketingHQ } from "../../images/marketing-hq.svg";
import styles from "../styles";
import { formatDistanceToNow } from "date-fns";
import { formatDurationShort, formatNumber, getTapInfo } from "../utils";
import { Education } from "../generated/education";
import ResearchInstitute from "./buildings/ResearchInstitute";
import { useTimeoutFn } from "react-use";
import { confetti } from "../confetti";
import classes from "./town-lot.module.css";
import classNames from "classnames";
import TapStreak from "./TapStreak";
import ReactDOM from "react-dom";
import { Coin, Pizza } from "../icons";
import { TownCentre } from "./buildings/TownCentre";
import { CountDown } from "./CountDown";
import { useEducationCount } from "./useEducationCount";
import { UpgradeSection } from "./town/UpgradeSection";
import { ResearchDiscovery } from "../generated/research";

const label = classnames("text-xs", "md:text-sm", "mr-1");
const value = classnames("text-sm", "md:text-lg", "ml-1");

const pizzaElement = document.createElement("div");
ReactDOM.render(<Pizza className={classnames("w-12 h-12")} />, pizzaElement);

const coinElement = document.createElement("div");
ReactDOM.render(<Coin className={classnames("w-12 h-12")} />, coinElement);

function getTapBonusFactor(discoveries: ResearchDiscovery[]): number {
  let tapBonusFactor = 1.0;

  if (discoveries.includes(ResearchDiscovery.SLAM)) {
    tapBonusFactor += 0.3;
  }
  if (discoveries.includes(ResearchDiscovery.HIT_IT)) {
    tapBonusFactor += 0.5;
  }
  if (discoveries.includes(ResearchDiscovery.GRAND_SLAM)) {
    tapBonusFactor += 0.75;
  }
  if (discoveries.includes(ResearchDiscovery.GODS_TOUCH)) {
    tapBonusFactor += 1.0;
  }
  return tapBonusFactor;
}

const TapSection: React.VFC<{ lotId: string; lot: Lot }> = ({ lot, lotId }) => {
  const [now, setNow] = useState(new Date());
  const [tapBackoff, setTapBackoff] = useState(false);
  const buttonConfettiRef = useRef<HTMLDivElement>(null);
  const discoveries = useStore((state) => state.gameState.discoveries);

  const tap = useStore((state) => state.tap);

  const { nextTapAt, canTap, taps, tapsRemaining, streak } = getTapInfo(
    lot,
    discoveries,
    now
  );

  const { 2: reset } = useTimeoutFn(
    () => setNow(new Date()),
    Math.max(Math.min(nextTapAt - Date.now(), 10_000), 16)
  );

  useEffect(() => {
    reset();
  }, [nextTapAt]);

  useEffect(() => {
    let timerId = -1;
    if (tapBackoff) {
      timerId = window.setTimeout(() => {
        setTapBackoff(false);
      }, 500);
    }

    return () => window.clearTimeout(timerId);
  }, [tapBackoff]);

  let tapResource: "pizzas" | "coins";
  let tapGains;
  const factor =
    Math.sqrt((lot.level + 1) * (streak + 1)) * getTapBonusFactor(discoveries);
  switch (lot.building) {
    case Building.KITCHEN:
      tapResource = "pizzas";
      tapGains = Math.round((80 * factor) / 5) * 5;
      break;
    case Building.SHOP:
      tapResource = "coins";
      tapGains = Math.round((35 * factor) / 5) * 5;
      break;
    default:
      return null;
  }

  const onClick = useCallback<React.MouseEventHandler>(
    (e) => {
      e.preventDefault();
      tap(lotId);
      setTapBackoff(true);

      if (buttonConfettiRef.current) {
        confetti(buttonConfettiRef.current, {
          elementCount: 3 * Math.ceil(Math.sqrt(lot.level + 1)),
          colors: [],
          startVelocity: 27,
          spread: 35,
          duration: 2000,
          createElement: () => {
            return (tapResource === "pizzas"
              ? pizzaElement.cloneNode(true)
              : coinElement.cloneNode(true)) as HTMLElement;
          },
        });
      }
    },
    [lotId, lot, pizzaElement, coinElement]
  );

  return (
    <section
      className={classnames(
        "m-4",
        "p-4",
        "bg-green-200",
        "flex",
        "items-center",
        "flex-col"
      )}
      data-cy="tap-section"
    >
      <div className="relative select-none">
        <button
          className={classNames(styles.primaryButton, classes.tapButton, {
            [classes.hasTapsRemaining]: tapsRemaining > 0,
          })}
          disabled={!canTap || tapBackoff}
          onClick={onClick}
        >
          <span>
            +{tapGains} {tapResource}
            <br />
            {taps} of 10
          </span>
        </button>
        <div
          ref={buttonConfettiRef}
          className={classnames("pointer-events-none")}
        />
      </div>
      <div className="mt-2 w-64 xs:w-80 h-6">
        <div
          className={classnames("text-center text-sm xs:text-base", {
            hidden: canTap || tapsRemaining > 0,
          })}
        >
          Next tap{" "}
          {formatDistanceToNow(new Date(nextTapAt), {
            includeSeconds: true,
            addSuffix: true,
          })}
        </div>
      </div>
      <div className={classnames("mt-2")}>
        <div className={classnames("text-center")}>Streak:</div>
        <TapStreak value={streak} max={12} />
      </div>
    </section>
  );
};

const RazeSection: React.VFC<{ lotId: string; lot: Lot }> = ({
  lot,
  lotId,
}) => {
  const coins = useStore((state) => state.gameState.resources?.coins ?? 0);
  const constructionQueue = useStore(
    (state) => state.gameState.constructionQueue
  );
  const gameData = useStore((state) => state.gameData);
  const buildingInfo = gameData?.buildings[lot.building];
  const razeBuilding = useStore((state) => state.razeBuilding);
  const cancelRazeBuilding = useStore((state) => state.cancelRazeBuilding);

  if (buildingInfo == null) {
    return null;
  }

  const onClick = () => {
    razeBuilding(lotId);
    window.scroll(0, 0);
  };

  const onCancelClick = () => {
    cancelRazeBuilding(lotId);
    window.scroll(0, 0);
  };

  const levelInfo = buildingInfo.levelInfos[lot.level];
  const constructionTime = Math.floor(levelInfo.constructionTime * 2);
  const cost = Math.floor(levelInfo.cost / 2);

  const canAfford = coins >= cost;

  const constr = constructionQueue.find((x) => x.lotId === lotId);
  if (constr) {
    return constr.razing ? (
      <section
        className={classnames("m-4", "p-4", "bg-red-200")}
        data-cy="raze-section"
      >
        <div>This building is being razed.</div>
        <button
          className={classnames(...styles.button, "bg-red-800")}
          onClick={onCancelClick}
          data-cy="cancel-raze-building-button"
        >
          Cancel
        </button>
      </section>
    ) : null;
  }

  return (
    <section
      className={classnames("m-4", "p-4", "bg-gray-300")}
      data-cy="raze-section"
    >
      <table>
        <tbody>
          <tr>
            <td className={classnames(label, "pr-2")}>Raze Cost:</td>
            <td className={classnames(value, "pr-2")}>
              {formatNumber(cost)} coins
            </td>
          </tr>
          <tr>
            <td className={classnames(label, "pr-2")}>Raze time:</td>
            <td className={classnames(value, "pr-2")}>
              {formatDurationShort(constructionTime)}
            </td>
          </tr>
        </tbody>
      </table>
      <hr className={classnames("border-t-2", "border-red-500", "my-2")} />
      {!canAfford && (
        <div className={classnames("m-2", "text-sm", "text-red-800")}>
          Not enough coins
        </div>
      )}
      <button
        className={classnames(...styles.button, "bg-red-800")}
        disabled={!canAfford}
        onClick={onClick}
        data-cy="raze-building-button"
      >
        Raze {buildingInfo.title}
      </button>
    </section>
  );
};

function TownLot() {
  const { id } = useParams();

  const lot = useStore(
    useCallback(
      (state) => (id !== undefined ? state.gameState.lots[id] : undefined),
      [id]
    )
  );
  const stats = useStore((state) => state.gameStats);
  const educationCount = useEducationCount();
  const gameData = useStore((state) => state.gameData);

  const ongoingConstruction = useStore(
    useCallback(
      (state) => state.gameState.constructionQueue.find((x) => x.lotId === id),
      [id]
    )
  );

  if (id === undefined) {
    return null;
  }

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
      data-cy="town-lot"
    >
      {!lot && !ongoingConstruction && <ConstructBuilding lotId={id} />}
      {!lot && ongoingConstruction && (
        <>
          <h2>Construction site</h2>
          <p className={classnames("mt-4", "text-gray-700")}>
            A{" "}
            {gameData?.buildings[
              ongoingConstruction.building
            ].title.toLowerCase()}{" "}
            is being constructed here.
          </p>
          <p className={classnames("mt-2", "text-gray-700")}>
            It will be ready <CountDown time={ongoingConstruction.completeAt} />
            .
          </p>
        </>
      )}
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
              {educationCount[Education.CHEF] ?? 0} educated chefs
            </span>
            .
          </p>
          <p className={classnames("my-4", "text-gray-700")}>
            If you upgrade or build more kitchens you can have even more
            employed chefs!
          </p>
          <UpgradeSection lot={lot} lotId={id} />
          <RazeSection lot={lot} lotId={id} />
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
          <RazeSection lot={lot} lotId={id} />
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
              {educationCount[Education.SALESMOUSE] ?? 0} educated salesmice
            </span>
            .
          </p>
          <p className={classnames("my-4", "text-gray-700")}>
            If you upgrade or build more shops you can have even more employed
            chefs!
          </p>
          <UpgradeSection lot={lot} lotId={id} />
          <RazeSection lot={lot} lotId={id} />
        </>
      )}
      {lot?.building === Building.MARKETINGHQ && (
        <>
          <h2>Marketing HQ</h2>
          <SvgMarketingHQ height={100} width={100} />
          <p className={classnames("my-4", "text-gray-700")}>
            This is were your marketing personel work.
          </p>
          <UpgradeSection lot={lot} lotId={id} />
          <RazeSection lot={lot} lotId={id} />
        </>
      )}
      {lot?.building === Building.RESEARCH_INSTITUTE && <ResearchInstitute />}
      {lot?.building === Building.TOWN_CENTRE && (
        <TownCentre lot={lot} lotId={id} />
      )}
      {lot?.building === Building.SCHOOL && <School />}
    </div>
  );
}

export default TownLot;
