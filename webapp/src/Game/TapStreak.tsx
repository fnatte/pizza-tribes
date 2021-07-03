import React from "react";
import { usePrevious } from "react-use";
import { classnames, TArg } from "tailwindcss-classnames";

const TapStreak: React.VFC<{ value: number; max: number, className?: string|undefined }> = ({
  value,
  max,
  className,
}) => {

  const prevValue = usePrevious(value);

  return (
    <div className={classnames("flex", "space-x-0.5", "justify-center", className as TArg)}>
      {Array.from({ length: max}, (_, i) => {
        return (
          <div
            key={i}
            className={classnames("w-4", "h-4", "border", "border-gray-500", "transition-colors", {
              "bg-green-600": i < value,
              "bg-green-100": i >= value,
              ["animate-wiggle-short" as any]: prevValue !== undefined && i < value && i >= prevValue,
            })}
          />
        );
      })}
    </div>
  );
};

export default TapStreak;
