import React from "react";
import { useQuery } from "@tanstack/react-query";
import { api, isAuthError } from "../utils/api";
import { Link, useNavigate } from "react-router-dom";

interface UserInfoResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
}

const fetchUserInfo = async (): Promise<UserInfoResponse> => {
  return await api.url("/api/user/profile").get().json<UserInfoResponse>();
};

const UserInfo: React.FC = () => {
  const navigate = useNavigate();
  const { data, isPending, error } = useQuery({
    queryKey: ["userInfo"],
    queryFn: fetchUserInfo,
    retry: (failureCount, error) => !isAuthError(error) && failureCount < 3,
    staleTime: 1000 * 60 * 5,
  });

  React.useEffect(() => {
    if (isAuthError(error)) {
      navigate("/login");
    }
  }, [error, navigate]);

  if (isPending) return <div>Loading user info...</div>;
  if (error) return <div>An error occurred. Please try again later.</div>;
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
