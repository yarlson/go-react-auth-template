import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { api } from "../api/apiClient";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const location = useLocation();

  const { data: isAuthenticated, isLoading } = useQuery({
    queryKey: ["isAuthenticated"],
    queryFn: () =>
      api
        .url("/api/user/profile")
        .get()
        .json()
        .then(() => true)
        .catch(() => false),
    retry: false,
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute;
