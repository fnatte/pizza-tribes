import React from "react";
import {useNavigate} from "react-router-dom";
import {classnames} from "tailwindcss-classnames";
import Header from "./Header";
import LoginForm from "./LoginForm";
import {useStore} from "./store";
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
  const onLogin = () => {
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

export default LoginPage;
