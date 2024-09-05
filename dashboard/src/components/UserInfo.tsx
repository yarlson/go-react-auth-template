import React from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "../utils/api";
import { Link } from "react-router-dom";

interface UserInfoResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
}

const fetchUserInfo = async (): Promise<UserInfoResponse> => {
  try {
    const response = await api.url("/api/user/profile").get().json();
    return response as UserInfoResponse;
  } catch (error) {
    if (
      error instanceof Error &&
      "status" in error &&
      (error as any).status === 401
    ) {
      // If unauthorized, try to refresh the session
      await api.url("/auth/refresh").post().json();
      // Retry the original request
      const retryResponse = await api.url("/api/user/profile").get().json();
      return retryResponse as UserInfoResponse;
    }
    throw error;
  }
};

const UserInfo: React.FC = () => {
  const { data, isPending, error } = useQuery({
    queryKey: ["userInfo"],
    queryFn: fetchUserInfo,
    retry: false, // We're handling retry logic in fetchUserInfo
  });

  if (isPending) return <div>Loading user info...</div>;
  if (error) return <div>Error: {(error as Error).message}</div>;
  if (!data) return <div>No user data available</div>;

  return (
    <div>
      <h1>User Info</h1>
      <p>Email: {data.email}</p>
      <p>First Name: {data.firstName}</p>
      <p>Last Name: {data.lastName}</p>
      <p>
        <Link to="/logout">Logout</Link>
      </p>
    </div>
  );
};

export default UserInfo;
