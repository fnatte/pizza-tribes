import React, {
  ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import classnames from "classnames";
import { useStore } from "../store";
import { Quest } from "../generated/quest";
import ReactMarkdown from "react-markdown";
import { Coin, Pizza } from "../icons";
import { QuestState } from "../generated/gamestate";
import { primaryButton } from "../styles";

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

function getQuestStateImportance(q: QuestState): number {
  if (q.completed && !q.claimedReward) {
    return 2;
  }

  if (!q.completed) {
    return 1;
  }

  return 0;
}

export default function QuestsView() {
  const questStates = useStore((state) => state.gameState.quests);
  const gameData = useStore((state) => state.gameData);
  const claimQuestReward = useStore((state) => state.claimQuestReward);
  const openQuest = useStore(state=>state.openQuest);

  const quests = useMemo(() => {
    const map = new Map<string, Quest>();
    gameData?.quests.forEach((quest) => map.set(quest.id, quest));
    return map;
  }, [gameData]);

  const sortedQuestStates = Object.entries(questStates).reverse().sort(
    ([_a, a], [_b, b]) =>
      getQuestStateImportance(b) - getQuestStateImportance(a)
  );

  const [expanded, setExpanded] = useState<string>();

  const [didSetInitialExpanded, setDidSetInitialExpanded] = useState(false);

  useEffect(() => {
    if (!expanded && !didSetInitialExpanded && sortedQuestStates.length > 0) {
      setDidSetInitialExpanded(true);
      setExpanded(sortedQuestStates[0][0]);
    }
  }, [questStates, expanded, didSetInitialExpanded]);

  console.log(questStates)

  useEffect(() => {
    if (expanded && !questStates[expanded].opened) {
      openQuest(expanded);
    }
  }, [expanded, questStates])

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
              <li key={id} className="my-6">
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
                      "font-bold": !questState.opened,
                    })}
                  >
                    <span className="underline">{quest.title}</span>
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
                    <div className="my-4">
                      <h4>Reward</h4>
                      <div className="flex gap-4 items-center flex-wrap">
                        {(quest.reward?.coins ?? 0) > 0 && (
                          <div className="flex gap-1 items-center">
                            <Coin className={"h-[3em] w-[3em]"} />{" "}
                            <span className="text-xl">
                              {quest.reward?.coins}
                            </span>
                          </div>
                        )}
                        {(quest.reward?.pizzas ?? 0) > 0 && (
                          <div className="flex gap-1 items-center">
                            <Pizza className={"h-[3em] w-[3em]"} />{" "}
                            <span className="text-xl">
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
