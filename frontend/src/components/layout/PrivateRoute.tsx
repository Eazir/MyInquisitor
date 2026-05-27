import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { Loading } from '../ui/Loading';

export function PrivateRoute() {
  const { isAuthenticated, loading } = useAuth();

  if (loading) return <Loading size="lg" text="Checking authentication..." />;
  if (!isAuthenticated) return <Navigate to="/login" replace />;

  return <Outlet />;
}
