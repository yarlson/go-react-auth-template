import { useLogout } from "../hooks/useLogout.ts";
import { useEffect } from "react";
import { ErrorMessage } from "../components/ErrorMessage.tsx";

export const Logout = () => {
  const { logout, isPending, isError, errorMessage } = useLogout();

  useEffect(() => {
    logout();
  }, [logout]);

  if (isPending) {
    return <div>Logging out...</div>;
  }

  if (isError) {
    return <ErrorMessage message={errorMessage} />;
  }

  return null;
};
