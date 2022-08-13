import React, { ReactNode, useEffect, useMemo, useState } from "react";
import classnames from "classnames";
import ReactMarkdown from "react-markdown";
import { useStore } from "../../store";
import { Quest } from "../../generated/quest";
import { Coin, Pizza } from "../../icons";
import { primaryButton } from "../../styles";
import { LEADERBOARD_QUEST_ID, STATS_QUEST_ID } from "./ids";
import { StatsQuestForm } from "./StatsQuestForm";
import { LeaderboardQuestForm } from "./LeaderboardQuestForm";

function Container({ children }: { children?: ReactNode | undefined }) {
  return (
    <div
      className={classnames(
        "container",
        "max-w-screen-sm",
        "mx-auto",
        "mt-2",
        "p-2"
      )}
    >
      {children}
    </div>
  );
}

export default function QuestsView() {
  const questStates = useStore((state) => state.gameState.quests);
  const gameData = useStore((state) => state.gameData);
  const claimQuestReward = useStore((state) => state.claimQuestReward);
  const openQuest = useStore((state) => state.openQuest);

  const quests = useMemo(() => {
    const map = new Map<string, Quest>();
    gameData?.quests.forEach((quest) => map.set(quest.id, quest));
    return map;
  }, [gameData]);

  const sortedQuestStates = Object.entries(questStates)
    .reverse()
    .sort(([a], [b]) => {
      return (quests.get(b)?.order ?? 0) - (quests.get(a)?.order ?? 0);
    });

  const [expanded, setExpanded] = useState<string>();

  const [didSetInitialExpanded, setDidSetInitialExpanded] = useState(false);

  useEffect(() => {
    if (!expanded && !didSetInitialExpanded && sortedQuestStates.length > 0) {
      setDidSetInitialExpanded(true);
      setExpanded(sortedQuestStates[0]![0]);
    }
  }, [questStates, expanded, didSetInitialExpanded]);

  useEffect(() => {
    if (expanded && questStates[expanded]?.opened === false) {
      openQuest(expanded);
    }
  }, [expanded, questStates]);

  if (!gameData) {
    return null;
  }

  return (
    <Container>
      <h2>Quests</h2>

      <section className="my-6">
        <ul>
          {sortedQuestStates.map(([id, questState]) => {
            const quest = quests.get(id);
            if (!quest) {
              return null;
            }
            const canClaimReward =
              questState.completed && !questState.claimedReward;
            return (
              <li key={id} className="my-6" data-cy="quest-item">
                <div className="flex items-center">
                  <input
                    type="checkbox"
                    checked={questState.completed && questState.claimedReward}
                    readOnly
                    className="mr-2 focus:ring-offset-0 focus:ring-0"
                  />
                  <button
                    onClick={() =>
                      setExpanded(expanded !== id ? id : undefined)
                    }
                    className={classnames({
                      "font-bold": !questState.opened || canClaimReward,
                    })}
                    aria-expanded={expanded === id}
                    data-cy="quest-item-expand-toggle"
                  >
                    <span className="underline" data-cy="quest-item-title">
                      {quest.title}
                    </span>
                    {expanded !== id && canClaimReward ? (
                      <span className="ml-1">(Claim reward)</span>
                    ) : null}
                  </button>
                </div>
                {expanded === id && (
                  <div className="mt-2 p-6 bg-green-50">
                    <ReactMarkdown
                      className={classnames("prose", "text-gray-700")}
                      disallowedElements={["img"]}
                    >
                      {quest.description}
                    </ReactMarkdown>

                    {id === STATS_QUEST_ID && !questState.completed && (
                      <div className="my-4">
                        <StatsQuestForm />
                      </div>
                    )}

                    {id === LEADERBOARD_QUEST_ID && !questState.completed && (
                      <div className="my-4">
                        <LeaderboardQuestForm />
                      </div>
                    )}

                    <div className="my-4">
                      <h4>Reward</h4>
                      <div className="flex gap-4 items-center flex-wrap">
                        {(quest.reward?.coins ?? 0) > 0 && (
                          <div className="flex gap-1 items-center">
                            <Coin className={"h-[3em] w-[3em]"} />{" "}
                            <span
                              className="text-xl"
                              data-cy="quest-item-reward-coins"
                            >
                              {quest.reward?.coins}
                            </span>
                          </div>
                        )}
                        {(quest.reward?.pizzas ?? 0) > 0 && (
                          <div className="flex gap-1 items-center">
                            <Pizza className={"h-[3em] w-[3em]"} />{" "}
                            <span
                              className="text-xl"
                              data-cy="quest-item-reward-pizzas"
                            >
                              {quest.reward?.pizzas}
                            </span>
                          </div>
                        )}
                        {questState.claimedReward && (
                          <div className="text-gray-700">
                            (Reward has been claimed)
                          </div>
                        )}
                      </div>
                      {canClaimReward && (
                        <div className="mt-4">
                          <button
                            className={primaryButton}
                            onClick={() => claimQuestReward(id)}
                            data-cy="quest-item-claim-reward-button"
                          >
                            Claim Reward
                          </button>
                        </div>
                      )}
                    </div>
                  </div>
                )}
              </li>
            );
          })}
        </ul>
      </section>
    </Container>
  );
}
