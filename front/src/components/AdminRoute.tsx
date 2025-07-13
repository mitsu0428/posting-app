import React, { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

interface AdminRouteProps {
  children: ReactNode;
}

const AdminRoute: React.FC<AdminRouteProps> = ({ children }) => {
  const { user } = useAuth();

  if (!user) {
    return <Navigate to="/admin-login-page" replace />;
  }

  if (!user.is_admin) {
    return <Navigate to="/admin-login-page" replace />;
  }

  return <>{children}</>;
};

export default AdminRoute;