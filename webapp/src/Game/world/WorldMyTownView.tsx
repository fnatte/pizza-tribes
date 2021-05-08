import React from "react";
import { classnames } from "tailwindcss-classnames";
import {useStore} from "../../store";

type Props = {
};

const WorldMyTownView: React.FC<Props> = () => {
  const username = useStore(state => state.user?.username);
  return (
    <div className={classnames("flex", "items-center", "flex-col", "mt-2")}>
      <h2>{username ? `${username}'s town` : "Your town"}</h2>
        <p>This is your own town</p>
    </div>
  );
};

export default WorldMyTownView;
