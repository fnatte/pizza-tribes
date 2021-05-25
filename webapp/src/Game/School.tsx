import { formatDistanceToNow } from "date-fns";
import JSBI from "jsbi";
import React, { useState } from "react";
import { useInterval } from "react-use";
import { classnames } from "tailwindcss-classnames";
import { Education } from "../generated/education";
import { useStore } from "../store";
import styles from "../styles";
import PlaceholderImage from "./PlaceholderImage";
import { ReactComponent as ChefSvg } from "../../images/chef.svg";
import { ReactComponent as SalesmouseSvg } from "../../images/salesmouse.svg";
import { ReactComponent as SecuritySvg } from "../../images/security.svg";
import { ReactComponent as ThiefSvg } from "../../images/thief.svg";
import { ReactComponent as PublicistSvg } from "../../images/publicist.svg";
import { ReactComponent as SvgSchool } from "../../images/school.svg";
import { formatDurationShort } from "../utils";
import { useForm } from "react-hook-form";
import { EducationInfo } from "../generated/game_data";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { RemoveIndex } from "../utils";

const title = classnames("text-lg", "md:text-xl", "mb-2");
const label = classnames("text-xs", "md:text-sm");
const value = classnames("text-sm");
const descriptionStyle = classnames("text-sm", "text-gray-600");

const svgs: Record<number, React.VFC | undefined> = {
  [Education.CHEF]: ChefSvg,
  [Education.SALESMOUSE]: SalesmouseSvg,
  [Education.GUARD]: SecuritySvg,
  [Education.THIEF]: ThiefSvg,
  [Education.PUBLICIST]: PublicistSvg,
};

const descriptions: Record<number, React.VFC | undefined> = {
  [Education.CHEF]: () => (
    <p className={descriptionStyle}>Chefs work in kitchens to make pizzas.</p>
  ),
  [Education.SALESMOUSE]: () => (
    <p className={descriptionStyle}>Salesmice work in shops to sell pizzas.</p>
  ),
  [Education.GUARD]: () => (
    <p className={descriptionStyle}>
      Security guards help protect your town against thieves.
    </p>
  ),
  [Education.THIEF]: () => (
    <p className={descriptionStyle}>
      Thieves can be sent to other towns to steal coins.
    </p>
  ),
  [Education.PUBLICIST]: () => (
    <p className={descriptionStyle}>
      Publicists work in marketing HQ to increase your popularity and demand.
    </p>
  ),
};

const numberFormat = new Intl.NumberFormat();

const schema = yup.object().shape({
  amount: yup
    .number()
    .typeError("amount must be a number")
    .integer()
    .positive()
    .required(),
});

type FormFields = RemoveIndex<yup.Asserts<typeof schema>>;

const SchoolEducation: React.VFC<{
  education: Education;
  educationInfo: EducationInfo;
}> = ({ education, educationInfo }) => {
  const coins = useStore((state) => state.gameState.resources.coins);
  const train = useStore((state) => state.train);
  const uneducated = useStore((state) => state.gameState.population.uneducated);

  const maxByCost = educationInfo.cost
    ? Math.floor(coins / educationInfo.cost)
    : 1_000;
  const max = Math.min(uneducated, maxByCost);

  const schema2 = schema.shape({
    amount: schema.fields.amount.max(
      max,
      uneducated < maxByCost
        ? `Cannot train more than available uneducated mice (${uneducated}).`
        : `Not enough coins. Can only afford ${maxByCost}).`
    ),
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormFields>({ resolver: yupResolver(schema2) });

  const onSubmit = ({ amount }: FormFields) => {
    train(education, amount);
    window.scroll(0, 0);
  };

  const disabled = max <= 0;

  const SvgImage = svgs[education];
  const Description = descriptions[education];

  return (
    <div className={classnames("flex", "mb-8")} key={educationInfo.title}>
      <div style={{ width: 110 }}>
        {SvgImage ? <SvgImage /> : <PlaceholderImage />}
      </div>
      <div className={classnames("ml-6", "lg:ml-8")}>
        <div className={title}>{educationInfo.title}</div>
        {Description && <Description />}
        <table>
          <tbody>
            <tr>
              <td className={classnames("px-2")}>
                <span className={label}>Train time:</span>
              </td>
              <td className={classnames("px-2")}>
                <span className={value}>
                  {formatDurationShort(educationInfo.trainTime)}
                </span>
              </td>
            </tr>
            {educationInfo.cost > 0 && (
              <tr>
                <td className={classnames("px-2")}>
                  <span className={label}>Cost:</span>
                </td>
                <td className={classnames("px-2")}>
                  <span className={value}>
                    {numberFormat.format(educationInfo.cost)} coins
                  </span>
                </td>
              </tr>
            )}
          </tbody>
        </table>
        <form
          className={classnames("my-2")}
          onSubmit={handleSubmit(onSubmit)}
          noValidate
        >
          <input
            className={classnames("w-20", "mr-2", "text-right")}
            type="number"
            min={0}
            max={max}
            defaultValue={1}
            required
            disabled={disabled}
            {...register("amount")}
          />
          <button
            type="submit"
            className={classnames(styles.primaryButton)}
            disabled={disabled}
          >
            Train
          </button>
          {errors.amount && (
            <div className="pb-2 text-red-900">{errors.amount.message}</div>
          )}
        </form>
      </div>
    </div>
  );
};

function School() {
  const educations = useStore((state) => state.gameData?.educations) || [];
  const trainingQueue = useStore((state) => state.gameState.trainingQueue);
  const uneducated = useStore((state) => state.gameState.population.uneducated);

  const [_, setNow] = useState(Date.now());
  useInterval(() => {
    setNow(Date.now());
  }, 1000);

  return (
    <div className={classnames("max-w-full", "px-2")}>
      <h2>School</h2>
      <SvgSchool height={100} width={100} />
      <p className={classnames("my-4", "text-gray-700")}>
        Educate your mice so that they can start contributing to your town.
      </p>
      <p className={classnames("my-4", "text-gray-700")}>
        There are{" "}
        <span className={classnames("font-bold", "text-gray-900")}>
          {uneducated} uneducated
        </span>{" "}
        mice in your town.
      </p>

      {trainingQueue.length > 0 && (
        <>
          <h3>In Training</h3>
          <table>
            <tbody>
              {trainingQueue.map((training) => (
                <tr
                  key={
                    training.completeAt.toString() +
                    training.education +
                    training.amount
                  }
                >
                  <td className={classnames("p-2")}>{training.amount}</td>
                  <td className={classnames("p-2")}>
                    {educations[training.education].title}
                  </td>
                  <td className={classnames("p-2")}>
                    {formatDistanceToNow(
                      Number(
                        JSBI.divide(
                          JSBI.BigInt(training.completeAt),
                          JSBI.BigInt(1e6)
                        )
                      ),
                      {
                        includeSeconds: true,
                        addSuffix: true,
                      }
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}

      <h3>Train</h3>
      {Object.keys(educations)
        .map(Number)
        .map((eduKey) => {
          return (
            <SchoolEducation
              key={eduKey}
              education={eduKey}
              educationInfo={educations[eduKey]}
            />
          );
        })}
    </div>
  );
}

export default School;
