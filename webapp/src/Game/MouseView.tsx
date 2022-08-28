import React, { ReactNode, useCallback } from "react";
import classnames from "classnames";
import * as yup from "yup";
import { useStore } from "../store";
import { Link, useNavigate, useParams } from "react-router-dom";
import { button } from "../styles";
import { useForm } from "react-hook-form";
import { RemoveIndex } from "../utils";
import { yupResolver } from "@hookform/resolvers/yup";
import { Education } from "../generated/education";
import { ReactComponent as ChefSvg } from "../../images/chef.svg";
import { ReactComponent as SalesmouseSvg } from "../../images/salesmouse.svg";
import { ReactComponent as SecuritySvg } from "../../images/security.svg";
import { ReactComponent as ThiefSvg } from "../../images/thief.svg";
import { ReactComponent as PublicistSvg } from "../../images/publicist.svg";
import { ReactComponent as UneducatedSvg } from "../../images/uneducated.svg";
import { MouseImage } from "./components/MouseImage";

const svgs: Record<number, React.VFC | undefined> = {
  [Education.CHEF]: ChefSvg,
  [Education.SALESMOUSE]: SalesmouseSvg,
  [Education.GUARD]: SecuritySvg,
  [Education.THIEF]: ThiefSvg,
  [Education.PUBLICIST]: PublicistSvg,
};

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
  name: yup
    .string()
    .trim()
    .min(2)
    .max(30)
    .matches(
      /^[a-zA-Z\s]+$/,
      "name must consist of letters (a-z) and up to two spaces"
    )
    .matches(
      /^[a-zA-Z]+(\s[a-zA-Z]+)?(\s[a-zA-Z]+)?$/,
      "name can have up to two spaces and there must be letter between"
    )
    .required(),
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
      {errors.name && <div className="p-2">{errors.name.message}</div>}
    </form>
  );
}

export default function MouseView() {
  const { id } = useParams();
  const navigate = useNavigate();

  const mouse = useStore(
    useCallback(
      (state) => (id !== undefined ? state.gameState.mice[id] : undefined),
      [id]
    )
  );

  const ambassador = useStore((state) => state.gameState.ambassadorMouseId);
  const isAmbassador = ambassador !== "" && ambassador === id;

  const reschool = useStore((store) => store.reschool);
  const setAmbassadorMouse = useStore((store) => store.setAmbassadorMouse);

  const gameData = useStore((state) => state.gameData);

  const onClickSetAmbassador = useCallback(() => id && setAmbassadorMouse(id), [
    id,
  ]);
  const onClickReschool = useCallback(() => id && reschool(id), [id]);
  const onClickChangeAppearance = () => navigate(`/mouse/${id}/appearance`);

  if (!gameData) {
    return null;
  }

  if (!mouse || !id) {
    return <Container>Could not find mouse.</Container>;
  }

  const { name, isEducated, isBeingEducated } = mouse;

  const education = gameData.educations[mouse.education];

  const FallbackSvg = mouse.isEducated
    ? svgs[mouse.education] ?? UneducatedSvg
    : UneducatedSvg;

  return (
    <Container>
      <h2 className="mb-4">Mouse: {name}</h2>
      <div className="flex gap-2">
        {
          <Link
            to={`/mouse/${id}/appearance`}
            title="Change Appearance"
            className="grow max-w-[250px]"
          >
            {mouse.appearance ? (
              <MouseImage
                appearance={mouse.appearance}
                shiftRight
                height={400}
                className="h-[400px] w-full"
              />
            ) : (
              <FallbackSvg
                height={400}
                className="translate-x-4 xs:translate-x-16 h-[400px] w-full"
              />
            )}
          </Link>
        }
        <div className="flex flex-col justify-center">
          <section className="my-6">{name} is a happy house.</section>
        </div>
      </div>

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
        <section className="my-6">
          <button
            className={classnames(...button, "bg-green-600")}
            onClick={onClickChangeAppearance}
          >
            Change Appearance
          </button>
        </section>
        <section className="flex gap-2 my-6 items-center">
          <button
            className={classnames(...button, "bg-green-600")}
            disabled={isAmbassador}
            onClick={onClickSetAmbassador}
            data-cy="make-ambassador-button"
          >
            Make Ambassador
          </button>
          {isAmbassador && `(${mouse.name} is already ambassador)`}
        </section>
      </section>
    </Container>
  );
}
