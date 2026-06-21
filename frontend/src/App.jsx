import { QueryClientProvider } from '@tanstack/react-query';
import { queryClientInstance } from '@/lib/query-client';
import { pagesConfig } from './pages.config';
import { BrowserRouter as Router, Route, Routes, useSearchParams } from 'react-router-dom';
import PageNotFound from './lib/PageNotFound';
import ProtectedRoute from '@/components/common/ProtectedRoute';
import { AuthProvider, useAuth } from '@/lib/AuthContext';

const PROTECTED_PAGES = new Set(['Profile', 'ServiceAppointment']);
const CONFIGURATOR_PAGE = 'Configurator';

const { Pages, Layout, mainPage } = pagesConfig;
const mainPageKey = mainPage ?? Object.keys(Pages)[0];
const MainPage = mainPageKey ? Pages[mainPageKey] : <></>;

const LayoutWrapper = ({ children, currentPageName }) => {
  if (currentPageName === 'Login' || currentPageName === 'Register') {
    return <>{children}</>;
  }
  return Layout ? <Layout currentPageName={currentPageName}>{children}</Layout> : <>{children}</>;
};

function ConfiguratorPage({ Page }) {
  const [searchParams] = useSearchParams();
  const requireAuth = Boolean(searchParams.get('config_id'));
  return (
    <ProtectedRoute when={requireAuth}>
      <Page />
    </ProtectedRoute>
  );
}

function renderPage(path, Page) {
  if (PROTECTED_PAGES.has(path)) {
    return (
      <ProtectedRoute>
        <Page />
      </ProtectedRoute>
    );
  }
  if (path === CONFIGURATOR_PAGE) {
    return <ConfiguratorPage Page={Page} />;
  }
  return <Page />;
}

const AuthenticatedApp = () => {
  const { isLoadingAuth } = useAuth();

  if (isLoadingAuth) {
    return (
      <div className="fixed inset-0 flex items-center justify-center">
        <div className="w-8 h-8 border-4 border-slate-200 border-t-slate-800 rounded-full animate-spin"></div>
      </div>
    );
  }

  return (
    <Routes>
      <Route path="/" element={
        <LayoutWrapper currentPageName={mainPageKey}>
          <MainPage />
        </LayoutWrapper>
      } />
      {Object.entries(Pages).map(([path, Page]) => (
        <Route
          key={path}
          path={`/${path}`}
          element={
            <LayoutWrapper currentPageName={path}>
              {renderPage(path, Page)}
            </LayoutWrapper>
          }
        />
      ))}
      <Route path="*" element={<PageNotFound />} />
    </Routes>
  );
};

function App() {
  return (
    <AuthProvider>
      <QueryClientProvider client={queryClientInstance}>
        <Router>
          <AuthenticatedApp />
        </Router>
      </QueryClientProvider>
    </AuthProvider>
  );
}

export default App;
