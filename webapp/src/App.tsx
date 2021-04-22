import React from "react";
import { classnames } from "tailwindcss-classnames";
import LoginForm from "./LoginForm";
import { useStore } from "./store";
import {
  BrowserRouter,
  Routes,
  Route,
  Link,
  useNavigate,
  Navigate,
} from "react-router-dom";
import styles from "./styles";
import CreateAccountPage from "./CreateAccountPage";
import Header from "./Header";

function Welcome() {
  return <Header />;
}

function GamePage() {
  const user = useStore((state) => state.user);
  const logout = useStore((state) => state.logout);
  const tap = useStore((state) => state.tap);
  const { pizzas, coins } = useStore((state) => state.gameState.resources);

  if (user === null) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div>
      <button className={classnames(styles.button)} onClick={() => logout()}>
        Logout
      </button>
      <button className={classnames(styles.button)} onClick={() => tap()}>
        Tap
      </button>
      <div>
        <span>Coins: {coins.toString()}</span>
        <span> | </span>
        <span>Pizzas: {pizzas.toString()}</span>
      </div>
    </div>
  );
}

function CreateAccountPromotion() {
  const navigate = useNavigate();
  const onClick = (e: React.MouseEvent) => {
    e.preventDefault();
    navigate("/create-account");
  };

  return (
    <div
      className={classnames(
        "mt-10",
        "flex",
        "justify-center",
        "flex-col",
        "items-center"
      )}
    >
      No account yet?
      <button className={classnames(styles.button)} onClick={onClick}>
        Create Account
      </button>
    </div>
  );
}

function LoginPage() {
  const navigate = useNavigate();
  const start = useStore((state) => state.start);
  const onLogin = () => {
    start("");
    navigate("/");
  };
  return (
    <div>
      <Welcome />
      <LoginForm onLogin={onLogin} />
      <CreateAccountPromotion />
    </div>
  );
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<GamePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/create-account" element={<CreateAccountPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
