import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useStore } from "./store";

function LogoutPage() {
  const logout = useStore((state) => state.logout);
  const navigate = useNavigate();

  useEffect(() => {
    logout();
    navigate("/login");
  }, []);

  return null;
}

export default LogoutPage;
