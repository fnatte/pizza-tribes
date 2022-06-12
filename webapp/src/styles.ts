import classnames from "classnames";

const baseButton = [
  "text-white",
  "disabled:bg-gray-600",
  "disabled:cursor-default",
] as const;

export const button = [
  ...baseButton,
  "my-2",
  "py-2",
  "px-4",
  "md:px-8",
] as const;

export const smallButton = [...baseButton, "py-1", "px-2"] as const;

export const primaryButton = classnames(...button, "bg-green-600");
export const smallPrimaryButton = classnames(...smallButton, "bg-green-600");

export default { primaryButton, button };
