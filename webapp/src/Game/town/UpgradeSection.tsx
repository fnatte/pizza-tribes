import React from "react";
import classnames from "classnames";
import { Education } from "../../generated/education";
import { useStore, Lot } from "../../store";
import { CountDown } from "../CountDown";
import { formatDurationShort, formatNumber } from "../../utils";
import styles from "../../styles";
import ReactMarkdown from "react-markdown";

const label = classnames("text-xs", "md:text-sm", "mr-1");
const value = classnames("text-sm", "md:text-lg", "ml-1");

export const UpgradeSection: React.VFC<{ lotId: string; lot: Lot }> = ({
  lot,
  lotId,
}) => {
  const coins = useStore((state) => state.gameState.resources?.coins ?? 0);
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
      <section
        className={classnames("m-4", "p-4", "bg-green-200")}
        data-cy="upgrade-section"
      >
        <span>Already at max level</span>
      </section>
    );
  }

  const constr = constructionQueue.find((x) => x.lotId === lotId);
  if (constr) {
    return !constr.razing ? (
      <section
        className={classnames("m-4", "p-4", "bg-green-200")}
        data-cy="upgrade-section"
      >
        <p>This building is being upgraded.</p>
        <p>
          It will be ready <CountDown time={constr.completeAt} />.
        </p>
      </section>
    ) : null;
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
    <section
      className={classnames("m-4", "p-4", "bg-green-200")}
      data-cy="upgrade-section"
    >
      {nextLevelInfo.description && (
        <ReactMarkdown
          className={"mb-2 prose text-black text-sm md:text-lg"}
          disallowedElements={["img"]}
        >
          {nextLevelInfo.description}
        </ReactMarkdown>
      )}
      <table>
        <tbody>
          <tr>
            <td className={classnames(label, "pr-2")}>Cost:</td>
            <td className={classnames(value, "pr-2")}>
              {formatNumber(cost)} coins
            </td>
          </tr>
          <tr>
            <td className={classnames(label, "pr-2")}>Build time:</td>
            <td className={classnames(value, "pr-2")}>
              {formatDurationShort(constructionTime)}
            </td>
          </tr>
          {increasedPopulation > 0 && (
            <tr>
              <td className={classnames(label, "pr-2")}>Population:</td>
              <td className={classnames(value, "pr-2")}>
                +{formatNumber(increasedPopulation)}
              </td>
            </tr>
          )}
          {increasesWorkforce > 0 && (
            <tr>
              <td className={classnames(label, "pr-2")}>
                {employsChefs && "Chef positions:"}
                {employsSalesmice && "Salesmouse positions:"}
              </td>
              <td className={classnames(value, "pr-2")}>
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
      <button
        className={styles.primaryButton}
        disabled={!canAfford}
        onClick={onClick}
        data-cy="upgrade-building-button"
      >
        Upgrade to level {lot.level + 2}
      </button>
    </section>
  );
};
