import React, { useEffect } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import classnames from "classnames";
import Header from "./Header";
import LoginForm from "./LoginForm";
import { useStore } from "./store";
import styles from "./styles";

function Welcome() {
  return <Header />;
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
        "mt-5",
        "flex",
        "justify-center",
        "flex-col",
        "items-center"
      )}
    >
      No account yet?
      <button className={classnames(styles.primaryButton)} onClick={onClick}>
        Create Account
      </button>
    </div>
  );
}

function LoginPage() {
  const navigate = useNavigate();
  const onLogin = () => {
    navigate("/games");
  };
  const connectionState = useStore((state) => state.connectionState);
  const connection = useStore((state) => state.connection);
  const user = useStore((state) => state.user);

  const shouldRedirectToLoggedInState =
    connectionState?.connected || connectionState?.connecting || user !== null;

  // Make sure we reset any existing connection state so that
  // no previous connection errors affects the new login
  useEffect(() => {
    if (!shouldRedirectToLoggedInState && connection) {
      connection.reset();
    }
  }, [shouldRedirectToLoggedInState, connection]);

  if (shouldRedirectToLoggedInState) {
    return <Navigate to="/" replace />;
  }

  return (
    <div
      className={classnames(
        "flex",
        "flex-col",
        "justify-center",
        "items-center"
      )}
    >
      <Welcome />
      <LoginForm onLogin={onLogin} />
      <div className={classnames("text-2xl", "mt-5")}>🍕🍕🍕🍕</div>
      <CreateAccountPromotion />
    </div>
  );
}

export default LoginPage;
