const Login = () => {
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold">Login</h1>
      <a href={`${import.meta.env.VITE_API_BASE_URL}/auth/google`}>
        Login with Google
      </a>
    </div>
  );
};

export default Login;
