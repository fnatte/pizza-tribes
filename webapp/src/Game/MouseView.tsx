import React, { ReactNode, useCallback } from "react";
import classnames from "classnames";
import * as yup from "yup";
import { useStore } from "../store";
import { useParams } from "react-router-dom";
import { button } from "../styles";
import { useForm } from "react-hook-form";
import { RemoveIndex } from "../utils";
import { yupResolver } from "@hookform/resolvers/yup";

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

const schema = yup.object().shape({
  name: yup.string().trim().min(2).max(30).required(),
});

type FormFields = RemoveIndex<yup.Asserts<typeof schema>>;

function RenameForm({ mouseId }: { mouseId: string }) {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormFields>({ resolver: yupResolver(schema) });

  const renameMouse = useStore((state) => state.renameMouse);

  const onSubmit = async (data: FormFields) => {
    renameMouse(mouseId, data.name);
    reset();
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {errors.name && <div className="p-2">{errors.name.message}</div>}
      <label className="mx-2">
        New Name:
        <input className="mx-2" type="text" {...register("name")} />
      </label>
      <button
        type="submit"
        disabled={isSubmitting}
        className={classnames(...button, "bg-green-600")}
      >
        Rename
      </button>
    </form>
  );
}

export default function MouseView() {
  const { id } = useParams();

  const mouse = useStore(
    useCallback(
      (state) => (id !== undefined ? state.gameState.mice[id] : undefined),
      [id]
    )
  );

  const reschool = useStore((store) => store.reschool);

  const gameData = useStore((state) => state.gameData);

  const onClickReschool = useCallback(() => id && reschool(id), [id]);

  if (!gameData) {
    return null;
  }

  if (!mouse || !id) {
    return <Container>Could not find mouse.</Container>;
  }

  const { name, isEducated, isBeingEducated } = mouse;

  const education = gameData.educations[mouse.education];

  return (
    <Container>
      <h2>Mouse: {name}</h2>

      <section className="my-6">{name} is a happy house.</section>

      <section className="my-6">
        <h3>Education</h3>
        <p>
          {name}{" "}
          {!isEducated ? (
            "is uneducated"
          ) : isBeingEducated ? (
            "is in school"
          ) : (
            <>
              is a{" "}
              <span className="font-bold">{education.title.toLowerCase()}</span>
            </>
          )}
          .
        </p>
        {isEducated && (
          <div
            className={classnames(
              "inline-block",
              "my-4",
              "mx-auto",
              "p-4",
              "bg-gray-200"
            )}
          >
            <p>Do you want {name} to have another education?</p>
            <hr
              className={classnames("border-t-2", "border-red-500", "my-2")}
            />
            <button
              className={classnames(...button, "bg-red-800")}
              onClick={onClickReschool}
            >
              Re-school
            </button>
          </div>
        )}
      </section>

      <section className="my-6">
        <h3>Manage</h3>
        <RenameForm mouseId={id} />
      </section>
    </Container>
  );
}
