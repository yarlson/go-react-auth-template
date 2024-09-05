import React from "react";
import { Route, Redirect, RouteProps } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { api } from "../utils/api";

interface ProtectedRouteProps extends RouteProps {
  component: React.ComponentType<any>;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  component: Component,
  ...rest
}) => {
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

  return (
    <Route
      {...rest}
      render={(props) =>
        isAuthenticated ? <Component {...props} /> : <Redirect to="/login" />
      }
    />
  );
};

export default ProtectedRoute;
