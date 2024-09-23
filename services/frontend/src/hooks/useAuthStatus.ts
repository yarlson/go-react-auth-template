import { useQuery } from "@tanstack/react-query";
import { api } from "../api/apiClient";

export const useAuthStatus = () => {
  return useQuery({
    queryKey: ["authStatus"],
    queryFn: () =>
      api
        .url("/api/ping")
        .get()
        .json()
        .then(() => true)
        .catch(() => false),
    retry: false,
  });
};
