import { yupResolver } from "@hookform/resolvers/yup";
import React from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { useAsync } from "react-use";
import classnames from "classnames";
import { WorldEntry_Town } from "../../generated/world";
import * as yup from "yup";
import { RemoveIndex } from "../../utils";
import styles from "../../styles";
import { useStore } from "../../store";
import { apiFetch } from "../../api";

type Props = {
  town: WorldEntry_Town;
  x: number;
  y: number;
};

const schema = yup.object().shape({
  count: yup.number().required().positive().integer(),
});

type FormFields = RemoveIndex<yup.Asserts<typeof schema>>;

const WorldTownView: React.FC<Props> = ({ x, y, town }) => {
  const username = useAsync(async () => {
    if (town === null) {
      return null;
    }
    const response = await apiFetch(`/user/${town?.userId}`);
    if (
      !response.ok ||
      response.headers.get("Content-Type") !== "application/json"
    ) {
      throw new Error("Failed to get user");
    }

    const data = await response.json();
    return data.username as string;
  }, [town]);

  const thieves = useStore((state) => state.gameState.population.thieves);
  const thievesAvailable = thieves; // TODO: subtract thieves on mission
  const steal = useStore((state) => state.steal);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormFields>({ resolver: yupResolver(schema) });

  const navigate = useNavigate();

  const onSubmit = async (data: FormFields) => {
    steal(x, y, data.count);
    navigate("/");
  };

  return (
    <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
      <h2>{username.value ? `${username.value}'s town` : "Town"}</h2>

      {thievesAvailable === 0 && (
        <p className={classnames("max-w-sm", "text-gray-700")}>
          If you train some thieves, you can send them on a heist to other
          players towns to steal their coins.
        </p>
      )}

      {thievesAvailable > 0 && (
        <form
          className={classnames(
            "flex",
            "flex-col",
            "max-w-screen-sm",
            "items-center",
            "mx-auto",
            "my-4"
          )}
          onSubmit={handleSubmit(onSubmit)}
        >
          <label>
            <div className={classnames("p-1")}>Number of thieves to send: </div>
            <input
              type="text"
              className={classnames("my-1")}
              disabled={isSubmitting}
              defaultValue={thievesAvailable}
              {...register("count")}
            />
            {errors.count && <div className="p-2">{errors.count.message}</div>}
          </label>
          <button
            type="submit"
            className={styles.primaryButton}
            disabled={isSubmitting}
          >
            Send thieves
          </button>
        </form>
      )}
    </div>
  );
};

export default WorldTownView;
