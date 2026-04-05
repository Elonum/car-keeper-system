import React, { createContext, useState, useContext, useEffect } from 'react';
import { authService } from '@/services/authService';
import { getApiErrorMessage } from '@/lib/apiErrors';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoadingAuth, setIsLoadingAuth] = useState(true);
  const [authError, setAuthError] = useState(null);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      setIsLoadingAuth(true);
      setAuthError(null);

      if (!authService.isAuthenticated()) {
        setIsAuthenticated(false);
        setUser(null);
        setIsLoadingAuth(false);
        return;
      }

      try {
        const currentUser = await authService.getCurrentUser();
        setUser(currentUser);
        setIsAuthenticated(true);
      } catch (error) {
        if (error.status === 401) {
          authService.logout();
        }
        setIsAuthenticated(false);
        setUser(null);
      }
    } catch (error) {
      setAuthError({
        type: 'unknown',
        message: getApiErrorMessage(error, 'Не удалось проверить авторизацию'),
      });
      setIsAuthenticated(false);
      setUser(null);
    } finally {
      setIsLoadingAuth(false);
    }
  };

  const login = async (email, password) => {
    try {
      const response = await authService.login(email, password);
      if (response && response.user) {
        setUser(response.user);
        setIsAuthenticated(true);
      }
      setAuthError(null);
      return response;
    } catch (error) {
      setAuthError({
        type: 'login_failed',
        message: getApiErrorMessage(error, 'Вход не выполнен'),
      });
      throw error;
    }
  };

  /** Creates account only; no JWT — client must call login() after redirect to Login. */
  const register = async (userData) => {
    try {
      const response = await authService.register(userData);
      setAuthError(null);
      return response;
    } catch (error) {
      setAuthError({
        type: 'register_failed',
        message: getApiErrorMessage(error, 'Регистрация не выполнена'),
      });
      throw error;
    }
  };

  const logout = (shouldRedirect = true) => {
    authService.logout();
    setUser(null);
    setIsAuthenticated(false);
    setAuthError(null);
  };

  const navigateToLogin = () => {
    window.location.href = '/Login';
  };

  /** Reloads user from GET /auth/me (e.g. after profile PATCH). */
  const refreshUser = async () => {
    const currentUser = await authService.getCurrentUser();
    if (currentUser) {
      setUser(currentUser);
      setIsAuthenticated(true);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated,
        isLoadingAuth,
        authError,
        login,
        register,
        logout,
        navigateToLogin,
        checkAuth,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
