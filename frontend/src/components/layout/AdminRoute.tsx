import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { Loading } from '../ui/Loading';

export function AdminRoute() {
  const { isAdmin, loading } = useAuth();

  if (loading) return <Loading size="lg" text="Checking permissions..." />;
  if (!isAdmin) return <Navigate to="/dashboard" replace />;

  return <Outlet />;
}
