import React from "react";
import classnames from "classnames";
import { useStore } from "../../store";
import { MouseImage } from "../components/MouseImage";

type Props = {};

const WorldMyTownView: React.FC<Props> = () => {
  const username = useStore((state) => state.user?.username);
  const ambassador = useStore(
    (state) => state.gameState.mice[state.gameState.ambassadorMouseId]
  );

  return (
    <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
      <h2>{username ? `${username}'s town` : "Your town"}</h2>
      <p>This is your own town</p>
      {ambassador && (
        <MouseImage
          appearance={ambassador.appearance}
          className="mt-6 translate-x-16"
          data-cy="ambassador-mouse"
        />
      )}
    </div>
  );
};

export default WorldMyTownView;
