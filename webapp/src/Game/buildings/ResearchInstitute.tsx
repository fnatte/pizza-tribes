import React, { useState } from "react";
import { classnames } from "tailwindcss-classnames";
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
  [ResearchDiscovery.WEBSITE]: PlaceholderImage,
  [ResearchDiscovery.DIGITAL_ORDERING_SYSTEM]: PlaceholderImage,
  [ResearchDiscovery.MOBILE_APP]: PlaceholderImage,
  [ResearchDiscovery.MASONRY_OVEN]: PlaceholderImage,
  [ResearchDiscovery.GAS_OVEN]: PlaceholderImage,
  [ResearchDiscovery.HYBRID_OVEN]: PlaceholderImage,
  [ResearchDiscovery.DURUM_WHEAT]: PlaceholderImage,
  [ResearchDiscovery.DOUBLE_ZERO_FLOUR]: PlaceholderImage,
  [ResearchDiscovery.SAN_MARZANO_TOMATOES]: PlaceholderImage,
  [ResearchDiscovery.OCIMUM_BASILICUM]: PlaceholderImage,
  [ResearchDiscovery.EXTRA_VIRGIN]: PlaceholderImage,
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
  opened: boolean;
  onToggleOpen: () => void;
}> = ({ node, discovered, opened, onToggleOpen }) => {
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
  const showResearchButton = !isBeingResearched && !discovered;

  if (!opened) {
    return (
      <button
        className={classnames("flex", "w-16", "h-16", "self-center")}
        onClick={() => onToggleOpen()}
      >
        <SvgImage className={classnames("w-full", "h-full")} />
      </button>
    );
  }

  return (
    <div
      className={classnames(
        "flex",
        "justify-center",
        "mb-8",
        "flex-wrap",
        "gap-6",
        "lg:gap-8"
      )}
    >
      <button
        className={classnames("flex-shrink-0", "w-24", "h-24")}
        onClick={() => onToggleOpen()}
      >
        <SvgImage className={classnames("w-full", "h-full")} />
      </button>
      <div>
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
  openedDiscovery: ResearchDiscovery | null;
  setOpenedDiscovery: (x: ResearchDiscovery | null) => void;
}> = ({ node, discoveries, openedDiscovery, setOpenedDiscovery }) => {
  return (
    <div className={classnames("flex", "flex-col", "gap-4")}>
      <ResearchNodeView
        key={node.discovery}
        node={node}
        discovered={discoveries.includes(node.discovery)}
        opened={openedDiscovery === node.discovery}
        onToggleOpen={() =>
          setOpenedDiscovery(
            openedDiscovery !== node.discovery ? node.discovery : null
          )
        }
      />
      <div className={classnames("flex", "gap-4", "justify-center")}>
        {node.nodes.map((n) => (
          <ResearchNodesView
            key={n.discovery}
            node={n}
            discoveries={discoveries}
            openedDiscovery={openedDiscovery}
            setOpenedDiscovery={setOpenedDiscovery}
          />
        ))}
      </div>
    </div>
  );
};

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
    <div className={classnames("max-w-full", "px-2")}>
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
        {researchTracks.map((track) => (
          <div key={track.title} className={classnames("bg-green-400", "p-1")}>
            <div className={classnames("flex", "items-center", "p-1")}>
              <div className={classnames("ml-4")}>{track.title}</div>
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
                    openedDiscovery={openedDiscovery}
                    setOpenedDiscovery={setOpenedDiscovery}
                  />
                )}
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

export default ResearchInstitute;
