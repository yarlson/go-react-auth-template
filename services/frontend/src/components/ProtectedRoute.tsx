import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuthStatus } from "../hooks/useAuthStatus"; // Adjust the import path as needed

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const location = useLocation();
  const { data: isAuthenticated, isLoading } = useAuthStatus();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute;
