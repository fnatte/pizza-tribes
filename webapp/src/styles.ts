import classnames from "classnames";

export const button = [
  "my-2",
  "py-2",
  "px-4",
  "md:px-8",
  "text-white",
  "disabled:bg-gray-600",
  "disabled:cursor-default",
] as const;

export const primaryButton = classnames(...button, "bg-green-600");

export default { primaryButton, button };
