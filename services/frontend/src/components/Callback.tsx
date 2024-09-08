// Callback.tsx
import { useEffect, useRef } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { api } from "../api/apiClient";
import { useApiError } from "../hooks/useApiError";
import { ErrorMessage } from "../components/ErrorMessage";

export const Callback = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const hasMutated = useRef(false);

  const mutation = useMutation<Response, unknown, string>({
    mutationFn: (callbackUrl: string) =>
      api
        .url(`/auth/google/callback${callbackUrl}&provider=google`)
        .get()
        .json(),
    onSuccess: () => navigate("/user-info"),
  });

  const errorMessage = useApiError(mutation.error);

  useEffect(() => {
    if (!hasMutated.current && location.search && !mutation.isSuccess) {
      hasMutated.current = true;
      mutation.mutate(location.search);
    }
  }, [location.search, mutation]);

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

export const Logout = () => {
  const navigate = useNavigate();

  const mutation = useMutation<Response, unknown>({
    mutationFn: () => api.url("/auth/logout").get().json(),
    onSuccess: () => navigate("/"),
  });

  const errorMessage = useApiError(mutation.error);

  useEffect(() => {
    mutation.mutate();
  }, []);

  if (mutation.isPending) {
    return <div>Logging out...</div>;
  }

  if (mutation.isError) {
    return <ErrorMessage message={errorMessage} />;
  }

  return null;
};
