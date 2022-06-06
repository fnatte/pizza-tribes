import React, { useState } from "react";
import classnames from "classnames";
import { useStore } from "../../store";
import styles from "../../styles";
import {
  formatDurationShort,
  formatNanoTimestampToNowShort,
} from "../../utils";
import {
  ResearchDiscovery,
  ResearchNode,
  ResearchTrack,
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

const title = classnames("text-lg", "md:text-xl", "mb-2");
const label = classnames("text-xs", "md:text-sm");
const value = classnames("text-sm");
const descriptionStyle = classnames("text-sm", "text-gray-600");

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
};

const descriptions: Record<ResearchDiscovery, React.VFC | undefined> = {
  [ResearchDiscovery.WEBSITE]: () => (
    <p className={descriptionStyle}>
      If only there was some kind of online medium that could increase our
      popularity.
    </p>
  ),
  [ResearchDiscovery.DIGITAL_ORDERING_SYSTEM]: () => (
    <p className={descriptionStyle}>
      With a digital ordering system the salesmice could work more effectively.
    </p>
  ),
  [ResearchDiscovery.MOBILE_APP]: () => (
    <p className={descriptionStyle}>
      A mobile app would increase our reach even further which in turn would
      increase demand of our fine pizzas.
    </p>
  ),
  [ResearchDiscovery.MASONRY_OVEN]: () => (
    <p className={descriptionStyle}>
      If we learned how to master the traditional pizza oven our pizzas would
      taste better &mdash; and that would lead to increased demand!
    </p>
  ),
  [ResearchDiscovery.GAS_OVEN]: () => (
    <p className={descriptionStyle}>
      A gas oven would heat much faster than the traditional ones. If we had gas
      ovens we would be able to bake pizzas faster.
    </p>
  ),
  [ResearchDiscovery.HYBRID_OVEN]: () => (
    <p className={descriptionStyle}>
      If we just could get the taste of traditional masonry ovens with the speed
      of gas ovens...
    </p>
  ),
  [ResearchDiscovery.DURUM_WHEAT]: () => (
    <p className={descriptionStyle}>
      We should deepen our knowledge of durum wheat to improve taste of our
      pizzas.
    </p>
  ),
  [ResearchDiscovery.DOUBLE_ZERO_FLOUR]: () => (
    <p className={descriptionStyle}>
      Lets continue the search for the perfect dough!
    </p>
  ),
  [ResearchDiscovery.SAN_MARZANO_TOMATOES]: () => (
    <p className={descriptionStyle}>
      Our tomatoes have no taste! To improve our tomato sauce we need to find
      the best tomatoes.
    </p>
  ),
  [ResearchDiscovery.OCIMUM_BASILICUM]: () => (
    <p className={descriptionStyle}>
      A key ingredient in tomato sauce is basil. Let us learn more about the
      herb.
    </p>
  ),
  [ResearchDiscovery.EXTRA_VIRGIN]: () => (
    <p className={descriptionStyle}>
      If we could find the perfect olive oil our tomato sauce would be even
      tastier!
    </p>
  ),
};

const numberFormat = new Intl.NumberFormat();

const ResearchNodeView: React.VFC<{
  node: ResearchNode;
  discovered: boolean;
  parentDiscovered: boolean;
  opened: boolean;
  onToggleOpen: () => void;
}> = ({ node, discovered, opened, onToggleOpen, parentDiscovered }) => {
  const coins = useStore((state) => state.gameState.resources.coins);
  const startResearch = useStore((state) => state.startResearch);
  const researchQueue = useStore((state) => state.gameState.researchQueue);

  const onSubmit = () => {
    startResearch(node.discovery);
    window.scroll(0, 0);
  };

  const SvgImage = svgs[node.discovery] || PlaceholderImage;
  const Description = descriptions[node.discovery];

  const isBeingResearched = researchQueue.some(
    (x) => x.discovery === node.discovery
  );

  const canAfford = node.cost < coins;
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
      >
        <SvgImage className={classnames("w-full", "h-full")} />
      </button>
      <div className={classnames("px-2")}>
        <div className={title}>{node.title}</div>
        {Description && <Description />}
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
                <span className={value}>
                  {numberFormat.format(node.cost)} coins
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        {discovered && (
          <section className={classnames("m-4", "p-4", "bg-green-200")}>
            <span>{node.title} has already been researched.</span>
          </section>
        )}
        {isBeingResearched && (
          <section className={classnames("m-4", "p-4", "bg-green-200")}>
            <span>This is currently being researched.</span>
          </section>
        )}
        {!parentDiscovered && (
          <section className={classnames("m-4", "p-4", "bg-yellow-200")}>
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
            >
              Start Research
            </button>
          </form>
        )}
      </div>
    </div>
  );
};

const ResearchNodesView: React.VFC<{
  node: ResearchNode;
  discoveries: ResearchDiscovery[];
  parentDiscovered: boolean;
  openedDiscovery: ResearchDiscovery | null;
  setOpenedDiscovery: (x: ResearchDiscovery | null) => void;
}> = ({
  node,
  discoveries,
  openedDiscovery,
  setOpenedDiscovery,
  parentDiscovered,
}) => {
  const discovered = discoveries.includes(node.discovery);
  return (
    <div className={classnames("flex", "flex-col", "gap-4")}>
      <ResearchNodeView
        key={node.discovery}
        node={node}
        parentDiscovered={parentDiscovered}
        discovered={discovered}
        opened={openedDiscovery === node.discovery}
        onToggleOpen={() =>
          setOpenedDiscovery(
            openedDiscovery !== node.discovery ? node.discovery : null
          )
        }
      />
      <div
        className={classnames("flex", "gap-4", "justify-center", "flex-wrap")}
      >
        {node.nodes.map((n) => (
          <ResearchNodesView
            key={n.discovery}
            node={n}
            parentDiscovered={discovered}
            discoveries={discoveries}
            openedDiscovery={openedDiscovery}
            setOpenedDiscovery={setOpenedDiscovery}
          />
        ))}
      </div>
    </div>
  );
};

function recursiveLoop<TNode extends Record<"nodes", TNode[]>>(
  root: TNode,
  callback: (node: TNode) => void
) {
  const stack = [root];
  while (stack.length) {
    const node = stack.shift();
    if (node !== undefined) {
      callback(node);
      node.nodes && stack.push(...node.nodes);
    }
  }
  return null;
}

function recursiveSearch<
  TNode extends Record<"nodes", TNode[]>,
  TKey extends keyof TNode
>(tree: TNode[], value: TNode[TKey], key: TKey): TNode | null {
  const stack = [...tree];
  while (stack.length) {
    const node = stack.shift();
    if (node !== undefined) {
      if (node[key] === value) return node;
      node.nodes && stack.push(...node.nodes);
    }
  }
  return null;
}

const getTrackCounts = (
  researchTrack: ResearchTrack,
  discoveries: ResearchDiscovery[]
): { count: number; discovered: number } => {
  let count = 0;
  let discovered = 0;

  if (researchTrack.rootNode) {
    recursiveLoop(researchTrack.rootNode, (node) => {
      count++;
      if (discoveries.includes(node.discovery)) {
        discovered++;
      }
    });
  }

  return { count, discovered };
};

const findNode = (
  researchTracks: ResearchTrack[],
  discovery: ResearchDiscovery
): ResearchNode | null => {
  for (let i = 0; i < researchTracks.length; i++) {
    const node = researchTracks[i].rootNode;
    if (node !== undefined) {
      if (node.discovery === discovery) {
        return node;
      }
      const n = recursiveSearch(node.nodes, discovery, "discovery");
      if (n !== null) {
        return n;
      }
    }
  }

  return null;
};

function ResearchInstitute() {
  const researchTracks =
    useStore((state) => state.gameData?.researchTracks) || [];
  const researchQueue = useStore((state) => state.gameState.researchQueue);

  const [treeOpen, setTreeOpen] = useState<string | null>(null);

  const discoveries = useStore((state) => state.gameState.discoveries);
  const [
    openedDiscovery,
    setOpenedDiscovery,
  ] = useState<ResearchDiscovery | null>(null);

  return (
    <div className={classnames("px-2", "w-full", "max-w-2xl", "mb-8")}>
      <h2>Research Institute</h2>
      <SvgResearchInstitute height={100} width={100} />
      <p className={classnames("my-4", "text-gray-700")}>
        Looking for the next big thing? Spend some coins on research!
      </p>

      {researchQueue.length > 0 && (
        <>
          <h3>Ongoing Research</h3>
          <table>
            <tbody>
              {researchQueue.map((ongoingResearch) => (
                <tr key={ongoingResearch.discovery}>
                  <td className={classnames("p-2")}>
                    {findNode(researchTracks, ongoingResearch.discovery)?.title}
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

      <h3>Areas</h3>
      <div className={classnames("flex", "flex-col", "gap-4")}>
        {researchTracks.map((track) => {
          const trackCounts = getTrackCounts(track, discoveries);
          return (
            <div
              key={track.title}
              className={classnames("bg-green-400", "p-1")}
            >
              <div className={classnames("flex", "items-center", "p-1")}>
                <div>
                  <span className={classnames("ml-4")}>{track.title}</span>
                  <span
                    className={classnames("ml-2", "text-sm", "text-gray-800")}
                  >
                    ({trackCounts.discovered} of {trackCounts.count})
                  </span>
                </div>
                <button
                  className={classnames(
                    "p-1",
                    "bg-white",
                    "ml-auto",
                    "flex",
                    "justify-center",
                    "items-center"
                  )}
                  onClick={() =>
                    setTreeOpen((x) => (x !== track.title ? track.title : null))
                  }
                >
                  {treeOpen === track.title ? "Close" : "Open"}
                </button>
              </div>
              {treeOpen === track.title && (
                <div className={classnames("p-2", "bg-green-200")}>
                  {track.rootNode && (
                    <ResearchNodesView
                      node={track.rootNode}
                      discoveries={discoveries}
                      parentDiscovered={true}
                      openedDiscovery={openedDiscovery}
                      setOpenedDiscovery={setOpenedDiscovery}
                    />
                  )}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default ResearchInstitute;
