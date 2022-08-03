import React, { ReactNode, useEffect, useMemo, useState } from "react";
import * as yup from "yup";
import classnames from "classnames";
import { useStore } from "../store";
import { Quest } from "../generated/quest";
import ReactMarkdown from "react-markdown";
import { Coin, Pizza } from "../icons";
import { QuestState } from "../generated/gamestate";
import { button, primaryButton } from "../styles";
import { useForm } from "react-hook-form";
import { RemoveIndex } from "../utils";
import { yupResolver } from "@hookform/resolvers/yup";

const STATS_QUEST_ID = "9";

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
  const openQuest = useStore((state) => state.openQuest);

  const quests = useMemo(() => {
    const map = new Map<string, Quest>();
    gameData?.quests.forEach((quest) => map.set(quest.id, quest));
    return map;
  }, [gameData]);

  const sortedQuestStates = Object.entries(questStates)
    .reverse()
    .sort(
      ([_a, a], [_b, b]) =>
        getQuestStateImportance(b) - getQuestStateImportance(a)
    );

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
                      "font-bold": !questState.opened,
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
                    {id === STATS_QUEST_ID && (
                      <div className="my-4">
                        <StatsQuestForm />
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

function StatsQuestForm() {
  const schema = yup.object().shape({
    answer: yup.string().trim().required(),
  });

  type FormFields = RemoveIndex<yup.Asserts<typeof schema>>;

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<{
    answer: string;
  }>({ resolver: yupResolver(schema) });

  const completeQuest = useStore((state) => state.completeQuest);
  const stats = useStore((state) => state.gameStats);
  const [answerResult, setAnswerResult] = useState<boolean>();

  const onSubmit = async (data: FormFields) => {
    const correctAnswer = stats?.pizzasProducedPerSecond;
    const answer = parseFloat(data.answer.replace("/s", "").replace(",", "."));
    if (
      correctAnswer !== undefined &&
      Math.abs(answer - correctAnswer) < 0.05
    ) {
      completeQuest(STATS_QUEST_ID);
      setAnswerResult(true);
    } else {
      setAnswerResult(false);
    }
    reset();
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} data-cy="stats-quest-form">
      {answerResult !== true && (
        <>
          {errors.answer && <div className="p-2">{errors.answer.message}</div>}
          <label className="mx-2">
            Answer:
            <input className="mx-2" type="text" {...register("answer")} />
          </label>
          <button
            type="submit"
            disabled={isSubmitting || answerResult}
            className={classnames(...button, "bg-green-600")}
          >
            Send Answer
          </button>
        </>
      )}
      {answerResult !== undefined && (
        <div className="p-2">{answerResult ? "Correct!" : "Wrong!"}</div>
      )}
    </form>
  );
}
