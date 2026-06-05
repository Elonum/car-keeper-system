import apiClient from '@/api/client';

export const authService = {
  login: async (email, password) => {
    const response = await apiClient.post('/auth/login', { email, password });
    if (response?.user) {
      sessionStorage.setItem('user', JSON.stringify(response.user));
    }
    return response;
  },

  /** POST /auth/register returns user profile only; session is established via login(). */
  register: async (userData) => {
    const response = await apiClient.post('/auth/register', userData);
    return response;
  },

  getCurrentUser: async () => {
    const user = await apiClient.get('/auth/me', { skipAuthRedirect: true });
    if (user) {
      sessionStorage.setItem('user', JSON.stringify(user));
    }
    return user;
  },

  /** Clears client session only (no redirect, no API). Use when /auth/me returns 401. */
  clearSession: () => {
    sessionStorage.removeItem('user');
  },

  logout: async () => {
    try {
      await apiClient.post('/auth/logout', null, { skipAuthRedirect: true });
    } catch {
      // clear local session even if network fails
    }
    authService.clearSession();
    window.location.href = '/';
  },
};
