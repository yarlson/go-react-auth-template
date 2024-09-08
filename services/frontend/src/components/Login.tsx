import { Button } from "@/components/ui/button";

const Login = () => {
  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="p-4 text-center">
        <h1 className="text-2xl font-bold mb-4">Login</h1>
        <Button asChild variant="default">
          <a href={`${import.meta.env.VITE_API_BASE_URL}/auth/google`}>
            Login with Google
          </a>
        </Button>
      </div>
    </div>
  );
};

export default Login;
