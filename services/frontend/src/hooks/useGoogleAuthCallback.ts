import { useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { api } from "../api/apiClient";
import { useApiError } from "../hooks/useApiError";

export const useGoogleAuthCallback = () => {
  const navigate = useNavigate();

  const mutation = useMutation<Response, unknown, string>({
    mutationFn: (callbackUrl: string) =>
      api
        .url(`/auth/google/callback${callbackUrl}&provider=google`)
        .get()
        .json(),
    onSuccess: () => navigate("/profile"),
    onError: (error) => {
      if (error instanceof Error) {
        return navigate("/login");
      }
    },
  });

  const errorMessage = useApiError(mutation.error);

  return {
    mutation,
    errorMessage,
    processCallback: mutation.mutate,
  };
};
