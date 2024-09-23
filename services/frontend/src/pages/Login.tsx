import React from "react";
import { Navigate } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Layout from "@/components/Layout";
import { useAuthStatus } from "@/hooks/useAuthStatus";

const Login: React.FC = () => {
  const { data: isAuthenticated, isLoading } = useAuthStatus();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (isAuthenticated) {
    return <Navigate to="/profile" replace />;
  }

  return (
    <Layout>
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>
            <h1 className="text-2xl font-bold">Login</h1>
          </CardTitle>
        </CardHeader>
        <CardFooter>
          <Button asChild variant="default">
            <a href={`${import.meta.env.VITE_API_BASE_URL}/auth/google`}>
              Login with Google
            </a>
          </Button>
        </CardFooter>
      </Card>
    </Layout>
  );
};

export default Login;
