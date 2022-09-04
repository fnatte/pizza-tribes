import React from "react";
import { QueryClientProvider } from "react-query";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import CreateAccountPage from "./CreateAccountPage";
import GamePage from "./GamePage";
import GamesPage from "./GamesPage";
import { Loading } from "./Loading";
import LoginPage from "./LoginPage";
import LogoutPage from "./LogoutPage";
import { useGames } from "./queries/useGames";
import { queryClient } from "./queryClient";

function AuthorizedPages() {
  const { data: games, error } = useGames();

  if (error) {
    return <Navigate to="/logout" />;
  }

  if (!games) {
    return <Loading />;
  }

  const joinedGames = games.filter((game) => game.joined);

  return (
    <Routes>
      <Route
        index
        element={
          <Navigate
            to={
              joinedGames.length > 0 ? `/game/${joinedGames[0].id}` : "/games"
            }
            replace
          />
        }
      />

      <Route path="/games" element={<GamesPage />} />
      <Route path="/game/:gameId/*" element={<GamePage />} />
    </Routes>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="*" element={<AuthorizedPages />} />

          <Route path="/login" element={<LoginPage />} />
          <Route path="/logout" element={<LogoutPage />} />
          <Route path="/create-account" element={<CreateAccountPage />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
