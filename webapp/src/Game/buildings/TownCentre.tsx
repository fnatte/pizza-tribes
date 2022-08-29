import React, { useCallback, useEffect } from "react";
import classnames from "classnames";
import { ReactComponent as SvgTownCentre } from "../../../images/town-centre.svg";
import { Lot, useStore } from "../../store";
import PopulationTable from "../PopulationTable";
import { Link } from "react-router-dom";
import { smallPrimaryButton } from "../../styles";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { RemoveIndex } from "../../utils";
import { UpgradeSection } from "../town/UpgradeSection";

function MouseTable({ className }: { className?: string }) {
  const mice = useStore((state) => state.gameState.mice);
  const ambassador = useStore((state) => state.gameState.ambassadorMouseId);
  const gameData = useStore((state) => state.gameData);
  const educations = gameData?.educations ?? {};

  return (
    <table className={className}>
      {
        <thead>
          <tr>
            <th colSpan={2} className="p-2 text-left text-sm text-gray-700">
              Mice
            </th>
          </tr>
        </thead>
      }
      <tbody>
        {Object.values(mice).length === 0 && (
          <tr>
            <td colSpan={2} className="p-2 w-80 text-left text-gray-700">
              When your town has some mice, they will show up here.
            </td>
          </tr>
        )}
        {Object.entries(mice).map(([id, mouse]) => {
          return (
            <tr key={id}>
              <td
                className={classnames("p-2", {
                  "font-bold": ambassador !== "" && ambassador === id,
                })}
              >
                {mouse.name}
              </td>
              <td className="p-2">
                {!mouse.isEducated
                  ? "Uneducated"
                  : mouse.isBeingEducated
                  ? "In school"
                  : educations[mouse.education].title}
              </td>
              <td className="p-2">
                <Link
                  to={`/mouse/${id}`}
                  onClick={() => {
                    window.scrollTo(0, 0);
                  }}
                >
                  <button className={classnames(smallPrimaryButton)}>
                    Visit
                  </button>
                </Link>
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}

function PopulationSection() {
  return (
    <section className="my-6">
      <h3>Population</h3>

      <div className="prose my-6 text-gray-700">
        <p>Your tribesmice have good morale and is working efficiently.</p>
      </div>

      <div className="flex flex-wrap items-start gap-24">
        <PopulationTable showHeader showZeroes showTotalCount />
        <MouseTable />
      </div>
    </section>
  );
}

const schema = yup.object().shape({
  pizzaPrice: yup
    .number()
    .label("Price")
    .typeError("Price must be a number")
    .integer()
    .positive()
    .min(1)
    .max(15)
    .required(),
});

type FormFields = RemoveIndex<yup.Asserts<typeof schema>>;

function EconomySection() {
  const setPizzaPrice = useStore((state) => state.setPizzaPrice);
  const pizzaPrice = useStore((state) => state.gameState.pizzaPrice);

  const {
    register,
    setValue,
    getValues,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<FormFields>({
    resolver: yupResolver(schema),
    defaultValues: {
      pizzaPrice,
    },
  });

  const increasePrice = () => {
    let value = getValues("pizzaPrice");
    if (typeof value === "string") {
      value = parseInt(value);
    }
    setValue("pizzaPrice", Math.max(1, Math.min(15, value + 1)));
  };
  const decreasePrice = () => {
    let value = getValues("pizzaPrice");
    if (typeof value === "string") {
      value = parseInt(value);
    }
    setValue("pizzaPrice", Math.max(1, Math.min(15, value - 1)));
  };

  const onSubmit = useCallback(({ pizzaPrice }: FormFields) => {
    setPizzaPrice(pizzaPrice);
  }, []);

  useEffect(() => {
    const subscription = watch(() => {
      handleSubmit(onSubmit)();
    });
    return () => subscription.unsubscribe();
  }, [handleSubmit, watch, onSubmit]);

  return (
    <section className="my-6">
      <h3>Economy</h3>

      <div className="prose my-6 text-gray-700">
        <p>A fool and their money are soon parted.</p>
      </div>

      <h4>Pricing</h4>
      <div className="prose my-2 text-gray-700">
        <p>A higher price will lower your demand.</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)}>
        <label className="flex gap-2 items-center">
          Pizza price:
          <div className="flex items-stretch focus-within:ring-1 ring-blue-500">
            <button
              className={smallPrimaryButton}
              type="button"
              onClick={() => decreasePrice()}
            >
              -
            </button>
            <input
              type="number"
              className="inline-block w-12 font-bold border-0 focus:border-0 focus:outline-0 focus:ring-0"
              {...register("pizzaPrice")}
            />
            <button
              className={smallPrimaryButton}
              type="button"
              onClick={() => increasePrice()}
            >
              +
            </button>
          </div>
          <span className="font-bold">coins</span>
        </label>
        <span>{errors.pizzaPrice?.message}</span>
      </form>
    </section>
  );
}
export function TownCentre({ lot, lotId }: { lot: Lot; lotId: string }) {
  return (
    <div className={classnames("container", "px-2", "max-w-2xl", "mb-8")}>
      <h2>Town Centre</h2>
      <div className="flex gap-6 items-center">
        <SvgTownCentre height={100} width={100} />
        <div className="prose my-6 text-gray-700">
          <p>
            This is where tribesmice go to decide on important tribe issues.
          </p>
        </div>
      </div>
      <UpgradeSection lotId={lotId} lot={lot} />
      {lot.level >= 1 && <EconomySection />}
      <PopulationSection />
    </div>
  );
}
