export const getIdx = (x: number, y: number): { zidx: number, eidx: number } => {
  const zidx = Math.floor(y / 10) * 11 + Math.floor(x / 10);
  const eidx = (y % 10) * 10 + (x % 10);

  return { zidx, eidx };
}

