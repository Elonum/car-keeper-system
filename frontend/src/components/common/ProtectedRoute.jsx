import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/lib/AuthContext';
import PageLoader from '@/components/common/PageLoader';

/**
 * Redirects unauthenticated users when `when` is true (default).
 * Uses useEffect to avoid navigation during render.
 */
export default function ProtectedRoute({ children, when = true }) {
  const { isAuthenticated, isLoadingAuth } = useAuth();
  const navigate = useNavigate();
  const shouldGuard = Boolean(when);

  useEffect(() => {
    if (!isLoadingAuth && shouldGuard && !isAuthenticated) {
      navigate('/Login', { replace: true });
    }
  }, [isLoadingAuth, shouldGuard, isAuthenticated, navigate]);

  if (isLoadingAuth || (shouldGuard && !isAuthenticated)) {
    return <PageLoader />;
  }

  return children;
}
