import { useCallback, useEffect, useState } from "react";

export const events = [
  "mousemove",
  "keydown",
  "wheel",
  "DOMMouseScroll",
  "mousewheel",
  "mousedown",
  "touchstart",
  "touchmove",
  "MSPointerDown",
  "MSPointerMove",
  "visibilitychange",
];

export function useActivity(
  callback: () => void,
  options?: { enabled: boolean }
) {
  const { enabled } = options ?? { enabled: true };
  const [lastReport, setLastReport] = useState(0);
  const [latestActivity, setLatestActivity] = useState(Date.now());
  const handler = useCallback(() => setLatestActivity(Date.now()), []);

  useEffect(() => {
    events.forEach((ev) => {
      document.addEventListener(ev, handler);
    });

    return () => {
      events.forEach((ev) => document.removeEventListener(ev, handler));
    };
  });

  useEffect(() => {
    if (enabled && latestActivity - lastReport > 60 * 1000) {
      setLastReport(Date.now());
      callback();
    }
  }, [enabled, latestActivity, lastReport, callback]);
}
