import { useEffect, useRef } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { api } from "../utils/api";

export const Callback = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const hasMutated = useRef(false);

  const mutation = useMutation<void, Error, string>({
    mutationFn: (callbackUrl: string) =>
      api
        .url(`/auth/google/callback${callbackUrl}&provider=google`)
        .get()
        .json(),
    onSuccess: () => {
      navigate("/user-info");
    },
    onError: (error) => {
      console.error("Callback error:", error);
    },
  });

  useEffect(() => {
    if (!hasMutated.current && location.search && !mutation.isSuccess) {
      hasMutated.current = true;
      mutation.mutate(location.search);
    }
  }, [location.search, mutation]);

  return (
    <div>
      <h1>Login</h1>
      {mutation.isPending && <div>Processing login...</div>}
      {mutation.isError && <div>Error: {mutation.error.message}</div>}
      {mutation.isSuccess && <div>Login successful, redirecting...</div>}
    </div>
  );
};

export const Logout = () => {
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationFn: () => api.url("/auth/logout").get().json(),
    onSuccess: () => {
      navigate("/");
    },
    onError: (error) => {
      console.error("Logout error:", error);
    },
  });

  useEffect(() => {
    mutation.mutate();
  }, []);

  return null;
};
