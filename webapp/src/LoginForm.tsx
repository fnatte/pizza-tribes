import React, { useState } from "react";
import { classnames, TArg } from "tailwindcss-classnames";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { RemoveIndex } from "./utils";
import { useForm } from "react-hook-form";

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

  const onSubmit = async (data: FormFields) => {
    let response: Response;
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 5000);
      response = await fetch("/api/auth/login", {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        credentials: "include",
        signal: controller.signal,
      });
      clearTimeout(timeoutId);
    } catch (e) {
      throw e;
    }

    if (response.status === 200) {
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
      <input
        type="password"
        placeholder="Password"
        className={classnames("my-1")}
        disabled={isSubmitting}
        {...register("password")}
      />
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
          "disabled:bg-gray-600" as TArg
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
