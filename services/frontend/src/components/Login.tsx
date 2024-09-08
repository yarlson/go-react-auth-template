import { Card, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Layout from "@/components/Layout";

const Login = () => {
  return (
    <Layout>
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
    </Layout>
  );
};

export default Login;
