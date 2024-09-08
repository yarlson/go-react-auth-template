import { Card, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

const Login = () => {
  return (
    <div className="flex items-center justify-center min-h-screen">
      <Card className="w-[350px]">
        <CardHeader>
          <CardTitle>
            <h1 className="text-2xl font-bold">Login</h1>
          </CardTitle>
        </CardHeader>
        <CardFooter className="flex justify-center">
          <Button asChild variant="default">
            <a href={`${import.meta.env.VITE_API_BASE_URL}/auth/google`}>
              Login with Google
            </a>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
};

export default Login;
