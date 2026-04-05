import axios from 'axios';
import { formatBackendErrorMessage } from '@/lib/apiErrors';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use(
  (config) => {
    if (config.data instanceof FormData) {
      delete config.headers['Content-Type'];
    }
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
      const raw = backendResponse.error || backendResponse.message || 'Ошибка запроса';
      return Promise.reject({
        status: response.status,
        message: formatBackendErrorMessage(raw) || raw,
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

      const backendResponse = data && typeof data === 'object' ? data : {};
      const raw = String(
        backendResponse.error ||
          backendResponse.message ||
          (typeof data === 'string' ? data : '') ||
          ''
      ).trim();
      const message = formatBackendErrorMessage(raw) || raw || 'Ошибка запроса';

      return Promise.reject({
        status,
        message,
        data: backendResponse,
      });
    }
    if (error.request) {
      return Promise.reject({
        status: 0,
        message: 'Не удалось связаться с сервером. Проверьте подключение к интернету.',
      });
    }
    return Promise.reject({
      status: 0,
      message: formatBackendErrorMessage(error.message) || 'Произошла ошибка. Попробуйте ещё раз.',
    });
  }
);

export default apiClient;
