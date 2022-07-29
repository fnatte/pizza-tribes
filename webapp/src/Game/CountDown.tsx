import React from "react";
import { useInterval, useUpdate } from "react-use";
import { formatNanoTimestampToNowShort } from "../utils";

export function CountDown({ time }: { time: string }) {
  const update = useUpdate();
  useInterval(() => update(), 1000);

  return <>{formatNanoTimestampToNowShort(time)}</>;
}
