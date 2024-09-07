import React, { useEffect } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  api,
  isAuthError,
  isOfflineResponse,
  isServerErrorResponse,
} from "../api/apiClient";
import { Link, useNavigate } from "react-router-dom";

interface UserInfoResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
}

const fetchUserInfo = async (): Promise<UserInfoResponse> => {
  const response = await api.url("/api/user/profile").get().res();

  if (isOfflineResponse(response)) {
    const data = await response.json();
    throw new Error(`Offline: ${data.debug}`);
  }

  if (isServerErrorResponse(response)) {
    const data = await response.json();
    throw new Error(`Server Error: ${data.status}`);
  }

  return await response.json();
};

const UserInfo: React.FC = () => {
  const navigate = useNavigate();
  const { data, isPending, error } = useQuery({
    queryKey: ["userInfo"],
    queryFn: fetchUserInfo,
    retry: (failureCount, error) => !isAuthError(error) && failureCount < 3,
    staleTime: 1000 * 60 * 5,
  });

  useEffect(() => {
    if (isAuthError(error)) {
      navigate("/login");
    }
  }, [error, navigate]);

  if (isPending) {
    return <LoadingMessage />;
  }
  if (error instanceof Error) {
    return <ErrorMessage error={error} />;
  }

  if (!data) {
    return <NoDataMessage />;
  }

  return <UserInfoDisplay data={data} />;
};

const LoadingMessage: React.FC = () => <div>Loading user info...</div>;

const ErrorMessage: React.FC<{ error: Error }> = ({ error }) => {
  if (error.message.startsWith("Offline:")) {
    return (
      <div>You are currently offline. Some features may be unavailable.</div>
    );
  }
  if (error.message.startsWith("Server Error:")) {
    return <div>A server error occurred. Please try again later.</div>;
  }
  return <div>An error occurred. Please try again later.</div>;
};

const NoDataMessage: React.FC = () => <div>No user data available</div>;

const UserInfoDisplay: React.FC<{ data: UserInfoResponse }> = ({ data }) => (
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

export default UserInfo;
