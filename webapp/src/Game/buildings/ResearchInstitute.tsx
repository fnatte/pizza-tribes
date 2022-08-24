import React, { useState } from "react";
import classnames from "classnames";
import { Link, Route, Routes, useNavigate, useParams } from "react-router-dom";
import { useStore } from "../../store";
import styles, { button, primaryButton } from "../../styles";
import {
  formatDurationShort,
  formatNanoTimestampToNowShort,
  formatNumber,
} from "../../utils";
import { ReactComponent as SvgArrowRight } from "images/icons/arrow-right.svg";
import {
  ResearchDiscovery,
  ResearchInfo,
  ResearchTree,
} from "../../generated/research";
import { ReactComponent as SvgResearchInstitute } from "../../../images/research-institute.svg";
import { ReactComponent as SvgWebsite } from "../../../images/research/website.svg";
import { ReactComponent as SvgDigitalOrderingSystem } from "../../../images/research/digital-ordering-system.svg";
import { ReactComponent as SvgMobileApp } from "../../../images/research/mobile-app.svg";
import { ReactComponent as SvgMasoneryOven } from "../../../images/research/masonery-oven.svg";
import { ReactComponent as SvgGasOven } from "../../../images/research/gas-oven.svg";
import { ReactComponent as SvgHybridOven } from "../../../images/research/hybrid-oven.svg";
import { ReactComponent as SvgDurumWheat } from "../../../images/research/durum-wheat.svg";
import { ReactComponent as SvgDoubleZeroFlour } from "../../../images/research/double-zero-flour.svg";
import { ReactComponent as SvgSanMaraznoTomatoes } from "../../../images/research/san-marzano-tomatoes.svg";
import { ReactComponent as SvgOcimumBasilicum } from "../../../images/research/ocimum-basilicum.svg";
import { ReactComponent as SvgExtraVirgin } from "../../../images/research/extra-virgin.svg";
import { ReactComponent as SvgCheck } from "../../../images/icons/check.svg";
import { uniq } from "lodash";
import { Coin, GeniusFlash, Pizza } from "../../icons";
import ReactMarkdown from "react-markdown";

const title = classnames("text-lg", "md:text-xl", "mb-2");
const label = classnames("text-xs", "md:text-sm");
const value = classnames("text-sm");
const descriptionStyle = classnames("text-sm", "text-gray-800");

const PlaceholderImage: React.VFC<React.SVGProps<SVGSVGElement>> = (props) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="50"
    height="50"
    viewBox="0 0 50 50"
    {...props}
  >
    <rect fill="#ddd" width="50" height="50" />
    <text
      fill="rgba(0,0,0,0.5)"
      fontFamily="sans-serif"
      fontSize="10"
      dy="10.5"
      fontWeight="bold"
      x="50%"
      y="50%"
      textAnchor="middle"
    >
      50x50
    </text>
  </svg>
);

const svgs: Record<ResearchDiscovery, React.VFC | undefined> = {
  [ResearchDiscovery.WEBSITE]: SvgWebsite,
  [ResearchDiscovery.DIGITAL_ORDERING_SYSTEM]: SvgDigitalOrderingSystem,
  [ResearchDiscovery.MOBILE_APP]: SvgMobileApp,
  [ResearchDiscovery.MASONRY_OVEN]: SvgMasoneryOven,
  [ResearchDiscovery.GAS_OVEN]: SvgGasOven,
  [ResearchDiscovery.HYBRID_OVEN]: SvgHybridOven,
  [ResearchDiscovery.DURUM_WHEAT]: SvgDurumWheat,
  [ResearchDiscovery.DOUBLE_ZERO_FLOUR]: SvgDoubleZeroFlour,
  [ResearchDiscovery.SAN_MARZANO_TOMATOES]: SvgSanMaraznoTomatoes,
  [ResearchDiscovery.OCIMUM_BASILICUM]: SvgOcimumBasilicum,
  [ResearchDiscovery.EXTRA_VIRGIN]: SvgExtraVirgin,

  [ResearchDiscovery.SLAM]: PlaceholderImage,
  [ResearchDiscovery.CARDIO]: PlaceholderImage,
  [ResearchDiscovery.COFFEE]: PlaceholderImage,
  [ResearchDiscovery.HIT_IT]: PlaceholderImage,
  [ResearchDiscovery.TIP_TOE]: PlaceholderImage,
  [ResearchDiscovery.ON_A_ROLL]: PlaceholderImage,
  [ResearchDiscovery.TRIP_WIRE]: PlaceholderImage,
  [ResearchDiscovery.GODS_TOUCH]: PlaceholderImage,
  [ResearchDiscovery.GRAND_SLAM]: PlaceholderImage,
  [ResearchDiscovery.WHITEBOARD]: PlaceholderImage,
  [ResearchDiscovery.BIG_POCKETS]: PlaceholderImage,
  [ResearchDiscovery.CONSECUTIVE]: PlaceholderImage,
  [ResearchDiscovery.LASER_ALARM]: PlaceholderImage,
  [ResearchDiscovery.NIGHTS_WATCH]: PlaceholderImage,
  [ResearchDiscovery.SHADOW_EXPERT]: PlaceholderImage,
  [ResearchDiscovery.STRESS_HANDLING]: PlaceholderImage,
  [ResearchDiscovery.BOOTS_OF_HASTE]: PlaceholderImage,
  [ResearchDiscovery.KITCHEN_STRATEGY]: PlaceholderImage,
  [ResearchDiscovery.THIEVES_FAVORITE_BAG]: PlaceholderImage,
};

const numberFormat = new Intl.NumberFormat();

const ResearchNodeView: React.VFC<{
  node: ResearchInfo;
  discovered: boolean;
  parentDiscovered: boolean;
  opened: boolean;
  onToggleOpen: () => void;
}> = ({ node, discovered, opened, onToggleOpen, parentDiscovered }) => {
  const geniusFlashes = useStore((state) => state.gameState.geniusFlashes ?? 0);
  const startResearch = useStore((state) => state.startResearch);
  const researchQueue = useStore((state) => state.gameState.researchQueue);

  const onSubmit = () => {
    startResearch(node.discovery);
    window.scroll(0, 0);
  };

  const SvgImage = svgs[node.discovery] || PlaceholderImage;

  const isBeingResearched = researchQueue.some(
    (x) => x.discovery === node.discovery
  );

  const canAfford = geniusFlashes > 0;
  const disabled = !canAfford;
  const showResearchButton =
    !isBeingResearched && !discovered && parentDiscovered;

  if (!opened) {
    return (
      <div className={classnames("self-center", "flex", "relative")}>
        <button
          className={classnames(
            "flex",
            "w-16",
            "h-16",
            "border-2",
            "rounded-md",
            "border-black",
            "bg-yellow-50",
            "p-1"
          )}
          onClick={() => onToggleOpen()}
          aria-expanded={opened}
          data-cy="research-node-expand-toggle"
        >
          <SvgImage className={classnames("w-full", "h-full")} />
        </button>
        {discovered && (
          <SvgCheck
            className={classnames(
              "w-8",
              "h-8",
              "absolute",
              "-right-12",
              "top-4"
            )}
          />
        )}
      </div>
    );
  }

  return (
    <div
      className={classnames(
        "flex",
        "flex-col",
        "items-center",
        "mb-8",
        "gap-6",
        "lg:gap-8",
        "bg-green-100",
        "p-2"
      )}
    >
      <button
        className={classnames(
          "flex-shrink-0",
          "w-24",
          "h-24",
          "bg-yellow-100",
          "border-4",
          "rounded-md",
          "border-black",
          "p-1"
        )}
        onClick={() => onToggleOpen()}
        aria-expanded={opened}
        data-cy="research-node-expand-toggle"
      >
        <SvgImage className={classnames("w-full", "h-full")} />
      </button>
      <div className={classnames("px-2")} data-cy="research-node-item">
        <div className={title}>{node.title}</div>
        <p>{node.description}</p>
        <table>
          <tbody>
            <tr>
              <td className={classnames("px-2")}>
                <span className={label}>Research time:</span>
              </td>
              <td className={classnames("px-2")}>
                <span className={value}>
                  {formatDurationShort(node.researchTime)}
                </span>
              </td>
            </tr>
            <tr>
              <td className={classnames("px-2")}>
                <span className={label}>Cost:</span>
              </td>
              <td className={classnames("px-2")}>
                <span className={value}>1 genius flash</span>
              </td>
            </tr>
          </tbody>
        </table>
        {discovered && (
          <section
            className={classnames("m-4", "p-4", "bg-green-200")}
            data-cy="research-node-already-researched"
          >
            <span>{node.title} has already been researched.</span>
          </section>
        )}
        {isBeingResearched && (
          <section
            className={classnames("m-4", "p-4", "bg-green-200")}
            data-cy="research-node-being-researched"
          >
            <span>This is currently being researched.</span>
          </section>
        )}
        {!parentDiscovered && (
          <section
            className={classnames("m-4", "p-4", "bg-yellow-200")}
            data-cy="research-node-parent-not-discovered"
          >
            <span>
              You must research all previous items in this tree to unlock this
              one.
            </span>
          </section>
        )}
        {showResearchButton && (
          <form
            className={classnames("my-2")}
            onSubmit={(e) => {
              e.preventDefault();
              onSubmit();
            }}
          >
            <button
              type="submit"
              className={classnames(styles.primaryButton)}
              disabled={disabled}
              data-cy="research-node-start-research-button"
            >
              Start Research
            </button>
          </form>
        )}
      </div>
    </div>
  );
};
function getAreaName(area: ResearchTree): string {
  switch (area) {
    case ResearchTree.DEMAND:
      return "Demand";
    case ResearchTree.GUARDS:
      return "Guards";
    case ResearchTree.TAPPING:
      return "Tapping";
    case ResearchTree.THIEVES:
      return "Thieves";
    case ResearchTree.PRODUCTION:
      return "Production";
    default:
      return "-";
  }
}

function getTrackCounts(
  area: ResearchTree,
  research: { [idx: number]: ResearchInfo },
  discoveries: ResearchDiscovery[]
): { count: number; discovered: number } {
  const discovered = discoveries.filter((x) => research[x].tree === area)
    .length;
  const count = Object.values(research).filter((x) => x.tree === area).length;

  return { count, discovered };
}

function GeniusFlashesSection() {
  const geniusFlashes = useStore((state) => state.gameState.geniusFlashes ?? 0);
  const discoveries = useStore((state) => state.gameState.discoveries);
  const costs = useStore((state) => state.gameData?.geniusFlashCosts);
  const buyGeniusFlash = useStore((state) => state.buyGeniusFlash);
  const coins = useStore((state) => state.gameState.resources?.coins);
  const pizzas = useStore((state) => state.gameState.resources?.pizzas);

  const current = geniusFlashes + discoveries.length;

  const currentCost = costs?.[current];
  const upcomingCost = costs?.[current + 1];

  const canAfford =
    coins !== undefined &&
    pizzas !== undefined &&
    currentCost !== undefined &&
    currentCost.coins <= coins &&
    currentCost.pizzas <= pizzas;

  return (
    <section>
      <h3>Genius Flashes</h3>
      <div
        className={classnames(
          "m-4",
          "p-4",
          "bg-green-200",
          "flex",
          "items-center",
          "justify-center",
          "gap-1"
        )}
      >
        <span className="text-md text-center mr-2">
          Available
          <br />
          Genius Flashes
        </span>
        <GeniusFlash className={"h-[3em] w-[3em]"} />
        <span className="text-2xl" data-cy="available-genius-flashes">{geniusFlashes}</span>
      </div>

      {currentCost && (
        <div className="flex flex-col items-center m-4 p-4 bg-green-200">
          <span className="text-md text-center">Next up</span>
          <div className="flex items-center justify-center gap-8 mt-4">
            <GeniusFlash className={"h-[3em] w-[3em]"} />
            <div className="text-xl">
              <div className="flex items-center gap-1">
                <Coin className={"h-[1.25em] w-[1.25em]"} />
                {formatNumber(currentCost.coins)}
              </div>
              <div className="flex items-center gap-1">
                <Pizza className={"h-[1.25em] w-[1.25em]"} />
                {formatNumber(currentCost.pizzas)}
              </div>
            </div>
          </div>
          <div className="mt-4">
            <button
              className={classnames(primaryButton, "w-28")}
              onClick={() => buyGeniusFlash(current)}
              disabled={!canAfford}
              data-cy="buy-genius-flash-button"
            >
              Buy
            </button>
          </div>
        </div>
      )}

      {upcomingCost && (
        <div className="flex flex-col items-center m-4 p-4 bg-green-200">
          <span className="text-md text-center">Upcoming</span>
          <div className="flex items-center justify-center gap-8 mt-4">
            <GeniusFlash className={"h-[3em] w-[3em]"} />
            <div className="text-xl">
              <div className="flex items-center gap-1">
                <Coin className={"h-[1.25em] w-[1.25em]"} />
                {formatNumber(upcomingCost.coins)}
              </div>
              <div className="flex items-center gap-1">
                <Pizza className={"h-[1.25em] w-[1.25em]"} />
                {formatNumber(upcomingCost.pizzas)}
              </div>
            </div>
          </div>
        </div>
      )}
    </section>
  );
}

function isArea(area: number): area is ResearchTree {
  return [
    ResearchTree.PRODUCTION,
    ResearchTree.THIEVES,
    ResearchTree.TAPPING,
    ResearchTree.GUARDS,
    ResearchTree.DEMAND,
  ].includes(area);
}

function ResearchSection() {
  const [selected, setSelected] = useState<ResearchInfo | null>(null);
  const research = useStore((state) => state.gameData?.research) || [];
  const researchQueue = useStore((state) => state.gameState.researchQueue);
  const discoveries = useStore((state) => state.gameState.discoveries);
  const geniusFlashes = useStore((state) => state.gameState.geniusFlashes);
  const startResearch = useStore((state) => state.startResearch);
  const { area: areaParam } = useParams();
  const navigate = useNavigate();

  if (!areaParam) {
    return null;
  }
  const area = parseInt(areaParam);
  if (!isArea(area)) {
    return null;
  }

  const nodes = Object.values(research).filter(
    (x) => x.tree === area && (x.x !== 0 || x.y !== 0)
  );

  const hasRequirementsForSelected = Boolean(
    selected && selected.requirements.every((x) => discoveries.includes(x))
  );

  const isResearchingSelected = Boolean(
    selected && researchQueue.some((x) => x.discovery === selected.discovery)
  );

  const hasResearchedSelected = Boolean(
    selected && discoveries.includes(selected.discovery)
  );

  return (
    <section>
      <h3>{getAreaName(area)}</h3>
      <div className="w-[340px] h-[750px] mt-6 mx-auto relative">
        <svg className="w-full h-full absolute">
          {nodes.flatMap((r) => {
            const req = r.requirements
              .map((x) => research[x])
              .filter((x) => x !== undefined);
            const to = r;
            const hasRequirements = Boolean(
              to && to.requirements.every((x) => discoveries.includes(x))
            );
            return req.map((from) => (
              <path
                key={`${from.discovery}-${to.discovery}`}
                d={`M${from.x + 50},${from.y + 50} L${to.x + 50},${to.y + 50}`}
                className={classnames("stroke-[4]", {
                  "stroke-green-600": hasRequirements,
                  "stroke-gray-600": !hasRequirements,
                })}
              />
            ));
          })}
        </svg>
        {nodes.map((r) => {
          const isResearched = discoveries.includes(r.discovery);

          return (
            <button
              key={r.discovery}
              onClick={() => setSelected(r)}
              className={classnames(
                "w-[100px]",
                "h-[100px]",
                "rounded-lg",
                "absolute",
                "border",
                "flex",
                "text-center",
                "justify-center",
                "items-center",
                "font-bold",
                "p-2",
                {
                  "ring-4 ring-offset-4 ring-green-600": r === selected,
                  "bg-green-50 text-black border-gray-400": !isResearched,
                  "bg-green-600 text-white font-bold border-green-600": isResearched,
                }
              )}
              style={{ transform: `translate3d(${r.x}px, ${r.y}px, 0)` }}
              data-cy="research-node"
            >
              {r.title}
            </button>
          );
        })}
      </div>
      {selected !== null && (
        <div className="fixed bottom-0 left-0 right-0 p-4 bg-green-50 shadow-[0_0_16px_0_rgba(0,0,0,0.3)]">
          <h4>{selected.title}</h4>
          <ReactMarkdown className={"prose text-black text-sm md:text-lg"}>
            {selected.description}
          </ReactMarkdown>
          <div className="flex justify-center mt-4 gap-8">
            {selected.rewards.map((reward) => (
              <div className="text-black text-center" key={reward.attribute}>
                <div className="text-md text-gray-800">{reward.attribute}</div>
                <div className="text-3xl">{reward.value}</div>
              </div>
            ))}
          </div>
          <div className="flex mt-4 gap-8 justify-center items-center">
            <button
              className={classnames(...button, "bg-cyan-600")}
              onClick={() => setSelected(null)}
            >
              Cancel
            </button>
            {hasResearchedSelected ? (
              <span>Already researched</span>
            ) : isResearchingSelected ? (
              <span>Researching...</span>
            ) : (
              <button
                className={primaryButton}
                onClick={() => {
                  startResearch(selected.discovery);
                  navigate("..");
                }}
                disabled={
                  isResearchingSelected ||
                  hasResearchedSelected ||
                  !hasRequirementsForSelected ||
                  geniusFlashes <= 0
                }
                data-cy="start-research-button"
              >
                Research
              </button>
            )}
          </div>
        </div>
      )}
    </section>
  );
}

function MainSection() {
  const geniusFlashes = useStore((state) => state.gameState.geniusFlashes ?? 0);
  const research = useStore((state) => state.gameData?.research) || [];
  const researchQueue = useStore((state) => state.gameState.researchQueue);
  const discoveries = useStore((state) => state.gameState.discoveries);

  const areas = uniq(Object.values(research).map((x) => x.tree));

  return (
    <section>
      {researchQueue.length > 0 && (
        <>
          <h3>Ongoing Research</h3>
          <table>
            <tbody>
              {researchQueue.map((ongoingResearch) => (
                <tr
                  key={ongoingResearch.discovery}
                  data-cy="ongoing-research-row"
                >
                  <td className={classnames("p-2")}>
                    {research[ongoingResearch.discovery]?.title}
                  </td>
                  <td className={classnames("p-2")}>
                    {formatNanoTimestampToNowShort(ongoingResearch.completeAt)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}

      <h3>Genius Flashes</h3>
      <section
        className={classnames(
          "m-4",
          "p-4",
          "bg-green-200",
          "flex",
          "items-center",
          "justify-center",
          "gap-1"
        )}
      >
        <GeniusFlash className={"h-[3em] w-[3em]"} />
        <span className="text-2xl">{geniusFlashes}</span>
        <Link to="genius-flashes">
          <button className={classnames(primaryButton, "ml-8")} data-cy="get-more-button">
            Get more
          </button>
        </Link>
      </section>

      <h3>Areas</h3>
      <div className={classnames("flex", "flex-col", "gap-4")}>
        {areas.map((area) => {
          const trackCounts = getTrackCounts(area, research, discoveries);
          return (
            <Link to={`research/${area}`} key={area}>
              <button
                className={classnames(
                  "bg-green-400",
                  "p-1",
                  "inline-block",
                  "w-full"
                )}
                data-cy="research-area"
              >
                <div className={classnames("flex", "items-center", "p-1")}>
                  <div>
                    <span className={classnames("ml-4")}>
                      {getAreaName(area)}
                    </span>
                    <span
                      className={classnames("ml-2", "text-sm", "text-gray-800")}
                    >
                      ({trackCounts.discovered} of {trackCounts.count})
                    </span>
                  </div>
                  <SvgArrowRight className="ml-auto" />
                </div>
              </button>
            </Link>
          );
        })}
      </div>
    </section>
  );
}

function ResearchInstitute() {
  return (
    <div className={classnames("px-2", "w-full", "max-w-2xl", "mb-8")}>
      <h2>Research Institute</h2>
      <Link to="">
        <SvgResearchInstitute height={100} width={100} />
      </Link>
      <p className={classnames("my-4", "text-gray-700")}>
        Looking for the next big thing? Spend some coins on research!
      </p>

      <Routes>
        <Route index element={<MainSection />} />
        <Route path="genius-flashes" element={<GeniusFlashesSection />} />
        <Route path="research/:area" element={<ResearchSection />} />
      </Routes>
    </div>
  );
}

export default ResearchInstitute;
