import React from "react";
import { Link } from "react-router-dom";
import { useApiError } from "../hooks/useApiError";
import { useUserInfo } from "../hooks/useUserInfo";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ErrorMessage } from "@/components/ErrorMessage";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";

const UserInfo: React.FC = () => {
  const { data, isPending, error } = useUserInfo();
  const errorMessage = useApiError(error);

  if (isPending) return <UserInfoSkeleton />;
  if (error) return <ErrorMessage message={errorMessage} />;
  if (!data)
    return (
      <div className="text-center text-gray-500">No user data available</div>
    );

  return (
    <Card className="shadow-lg max-w-md mx-auto">
      <CardHeader className="text-center">
        <Avatar className="w-24 h-24 mx-auto mb-4">
          <AvatarImage
            src={data.pictureUrl}
            alt={`${data.firstName} ${data.lastName}`}
          />
          <AvatarFallback>
            {getInitials(data.firstName, data.lastName)}
          </AvatarFallback>
        </Avatar>
        <CardTitle className="text-2xl font-bold">{`${data.firstName} ${data.lastName}`}</CardTitle>
      </CardHeader>
      <CardContent>
        <p>
          <span className="font-semibold">Email:</span> {data.email}
        </p>
      </CardContent>
      <CardFooter className="justify-center">
        <Button asChild>
          <Link to="/logout">Logout</Link>
        </Button>
      </CardFooter>
    </Card>
  );
};

const UserInfoSkeleton: React.FC = () => (
  <Card className="max-w-md mx-auto shadow-lg">
    <CardHeader className="text-center">
      <Skeleton className="w-24 h-24 rounded-full mx-auto mb-4" />
      <Skeleton className="h-8 w-3/4 mx-auto" />
    </CardHeader>
    <CardContent className="space-y-4">
      <Skeleton className="h-4 w-full" />
    </CardContent>
    <CardFooter className="justify-center">
      <Skeleton className="h-10 w-24" />
    </CardFooter>
  </Card>
);

const getInitials = (firstName: string, lastName: string): string =>
  `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase();

export default UserInfo;
