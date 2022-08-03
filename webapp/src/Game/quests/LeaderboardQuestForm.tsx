import React, { useState } from "react";
import * as yup from "yup";
import classnames from "classnames";
import { RemoveIndex } from "../../utils";
import { useStore } from "../../store";
import { yupResolver } from "@hookform/resolvers/yup";
import { useForm } from "react-hook-form";
import { LEADERBOARD_QUEST_ID } from "./ids";
import { button } from "../../styles";
import { useAsync } from "react-use";
import { apiFetch } from "../../api";

export function LeaderboardQuestForm() {
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
  const [answerResult, setAnswerResult] = useState<boolean>();

  const rank = useAsync(async () => {
    const response = await apiFetch("/leaderboard/me/rank");
      if (
        !response.ok ||
        response.headers.get("Content-Type") !== "application/json"
      ) {
        throw new Error("Failed to get leaderboard rank");
      }
      const body = await response.json();
    return body as number;
  }, [])

  const onSubmit = async (data: FormFields) => {
    if (rank.error || rank.loading)  {
      return;
    }

    const correctAnswer = rank.value;
    const answer = parseInt(data.answer.replace("th", ""));
    if (correctAnswer === answer) {
      completeQuest(LEADERBOARD_QUEST_ID);
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
