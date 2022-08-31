import * as React from "react";
import { Lot } from "../store";
import { ReactComponent as SvgKitchen } from "../../images/kitchen.svg";
import { ReactComponent as SvgHouse } from "../../images/house.svg";
import { ReactComponent as SvgShop } from "../../images/shop.svg";
import { ReactComponent as SvgSchool } from "../../images/school.svg";
import { ReactComponent as SvgMarketingHQ } from "../../images/marketing-hq.svg";
import { ReactComponent as SvgResearchInstitute } from "../../images/research-institute.svg";
import { ReactComponent as SvgTownCentre } from "../../images/town-centre.svg";
import { ReactComponent as SvgConstructingKitchen } from "../../images/constructing-kitchen.svg";
import { ReactComponent as SvgConstructingHouse } from "../../images/constructing-house.svg";
import { ReactComponent as SvgConstructingShop } from "../../images/constructing-shop.svg";
import { ReactComponent as SvgConstructingSchool } from "../../images/constructing-school.svg";
import { ReactComponent as SvgConstructingMarketingHQ } from "../../images/constructing-marketing-hq.svg";
import { ReactComponent as SvgConstructingResearchInstitute } from "../../images/constructing-research-institute.svg";
import { Building } from "../generated/building";
import classnames from "classnames";
import { getTapInfo } from "../utils";
import { useStore } from "../store";
import { Construction } from "../generated/gamestate";
import { useMemo, useState } from "react";
import { useInterval } from "react-use";
import { GameData } from "../generated/game_data";
import { ResearchDiscovery } from "../generated/research";

const Badge: React.FC<
  {
    position: "topleft" | "topright";
    animation?: "bounce" | undefined;
    background?: "red" | "white" | undefined;
    size?: "normal" | "big";
  } & React.HTMLAttributes<SVGElement>
> = ({
  position,
  animation,
  children,
  background = "red",
  size = "normal",
  ...rest
}) => {
  const { x, y } = position === "topleft" ? { x: -2, y: -2 } : { x: 18, y: -1 };
  const offset = useMemo(() => Math.random(), []);

  return (
    <g transform={`translate(${x}, ${y})`} {...rest}>
      {animation === "bounce" && (
        <animateTransform
          attributeName="transform"
          begin={`${(offset * 2.5).toFixed(2)}s`}
          dur="2.5s"
          type="translate"
          repeatCount="indefinite"
          fill="freeze"
          calcMode="spline"
          keySplines="0 0 1 1; 0.4 0 0.2 1; 0.4 0 0.2 1"
          keyTimes="0; .8; .85; 1"
          values={`${x} ${y}; ${x} ${y}; ${x} ${y - 2}; ${x} ${y}`}
        />
      )}
      <circle
        r={size === "normal" ? "3" : "4"}
        className={classnames({
          "fill-red-700": background === "red",
          "fill-gray-300": background === "white",
        })}
      />
      <text
        textAnchor="middle"
        dominantBaseline="central"
        className={classnames({
          "text-[4px]": size === "normal",
          "text-[5px]": size === "big",
          "text-gray-50": background === "red",
          "text-black": background === "white",
        })}
      >
        {children}
      </text>
    </g>
  );
};

const NotificationBadge: React.VFC = () => (
  <Badge position="topleft" animation="bounce" />
);

const LevelBadge: React.VFC<{ level: number }> = ({ level }) => (
  <Badge
    position="topright"
    background="white"
    size="big"
    data-cy="level-badge"
  >
    {level}
  </Badge>
);

function renderBuilding(
  building: Building | undefined,
  notification: boolean,
  level: number,
  gameData: GameData | null
) {
  const buildingInfo =
    building !== undefined ? gameData?.buildings[building] : undefined;

  switch (building) {
    case Building.KITCHEN:
      return (
        <g transform="translate(-10, -13)">
          {buildingInfo && <title>{buildingInfo.title}</title>}
          <SvgKitchen width={20} height={20} />
          {notification && <NotificationBadge />}
          <LevelBadge level={level} />
        </g>
      );
    case Building.HOUSE:
      return (
        <g transform="translate(-5, -8)">
          {buildingInfo && <title>{buildingInfo.title}</title>}
          <g transform="scale(0.5)">
            <SvgHouse width={20} height={20} />
          </g>
          <g transform="translate(-5, -2)">
            {notification && <NotificationBadge />}
            <LevelBadge level={level} />
          </g>
        </g>
      );
    case Building.SHOP:
      return (
        <g transform="translate(-10, -13)">
          {buildingInfo && <title>{buildingInfo.title}</title>}
          <SvgShop width={20} height={20} />
          {notification && <NotificationBadge />}
          <LevelBadge level={level} />
        </g>
      );
    case Building.SCHOOL:
      return (
        <g transform="translate(-10, -13)">
          {buildingInfo && <title>{buildingInfo.title}</title>}
          <SvgSchool width={20} height={20} />
          {notification && <NotificationBadge />}
          <LevelBadge level={level} />
        </g>
      );
    case Building.MARKETINGHQ:
      return (
        <g transform="translate(-10, -13)">
          {buildingInfo && <title>{buildingInfo.title}</title>}
          <SvgMarketingHQ width={20} height={20} />
          {notification && <NotificationBadge />}
          <LevelBadge level={level} />
        </g>
      );
    case Building.RESEARCH_INSTITUTE:
      return (
        <g transform="translate(-10, -13)">
          {buildingInfo && <title>{buildingInfo.title}</title>}
          <SvgResearchInstitute width={20} height={20} />
          {notification && <NotificationBadge />}
          <LevelBadge level={level} />
        </g>
      );
    case Building.TOWN_CENTRE:
      return (
        <g transform="translate(-10, -13)">
          {buildingInfo && <title>{buildingInfo.title}</title>}
          <SvgTownCentre width={20} height={20} />
          {notification && <NotificationBadge />}
          <LevelBadge level={level} />
        </g>
      );
  }
}

const renderConstructingBuilding = (building: Building | undefined) => {
  switch (building) {
    case Building.KITCHEN:
      return (
        <g transform="translate(-10, -13)">
          <SvgConstructingKitchen width={20} height={20} />
        </g>
      );
    case Building.HOUSE:
      return (
        <g transform="translate(-5, -8) scale(0.5)">
          <SvgConstructingHouse width={20} height={20} />
        </g>
      );
    case Building.SHOP:
      return (
        <g transform="translate(-10, -13)">
          <SvgConstructingShop width={20} height={20} />;
        </g>
      );
    case Building.SCHOOL:
      return (
        <g transform="translate(-10, -13)">
          <SvgConstructingSchool width={20} height={20} />
        </g>
      );
    case Building.MARKETINGHQ:
      return (
        <g transform="translate(-10, -13)">
          <SvgConstructingMarketingHQ width={20} height={20} />
        </g>
      );
    case Building.RESEARCH_INSTITUTE:
      return (
        <g transform="translate(-10, -13)">
          <SvgConstructingResearchInstitute width={20} height={20} />
        </g>
      );
  }
};

const renderLot = (
  lots: Record<string, Lot | undefined>,
  constructionQueue: Construction[],
  lotId: string,
  now: Date,
  discoveries: ResearchDiscovery[],
  gameData: GameData | null
) => {
  const lot = lots[lotId];
  const construction = constructionQueue.find((x) => x.lotId === lotId);

  if (construction && !construction.razing && construction.level <= 0) {
    return renderConstructingBuilding(construction.building);
  }

  if (lot) {
    return renderBuilding(
      lot.building,
      getTapInfo(lot, discoveries, now).canTap,
      lot.level + 1,
      gameData
    );
  }

  return null;
};

function SvgTown(
  {
    lots,
    constructionQueue,
    ...props
  }: React.SVGProps<SVGSVGElement> & {
    lots: Record<string, Lot | undefined>;
    constructionQueue: Construction[];
  },
  svgRef?: React.Ref<SVGSVGElement>
) {
  const [now, setNow] = useState(new Date());
  useInterval(() => setNow(new Date()), 10_000);

  const gameData = useStore((state) => state.gameData);
  const discoveries = useStore((state) => state.gameState.discoveries);

  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={813.744}
      height={648.868}
      viewBox="0 0 215.303 171.68"
      id="svg7459"
      ref={svgRef}
      {...props}
    >
      <g id="layer1" transform="translate(1.818 -62.327)">
        <ellipse
          id="path10"
          cx={105.833}
          cy={148.167}
          rx={107.652}
          ry={85.84}
          fill="#59a608"
          fillOpacity={1}
          fillRule="evenodd"
          strokeWidth={0.775}
        />
        <g
          id="g1059"
          transform="translate(-42.795 -61.723)"
          display="inline"
          fill="none"
          fillOpacity={1}
          stroke="#94831d"
          strokeWidth={0.794}
          strokeLinecap="butt"
          strokeLinejoin="miter"
          strokeMiterlimit={4}
          strokeDasharray="1.5875,.79375"
          strokeDashoffset={0}
          strokeOpacity={1}
        >
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M95.66 212.569c-2.493-3.823-3.825-8.58-5.291-13.23"
            id="path905"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M91.692 219.184c1.488-1.201 4.524-3.676 7.937-2.646"
            id="path926"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M184.296 205.954c-.292 2.654-.252 5.105 1.323 6.615"
            id="path964"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M197.525 219.184c-.482-1.92-1.438-3.504-3.969-3.97"
            id="path972"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M209.431 201.986c-2.357.02-3.804 2.714-5.291 5.291"
            id="path980"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M128.733 170.236v3.968"
            id="path996"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M176.358 172.882v1.322"
            id="path1004"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M130.056 249.61l1.323-3.968"
            id="path1020"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M166.477 247.191l-1.323-6.614"
            id="path1034"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M132.456 230.716c-.186.569-2.472 2.758.96 5.214 4.982 4.435-1.621 4.26-1.263 6.579"
            id="path942-5"
            vectorEffect="none"
            stopColor="#000"
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M137.372 203.535l1.323 3.969"
            id="path942"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
        </g>
        <g
          id="g1045"
          transform="translate(-43.512 -61.877)"
          display="inline"
          stroke="#94701d"
          strokeWidth={1.323}
          strokeMiterlimit={4}
          strokeDasharray="5.29167,1.32292"
          strokeDashoffset={0}
          strokeOpacity={1}
          fill="none"
          fillOpacity={1}
          strokeLinecap="butt"
          strokeLinejoin="miter"
        >
          <path
            d="M83.754 211.246c-.521.037 6.36 1.56 12.474 2.98 6.355 1.475 12.093 2.84 19.276 2.312 8.345-.614 15.693-5.92 23.813-7.938 5.206-1.293 10.514-2.86 15.875-2.646 10.9.436 20.842 8.12 31.75 7.938 6.42-.107 18.52-5.292 18.52-5.292"
            id="path895"
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M112.858 213.892c-1.11-12.776-5.907-36.648-1.01-37.078 11.2-.984 23.294-3.858 36.73-1.287 20.128 3.853 49.324-6.983 58.987 9.48 2.702 4.603-5.616 15.823-2.102 16.979"
            id="path988"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
          <path
            style={{
              fontVariationSettings: "normal",
            }}
            d="M112.858 219.184s8.684 19.966 15.451 23.667c7.604 4.159 15.998-2.695 24.21-4.731 6.002-1.489 11.89 2.87 17.383-.282 7.218-4.14 10.425-22.623 10.425-22.623"
            id="path1012"
            opacity={1}
            vectorEffect="none"
            stopColor="#000"
            stopOpacity={1}
          />
        </g>
        <g
          id="lot11"
          transform="translate(184.514 140.613)"
          data-type="lot"
          display="inline"
          data-cy="lot11"
        >
          <ellipse
            id="lot11ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "11", now, discoveries, gameData)}
          ;
        </g>
        <g
          id="lot10"
          data-cy="lot10"
          transform="translate(134.174 99.582)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot10ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "10", now, discoveries, gameData)}
        </g>
        <g
          id="lot9"
          data-cy="lot9"
          transform="translate(84.44 96.273)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot9ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "9", now, discoveries, gameData)}
        </g>
        <g
          id="lot8"
          data-cy="lot8"
          transform="translate(44.1 127.126)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot8ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "8", now, discoveries, gameData)}
        </g>
        <g
          id="lot7"
          data-cy="lot7"
          transform="translate(47.454 169.05)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot7ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "7", now, discoveries, gameData)}
        </g>
        <g
          id="lot6"
          data-cy="lot6"
          transform="translate(85.633 198.884)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot6ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "6", now, discoveries, gameData)}
        </g>
        <g
          id="lot5"
          data-cy="lot5"
          transform="translate(128.324 196.828)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot5ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "5", now, discoveries, gameData)}
        </g>
        <g
          id="lot4"
          data-cy="lot4"
          transform="translate(160.175 168.5)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot4ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "4", now, discoveries, gameData)}
        </g>
        <g
          id="lot3"
          data-cy="lot3"
          transform="translate(102.17 161.218)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot3ellipse"
            rx={16.495}
            ry={9.577}
            cx={0}
            cy={0}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "3", now, discoveries, gameData)}
        </g>
        <g
          id="lot2"
          data-cy="lot2"
          transform="translate(91.65 128.657)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot2ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "2", now, discoveries, gameData)}
        </g>
        <g
          id="lot1"
          data-cy="lot1"
          transform="translate(141.641 132.788)"
          data-type="lot"
          display="inline"
        >
          <ellipse
            id="lot1ellipse"
            cx={0}
            cy={0}
            rx={16.495}
            ry={9.577}
            fill="#87c43b"
            fillOpacity={1}
            strokeWidth={0.379}
          />
          {renderLot(lots, constructionQueue, "1", now, discoveries, gameData)}
        </g>
      </g>
    </svg>
  );
}

const ForwardRef = React.forwardRef(SvgTown);
export default ForwardRef;
