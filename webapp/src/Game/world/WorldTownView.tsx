import { yupResolver } from "@hookform/resolvers/yup";
import React from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { useAsync } from "react-use";
import classnames from "classnames";
import { ReactComponent as HeartsSvg } from "images/hearts.svg";
import { WorldEntry_Town } from "../../generated/world";
import * as yup from "yup";
import { RemoveIndex } from "../../utils";
import styles from "../../styles";
import { useStore } from "../../store";
import { apiFetch } from "../../api";
import { useEducationCount } from "../useEducationCount";
import { Education } from "../../generated/education";
import { ApiUserResponse } from "../../generated/responses";
import { MouseImage } from "../components/MouseImage";
import { useGameNavigate } from "../useGameNavigate";

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
  const user = useAsync(async () => {
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

    return ApiUserResponse.fromJsonString(await response.text());
  }, [town]);

  const educationCount = useEducationCount();
  const thieves = educationCount[Education.THIEF] ?? 0;
  const thievesAvailable = thieves; // TODO: subtract thieves on mission
  const steal = useStore((state) => state.steal);

  const schema2 = schema.shape({
    count: schema.fields.count.max(
      thievesAvailable,
      "Cannot send more thieves than available."
    ),
  });

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormFields>({ resolver: yupResolver(schema2) });

  const navigate = useGameNavigate();

  const onSubmit = async (data: FormFields) => {
    steal(x, y, data.count);
    navigate("town");
  };

  if (user.loading || !user.value) {
    return (
      <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
        <HeartsSvg />
      </div>
    );
  }

  return (
    <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
      <h2>{`${user.value.username}'s town`}</h2>

      {user.value.ambassador && (
        <MouseImage
          appearance={user.value.ambassador.appearance}
          shiftRight
          className="my-6"
          data-cy="ambassador-mouse"
        />
      )}

      {thievesAvailable === 0 && (
        <p className={classnames("max-w-sm", "text-gray-700", "my-4")}>
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
              data-cy="thieves-to-send-input"
              {...register("count")}
            />
            {errors.count && (
              <div className="pb-2 text-red-900" data-cy="error">
                {errors.count.message}
              </div>
            )}
          </label>
          <button
            type="submit"
            className={styles.primaryButton}
            disabled={isSubmitting}
            data-cy="send-thieves-button"
          >
            Send thieves
          </button>
        </form>
      )}
    </div>
  );
};

export default WorldTownView;
