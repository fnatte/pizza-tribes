import React, { useState } from "react";
import * as yup from "yup";
import classnames from "classnames";
import { RemoveIndex } from "../../utils";
import { useStore } from "../../store";
import { yupResolver } from "@hookform/resolvers/yup";
import { useForm } from "react-hook-form";
import { STATS_QUEST_ID } from "./ids";
import { button } from "../../styles";

export function StatsQuestForm() {
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
