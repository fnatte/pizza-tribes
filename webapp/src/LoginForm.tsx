import React, { useState } from "react";
import classnames from "classnames";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { RemoveIndex } from "./utils";
import { useForm } from "react-hook-form";
import { centralApiFetch, setAccessToken } from "./api";
import { ViewFilled, ViewOffFilled } from "./icons";

type Props = {
  onLogin: () => void;
};

const schema = yup.object().shape({
  username: yup.string().required(),
  password: yup.string().required(),
});

type FormFields = RemoveIndex<yup.Asserts<typeof schema>>;

const LoginForm: React.FC<Props> = ({ onLogin }) => {
  const [badCredentials, setBadCredentials] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormFields>({ resolver: yupResolver(schema) });

  const [viewPassword, setViewPassword] = useState(false);

  const onSubmit = async (data: FormFields) => {
    let response: Response;
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 5000);
      response = await centralApiFetch("/auth/login", {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        signal: controller.signal,
      });
      clearTimeout(timeoutId);
    } catch (e) {
      throw e;
    }

    if (response.status === 200) {
      // In case we get a access token in the response we store it ourselves, note however that for regular web browser users
      // the access token will be sent as a (http-only) cookie and thus we just rely on the browser to store and transfer cookie.
      // We avoid cookies for cross-origin usages (such as mobile apps) simply because it's easier as it avoids some headaches with
      // same-origin rules and different browser vendors doing things differently.
      if (response.headers.get("Content-Type") === "application/json") {
        const json = await response.json();
        if (typeof json === "object" && typeof json.accessToken === "string") {
          setAccessToken(json.accessToken);
        }
      }
      setBadCredentials(false);
      onLogin();
    } else if (response.status === 403) {
      setBadCredentials(true);
    }
  };

  return (
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
      <input
        type="text"
        placeholder="Username"
        className={classnames("my-1")}
        disabled={isSubmitting}
        {...register("username")}
      />
      {errors.username && (
        <span className="pb-2 text-red-900">{errors.username.message}</span>
      )}
      <div className="relative">
        <input
          type={viewPassword ? "text" : "password"}
          placeholder="Password"
          className={classnames("my-1")}
          disabled={isSubmitting}
          {...register("password")}
        />
        {!viewPassword ? (
          <ViewOffFilled
            className="absolute right-2 top-1/2 -translate-y-1/2 w-5 h-5"
            onClick={() => setViewPassword(true)}
          />
        ) : (
          <ViewFilled
            className="absolute right-2 top-1/2 -translate-y-1/2 w-5 h-5"
            onClick={() => setViewPassword(false)}
          />
        )}
      </div>
      {errors.password && (
        <span className="pb-2 text-red-900">{errors.password.message}</span>
      )}
      <button
        type="submit"
        className={classnames(
          "my-2",
          "py-2",
          "px-8",
          "text-white",
          "bg-green-600",
          "disabled:bg-gray-600"
        )}
        disabled={isSubmitting}
      >
        Login
      </button>
      {badCredentials && (
        <div className={classnames("text-red-900")}>
          Wrong username or password
        </div>
      )}
    </form>
  );
};

export default LoginForm;
