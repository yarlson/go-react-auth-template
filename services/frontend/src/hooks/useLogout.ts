import { useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { api } from "../api/apiClient.ts";
import { useApiError } from "../hooks/useApiError.ts";

export const useLogout = () => {
  const navigate = useNavigate();

  const mutation = useMutation<Response, unknown>({
    mutationFn: () => api.url("/auth/logout").get().json(),
    onSuccess: () => navigate("/"),
  });

  const errorMessage = useApiError(mutation.error);

  return {
    logout: mutation.mutate,
    isPending: mutation.isPending,
    isError: mutation.isError,
    errorMessage,
  };
};
