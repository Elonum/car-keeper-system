import { useAuth } from '@/lib/AuthContext';
import PageLoader from '@/components/common/PageLoader';

export default function ProtectedRoute({ children }) {
  const { isAuthenticated, isLoadingAuth, navigateToLogin } = useAuth();

  if (isLoadingAuth) {
    return <PageLoader />;
  }

  if (!isAuthenticated) {
    navigateToLogin();
    return <PageLoader />;
  }

  return children;
}
