import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { UserInfo, fetchUserInfo, isAuthError } from "../api/apiClient";

export const useUserInfo = (): UseQueryResult<UserInfo, Error> => {
  return useQuery<UserInfo, Error>({
    queryKey: ["userInfo"],
    queryFn: fetchUserInfo,
    retry: (failureCount, error) => !isAuthError(error) && failureCount < 3,
    staleTime: 1000 * 60 * 5,
  });
};
