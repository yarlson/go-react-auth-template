import { useEffect, useRef } from "react";
import { useLocation } from "react-router-dom";
import { ErrorMessage } from "../components/ErrorMessage";
import { useGoogleAuthCallback } from "../hooks/useGoogleAuthCallback";

export const Callback = () => {
  const location = useLocation();
  const hasMutated = useRef(false);
  const { mutation, errorMessage, processCallback } = useGoogleAuthCallback();

  useEffect(() => {
    if (!hasMutated.current && location.search && !mutation.isSuccess) {
      hasMutated.current = true;
      processCallback(location.search);
    }
  }, [location.search, mutation, processCallback]);

  if (mutation.isPending) {
    return <div>Processing login...</div>;
  }

  if (mutation.isError) {
    return <ErrorMessage message={errorMessage} />;
  }

  if (mutation.isSuccess) {
    return <div>Login successful, redirecting...</div>;
  }

  return null;
};
