import axios from 'axios';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => {
    const backendResponse = response.data;
    if (backendResponse.success) {
      return backendResponse.data;
    } else {
      return Promise.reject({
        status: response.status,
        message: backendResponse.error || backendResponse.message || 'Request failed',
        data: backendResponse,
      });
    }
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response;
      
      if (status === 401) {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('user');
        if (window.location.pathname !== '/Login') {
          window.location.href = '/Login';
        }
      }
      
      const backendResponse = data || {};
      // Extract error message from various possible locations
      const errorMessage = backendResponse.error || 
                          backendResponse.message || 
                          (typeof backendResponse === 'string' ? backendResponse : null) ||
                          error.message || 
                          `Request failed with status ${status}`;
      return Promise.reject({
        status,
        message: errorMessage,
        data: backendResponse,
      });
    } else if (error.request) {
      return Promise.reject({
        status: 0,
        message: 'Network error. Please check your connection.',
      });
    } else {
      return Promise.reject({
        status: 0,
        message: error.message || 'An unexpected error occurred',
      });
    }
  }
);

export default apiClient;
