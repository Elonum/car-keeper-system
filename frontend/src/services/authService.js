import apiClient from '@/api/client';

export const authService = {
  login: async (email, password) => {
    const response = await apiClient.post('/auth/login', { email, password });
    if (response && response.token) {
      localStorage.setItem('auth_token', response.token);
      if (response.user) {
        localStorage.setItem('user', JSON.stringify(response.user));
      }
    }
    return response;
  },

  register: async (userData) => {
    const response = await apiClient.post('/auth/register', userData);
    if (response && response.token) {
      localStorage.setItem('auth_token', response.token);
      if (response.user) {
        localStorage.setItem('user', JSON.stringify(response.user));
      }
    }
    return response;
  },

  getCurrentUser: async () => {
    const user = await apiClient.get('/auth/me');
    if (user) {
      localStorage.setItem('user', JSON.stringify(user));
    }
    return user;
  },

  isAuthenticated: () => {
    return !!localStorage.getItem('auth_token');
  },

  logout: () => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user');
    window.location.href = '/';
  },

  getStoredUser: () => {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  },
};

