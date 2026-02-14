import React, { createContext, useState, useContext, useEffect } from 'react';
import { authService } from '@/services/authService';

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
        message: error.message || 'Failed to check authentication',
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
        message: error.message || 'Login failed',
      });
      throw error;
    }
  };

  const register = async (userData) => {
    try {
      const response = await authService.register(userData);
      if (response && response.user) {
        setUser(response.user);
        setIsAuthenticated(true);
      }
      setAuthError(null);
      return response;
    } catch (error) {
      setAuthError({
        type: 'register_failed',
        message: error.message || 'Registration failed',
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
