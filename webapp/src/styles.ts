import { classnames, TArg } from "tailwindcss-classnames";

export const button = classnames(
  "my-2",
  "py-2",
  "px-8",
  "text-white",
  "bg-green-600",
  "disabled:bg-gray-600" as TArg,
);

export default { button };