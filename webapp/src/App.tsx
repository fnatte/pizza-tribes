import React from "react";
import {
  BrowserRouter,
  Routes,
  Route,
} from "react-router-dom";
import CreateAccountPage from "./CreateAccountPage";
import GamePage from "./GamePage";
import LoginPage from "./LoginPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/*" element={<GamePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/create-account" element={<CreateAccountPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
