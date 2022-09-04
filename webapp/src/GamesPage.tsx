import React from "react";
import classnames from "classnames";
import { useNavigate } from "react-router-dom";
import Header from "./Header";
import { useGames } from "./queries/useGames";
import { ArrowRight } from "./icons";
import { centralApiFetch } from "./api";
import { partition } from "lodash";

async function joinGame(gameId: string): Promise<void> {
  const response = await centralApiFetch(`/games/${gameId}/join`, {
    method: "POST",
  });
  if (!response.ok) {
    throw new Error("Failed to join game");
  }
  return;
}

function GamesPage() {
  const { data: games } = useGames();
  const navigate = useNavigate();

  const handleJoinNewGame = async (gameId: string) => {
    await joinGame(gameId);
    navigate(`/game/${gameId}`);
  };

  const handleContinueGame = async (gameId: string) => {
    navigate(`/game/${gameId}`);
  };

  const [joinedGames, unjoinedGames] = partition(games, (game) => game.joined);

  return (
    <div
      className={classnames(
        "flex",
        "flex-col",
        "justify-center",
        "items-center"
      )}
    >
      <Header />
      <div className={classnames("text-2xl", "mt-5")}>🍕🍕🍕🍕</div>

      {joinedGames.length > 0 && (
        <>
          <h2 className="mt-12">Continue:</h2>
          <ul className="mt-2 w-3/4 max-w-xs">
            {joinedGames.map((game) => (
              <li key={game.id} className="mb-4">
                <button
                  className="bg-green-600 w-full text-white fill-white font-bold p-4 flex items-center"
                  onClick={() => handleContinueGame(game.id)}
                >
                  <span>{game.title}</span>
                  <ArrowRight className="ml-auto" />
                </button>
              </li>
            ))}
          </ul>
        </>
      )}

      <h2 className="mt-16">Join a New Game</h2>
      {unjoinedGames.length > 0 ? (
        <ul className="mt-4 w-3/4 max-w-xs">
          {unjoinedGames.map((game) => (
            <li key={game.id} className="mb-4">
              <button
                className="bg-green-600 w-full text-white fill-white font-bold p-4 flex items-center"
                onClick={() => handleJoinNewGame(game.id)}
              >
                <span>{game.title}</span>
                <ArrowRight className="ml-auto" />
              </button>
            </li>
          ))}
        </ul>
      ) : (
        <div className="mt-2 mx-4 text-gray-700">
          There are no ongoing games that you can join.
        </div>
      )}
    </div>
  );
}

export default GamesPage;
