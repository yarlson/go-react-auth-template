import { Routes, Route, Navigate } from "react-router-dom";
import Login from "./components/Login";
import ProtectedRoute from "./components/ProtectedRoute";
import UserInfo from "./components/UserInfo";
import { Callback, Logout } from "./components/Callback";

const App = () => {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/callback" element={<Callback />} />
      <Route
        path="/user-info"
        element={
          <ProtectedRoute>
            <UserInfo />
          </ProtectedRoute>
        }
      />
      <Route
        path="/logout"
        element={
          <ProtectedRoute>
            <Logout />
          </ProtectedRoute>
        }
      />
      <Route path="/" element={<Navigate to="/user-info" replace />} />
    </Routes>
  );
};

export default App;