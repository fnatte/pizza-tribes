import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { Link } from "react-router-dom";
import { classnames, TArg } from "tailwindcss-classnames";
import Header from "./Header";
import styles from "./styles";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";

const schema = yup.object().shape({
  username: yup.string().required(),
  password: yup.string().required(),
  confirm: yup
    .string()
    .required()
    .oneOf([yup.ref("password")], "Passwords must match!"),
});

type FormFields = yup.Asserts<typeof schema>;

const CreateAccountPage: React.VFC<{}> = () => {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<FormFields>({ resolver: yupResolver(schema) });

  const onSubmit = async (data: FormFields) => {
    console.log(data);
  };

  return (
    <div>
      <Header />
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
        {errors.username && <span className="pb-2">{errors.username.message}</span>}
        <input
          type="password"
          placeholder="Password"
          className={classnames("my-1")}
          disabled={isSubmitting}
          {...register("password")}
        />
        {errors.password && <span className="pb-2">{errors.password.message}</span>}
        <input
          type="password"
          placeholder="Confirm password"
          className={classnames("my-1")}
          disabled={isSubmitting}
          {...register("confirm")}
        />
        {errors.confirm && <span className="pb-2">{errors.confirm.message}</span>}
        <button type="submit" className={styles.button} disabled={isSubmitting}>
          Create Account
        </button>
      </form>

      <nav className={classnames("mt-10", "flex", "justify-center")}>
        <Link to="/">Back</Link>
      </nav>
    </div>
  );
};

export default CreateAccountPage;
