import React from "react";
import { classnames, TArg } from "tailwindcss-classnames";

const TapStreak: React.VFC<{ value: number; max: number, className?: string|undefined }> = ({
  value,
  max,
  className,
}) => {
  return (
    <div className={classnames("flex", "space-x-0.5", "justify-center", className as TArg)}>
      {Array.from({ length: max}, (_, i) => {
        return (
          <div
            className={classnames("w-4", "h-4", "border", "border-gray-500", {
              "bg-green-600": i < value,
              "bg-green-100": i >= value,
            })}
          />
        );
      })}
    </div>
  );
};

export default TapStreak;
