import { Routes, Route, Navigate } from "react-router-dom";
import ProtectedRoute from "./components/ProtectedRoute";
import { ThemeProvider } from "@/components/ThemeProvider";
import Login from "@/pages/Login.tsx";
import { Callback } from "@/pages/Callback.tsx";
import UserInfo from "@/pages/UserInfo.tsx";
import Layout from "@/components/Layout.tsx";
import { Logout } from "@/pages/Logout.tsx";

const App = () => {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/callback" element={<Callback />} />
        <Route
          path="/logout"
          element={
            <ProtectedRoute>
              <Logout />
            </ProtectedRoute>
          }
        />
        <Route
          path="/profile"
          element={
            <ProtectedRoute>
              <Layout>
                <UserInfo />
              </Layout>
            </ProtectedRoute>
          }
        />
        <Route path="/" element={<Navigate to="/profile" />} />
      </Routes>
    </ThemeProvider>
  );
};

export default App;
