import React, {
  useEffect,
  useLayoutEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import classnames from "classnames";
import grassTileSrc from "images/v2/tiles/grass.png";
import forestTileSrc from "images/v2/tiles/forest.png";
import mountainTileSrc from "images/v2/tiles/mountain.png";
import town1TileSrc from "images/v2/tiles/town-1.png";
import town2TileSrc from "images/v2/tiles/town-2.png";
import town3TileSrc from "images/v2/tiles/town-3.png";

import { WorldEntry, WorldEntry_LandType } from "../../generated/world";

const getMapKey = (x: number, y: number) => `${x}:${y}`;

type TileType =
  | "grass"
  | "mountain"
  | "forest"
  | "town-1"
  | "town-2"
  | "town-3";

type Props = {
  x: number;
  y: number;
  zoom?: number;
  data: Map<string, WorldEntry | null>;
  onClick?: (x: number, y: number) => void;
  onOffsetChange?: (x: number, y: number) => void;
  className?: string;
  selection?: { x: number, y: number };
};

function useResizeObserver(element: HTMLElement | null) {
  const [result, setResult] = useState<ResizeObserverEntry>();

  const observer = useMemo(
    () =>
      new ResizeObserver((entries) => {
        if (entries[0]) {
          setResult(entries[0]);
        }
      }),
    []
  );

  useLayoutEffect(() => {
    if (!element) return;
    observer.observe(element);
    return () => {
      observer.disconnect();
    };
  }, [element]);

  return result;
}

function useImage(src: string, onLoad: (img: HTMLImageElement) => void) {
  useEffect(() => {
    const img = new Image();
    img.addEventListener("load", () => onLoad(img), false);
    img.src = src;
  }, []);
}

function useImages() {
  const [images, setImages] = useState<
    Record<TileType, HTMLImageElement | null>
  >({
    grass: null,
    forest: null,
    mountain: null,
    "town-1": null,
    "town-2": null,
    "town-3": null,
  });

  useImage(grassTileSrc, (img) =>
    setImages((images) => ({ ...images, grass: img }))
  );
  useImage(forestTileSrc, (img) =>
    setImages((images) => ({ ...images, forest: img }))
  );
  useImage(mountainTileSrc, (img) =>
    setImages((images) => ({ ...images, mountain: img }))
  );
  useImage(town1TileSrc, (img) =>
    setImages((images) => ({ ...images, "town-1": img }))
  );
  useImage(town2TileSrc, (img) =>
    setImages((images) => ({ ...images, "town-2": img }))
  );
  useImage(town3TileSrc, (img) =>
    setImages((images) => ({ ...images, "town-3": img }))
  );

  return images;
}

function HexGridCanvas({
  data,
  x: centerX,
  y: centerY,
  zoom = 1,
  onClick: onClickProp,
  onOffsetChange,
  className,
  selection,
}: Props) {
  const p = 0.2514;
  const w = Math.floor(1024 / (10 / zoom));
  const h = Math.floor(w * 0.52);

  const ref = useRef<HTMLCanvasElement>(null);

  const sizeResult = useResizeObserver(ref.current);
  const images = useImages();

  const showCoordinates = false;

  const [isMouseDown, setIsMouseDown] = useState(false);
  const [mousePos, setMousePos] = useState({ x: 0, y: 0 });
  const [offset, setOffset] = useState({
    x: centerX * w,
    y: (centerY + 0.5) * h,
  });

  useEffect(() => {
    if (sizeResult && ref.current) {
      ref.current.width = sizeResult.contentRect.width;
      ref.current.height = sizeResult.contentRect.height;
    }
  }, [sizeResult]);

  const mapOffsetX = Math.round(offset.x / (w * 0.75));
  const mapOffsetY = Math.round(offset.y / h);

  useEffect(() => {
    onOffsetChange?.(mapOffsetX, mapOffsetY);
  }, [mapOffsetX, mapOffsetY]);

  useEffect(() => {
    const canvas = ref.current;
    if (!canvas) {
      return;
    }

    const ctx = canvas.getContext("2d");
    if (!ctx) {
      return;
    }

    let raf: number;
    const onFrame = () => {
      const fromX = Math.floor((offset.x - canvas.width / 2) / (w * 0.75)) - 1;
      const toX = Math.floor((offset.x + canvas.width / 2) / (w * 0.75)) + 1;
      const fromY = Math.floor((offset.y - canvas.height / 2) / h) - 2;
      const toY = Math.floor((offset.y + canvas.height / 2) / h) + 1;

      ctx.clearRect(0, 0, canvas.width, canvas.height);
      ctx.fillStyle = "#dcfce7";
      ctx.fillRect(0, 0, canvas.width, canvas.height);

      ctx.save();
      ctx.translate(canvas.width / 2, canvas.height / 2);
      ctx.translate(-offset.x, -offset.y);

      const drawTileBottomLines = (dx: number, dy: number) => {
        ctx.beginPath();
        ctx.moveTo(dx + w, dy + h / 2);
        ctx.lineTo(dx + w * 3 * p, dy + h);
        ctx.lineTo(dx + w * p, dy + h);
        ctx.lineTo(dx, dy + h / 2);
        ctx.lineWidth = 5;
        ctx.lineCap = "round";
        ctx.strokeStyle = "#cfd8be";
        ctx.stroke();
      };

      const drawTileUpperLines = (dx: number, dy: number) => {
        ctx.beginPath();
        ctx.moveTo(dx, dy + h / 2);
        ctx.lineTo(dx + w * p, dy);
        ctx.lineTo(dx + w * 3 * p, dy);
        ctx.lineTo(dx + w, dy + h / 2);
        ctx.lineWidth = 5;
        ctx.lineCap = "round";
        ctx.strokeStyle = "#cfd8be";
        ctx.stroke();
      };

      const drawSelection = (dx: number, dy: number) => {
        ctx.beginPath();
        ctx.moveTo(dx, dy + h / 2);
        ctx.lineTo(dx + w * p, dy);
        ctx.lineTo(dx + w * 3 * p, dy);
        ctx.lineTo(dx + w, dy + h / 2);
        ctx.lineTo(dx + w * 3 * p, dy + h);
        ctx.lineTo(dx + w * p, dy + h);
        ctx.lineTo(dx, dy + h / 2);
        ctx.fillStyle = "#00000033";
        ctx.fill();
      };

      const getTile = (x: number, y: number): TileType | undefined => {
        const entry = data.get(getMapKey(x, y));
        if (entry) {
          if (entry.object.oneofKind === "town") {
            return "town-1";
          } else if (entry.landType === WorldEntry_LandType.GRASS) {
            return "grass";
          } else if (entry.landType === WorldEntry_LandType.FOREST) {
            return "forest";
          } else if (entry.landType === WorldEntry_LandType.MOUNTAIN) {
            return "mountain";
          }
        }

        return "grass";
      };

      const drawTile = (x: number, y: number) => {
        const dx = x * w * 0.75;
        const dy = h * (y + 0.5 * (x & 1));

        const tile = getTile(x, y);
        if (!tile) {
          return;
        }
        const img = images[tile];
        if (!img) {
          return;
        }

        const sw = w / 700;
        const sh = h / 360;
        const xOffset = (img.width - 700) * sw * 0.5;
        const yOffset = (img.height - 360) * sh * 0.5;
        const dw = img.width * sw;
        const dh = img.height * sh;

        drawTileUpperLines(dx, dy);
        ctx.drawImage(img, dx - xOffset, dy - yOffset, dw, dh);
        drawTileBottomLines(dx, dy);

        if (showCoordinates) {
          ctx.font = "18px serif";
          ctx.textAlign = "center";
          ctx.fillStyle = "black";
          ctx.fillText(`${x}, ${y}`, dx + w / 2, dy + h / 2);
        }

        if (selection && selection.x === x && selection.y === y) {
          drawSelection(dx, dy);
        }
      };

      for (let y = fromY; y <= toY; y++) {
        for (let x = fromX + Math.abs(fromX % 2); x <= toX; x += 2) {
          drawTile(x, y);
        }
        for (let x = fromX - Math.abs(fromX % 2) + 1; x <= toX; x += 2) {
          drawTile(x, y);
        }
      }

      ctx.restore();
    };

    raf = requestAnimationFrame(onFrame);

    return () => cancelAnimationFrame(raf);
  }, [sizeResult, offset, images, selection, w]);

  const onMouseDown = (e: React.MouseEvent) => {
    setIsMouseDown(true);
    setMousePos({ x: e.pageX, y: e.pageY });
  };
  const onTouchStart = (e: React.TouchEvent) => {
    setIsMouseDown(true);
    setMousePos({ x: e.touches[0].pageX, y: e.touches[0].pageY });
  };

  useEffect(() => {
    const onMouseUp = () => setIsMouseDown(false);
    const onTouchEnd = () => setIsMouseDown(false);

    const onMouseMove = (e: MouseEvent) => {
      if (isMouseDown) {
        setOffset(({ x, y }) => ({
          x: x + mousePos.x - e.pageX,
          y: y + mousePos.y - e.pageY,
        }));
        setMousePos({ x: e.pageX, y: e.pageY });
      }
    };

    const onTouchMove = (e: TouchEvent) => {
      if (isMouseDown) {
        setOffset(({ x, y }) => ({
          x: x + mousePos.x - e.touches[0].pageX,
          y: y + mousePos.y - e.touches[0].pageY,
        }));
        setMousePos({ x: e.touches[0].pageX, y: e.touches[0].pageY });
      }
    };

    document.body.addEventListener("mousemove", onMouseMove);
    document.body.addEventListener("mouseup", onMouseUp);
    document.body.addEventListener("touchmove", onTouchMove);
    document.body.addEventListener("touchend", onTouchEnd);
    return () => {
      document.body.removeEventListener("mousemove", onMouseMove);
      document.body.removeEventListener("mouseup", onMouseUp);
      document.body.removeEventListener("touchmove", onTouchMove);
      document.body.removeEventListener("touchend", onTouchEnd);
    };
  }, [mousePos, isMouseDown]);

  const onClick = (e: React.MouseEvent) => {
    const rect = ref.current?.getBoundingClientRect();
    if (!rect) {
      return;
    }

    // TODO: This is not really correct but it's good enough for now
    const x = e.clientX - rect.left - (rect.width / 2 - offset.x);
    const y = e.clientY - rect.top - (rect.height / 2 - offset.y);

    const tileX = Math.floor(x / (w * 0.75) - 0.2);
    const tileY = Math.floor(y / h - (Math.abs(tileX) % 2 === 1 ? 0.5 : 0));

    onClickProp?.(tileX, tileY);
  };

  return (
    <canvas
      ref={ref}
      className={classnames("touch-none", className)}
      onMouseDown={onMouseDown}
      onTouchStart={onTouchStart}
      onClick={onClick}
    />
  );
}

export default HexGridCanvas;
