import React, { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

interface PrivateRouteProps {
  children: ReactNode;
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
  const { user } = useAuth();

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  // Temporarily allow access without active subscription for testing
  // TODO: Re-enable subscription check once payment system is fully configured
  // if (user.subscription_status !== 'active') {
  //   return <Navigate to="/subscription" replace />;
  // }

  return <>{children}</>;
};

export default PrivateRoute;