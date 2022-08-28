import React, { ReactNode, useCallback } from "react";
import classnames from "classnames";
import { useStore } from "../store";
import { Link, useNavigate, useParams } from "react-router-dom";
import { MouseEditor } from "./components/MouseEditor";
import { MouseAppearance } from "../generated/appearance";

function Container({ children }: { children?: ReactNode | undefined }) {
  return (
    <div
      className={classnames(
        "container",
        "max-w-screen-lg",
        "mx-auto",
        "mt-2",
        "lg:px-4"
      )}
    >
      {children}
    </div>
  );
}

export default function MouseAppearanceView() {
  const { id } = useParams();
  const navigate = useNavigate();

  const mouse = useStore(
    useCallback(
      (state) => (id !== undefined ? state.gameState.mice[id] : undefined),
      [id]
    )
  );

  const saveMouseAppearance = useStore((state) => state.saveMouseAppearance);

  const gameData = useStore((state) => state.gameData);

  if (!gameData) {
    return null;
  }

  if (!mouse || !id) {
    return <Container>Could not find mouse.</Container>;
  }

  const { name } = mouse;

  const onCancel = () => {
    navigate(-1);
  };

  const onSave = (appearance: MouseAppearance) => {
    saveMouseAppearance(id, appearance);
    navigate(`/mouse/${id}`);
  };

  return (
    <Container>
      <Link to={`/mouse/${id}`}>
        <h2 className="text-center my-4">Mouse: {name}</h2>
      </Link>
      <div className="flex gap-2">
        <MouseEditor
          className="w-full"
          onCancel={onCancel}
          onSave={onSave}
          defaultAppearance={mouse.appearance}
        />
      </div>
    </Container>
  );
}
