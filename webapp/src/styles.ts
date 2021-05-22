import { classnames, TArg } from "tailwindcss-classnames";

export const button: TArg[] = [
  "my-2",
  "py-2",
  "px-4",
  "md:px-8",
  "text-white",
  "disabled:bg-gray-600" as TArg,
  "disabled:cursor-default" as TArg,
];

export const primaryButton = classnames(...button, "bg-green-600");

export default { primaryButton, button };
