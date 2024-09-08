import { useQuery } from "@tanstack/react-query";
import { api, isAuthError } from "../api/apiClient";
import { Link } from "react-router-dom";
import { useApiError } from "../hooks/useApiError";
import { ErrorMessage } from "./ErrorMessage";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface UserInfoResponse {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
}

const fetchUserInfo = async (): Promise<UserInfoResponse> => {
  return await api.url("/api/user/profile").get().json();
};

const UserInfo: React.FC = () => {
  const { data, isPending, error } = useQuery({
    queryKey: ["userInfo"],
    queryFn: fetchUserInfo,
    retry: (failureCount, error) => !isAuthError(error) && failureCount < 3,
    staleTime: 1000 * 60 * 5,
  });

  const errorMessage = useApiError(error);

  if (isPending) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        Loading user info...
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <ErrorMessage message={errorMessage} />
      </div>
    );
  }

  if (!data) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        No user data available
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center min-h-screen">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>
            <h1 className="text-2xl font-bold">User Info</h1>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p>Email: {data.email}</p>
          <p>First Name: {data.firstName}</p>
          <p>Last Name: {data.lastName}</p>
        </CardContent>
        <CardFooter>
          <Button>
            <Link to="/logout">Logout</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};

export default UserInfo;
