import React, { useContext } from "react";
import { Navigate } from "react-router-dom";
import { UserContext } from "../context/UserContext";

interface AuthGuardProps {
  children: React.ReactNode;
  requiredRole?: string; // Optional: for role-based access control
}

const AuthGuard: React.FC<AuthGuardProps> = ({ children, requiredRole }) => {
  const userContext = useContext(UserContext);

  if (!userContext) {
    console.error("AuthGuard must be used within a UserProvider");
    return <Navigate to="/login" replace />;
  }

  const { user, isLoading } = userContext;

  if (isLoading) {
    // While user data is being fetched, render nothing or a loading spinner
    return null; // Or a loading spinner component
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRole) {
    if (user.role === undefined) {
      // This case should ideally not happen after the fix in UserContext, but as a safeguard
      return <Navigate to="/dashboard" replace />; // Treat as unauthorized for role-based access
    }
    if (user.role !== requiredRole) {
      return <Navigate to="/dashboard" replace />; // Or a specific /forbidden page
    }
  }

  return <>{children}</>;
};

export default AuthGuard;
