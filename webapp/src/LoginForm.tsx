import React, {useState} from "react";
import {classnames, TArg} from "tailwindcss-classnames";

type Props = {
  onLogin: () => void;
}

const LoginForm: React.FC<Props> = ({ onLogin }) => {
  const [isLoading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    setLoading(true);
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const json = JSON.stringify(Object.fromEntries(formData));

    const result = await fetch("/api/auth/login", {
      method: "POST",
      body: json,
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      credentials: "include",
    });
    setLoading(false);
    if (result.status === 200) {
      onLogin();
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
      onSubmit={onSubmit}
    >
      <input
        type="text"
        name="username"
        placeholder="Username"
        className={classnames("my-1")}
        disabled={isLoading}
      />
      <input
        type="password"
        name="password"
        placeholder="Password"
        className={classnames("my-1")}
        disabled={isLoading}
      />
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
        disabled={isLoading}
      >
        Login
      </button>
    </form>
  );
}

export default LoginForm;
