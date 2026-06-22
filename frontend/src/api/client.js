/// <reference types="vite/client" />

import axios from 'axios';
import { formatBackendErrorMessage } from '@/lib/apiErrors';

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

function resolveApiOrigin() {
  const base = API_BASE_URL;
  if (String(base).startsWith('http')) {
    return String(base).replace(/\/api\/?$/, '');
  }
  if (typeof window !== 'undefined' && window.location?.origin) {
    return window.location.origin;
  }
  return import.meta.env.VITE_API_ORIGIN || 'http://localhost:8080';
}

export const API_BASE_ORIGIN = resolveApiOrigin();

/** Headers for fetch (binary download) — auth via HttpOnly session cookie. */
export function getApiAuthHeaders() {
  return {};
}

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use(
  (config) => {
    if (config.data instanceof FormData) {
      delete config.headers['Content-Type'];
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
      // Go nil slices encode as JSON null; normalize for list consumers.
      return backendResponse.data === null ? [] : backendResponse.data;
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
        const skipRedirect = Boolean(error.config?.skipAuthRedirect);
        const requestUrl = String(error.config?.url || '');
        const isAuthMe = requestUrl.includes('/auth/me');

        sessionStorage.removeItem('user');

        const onAuthPage =
          window.location.pathname === '/Login' ||
          window.location.pathname === '/Register';

        if (!skipRedirect && !isAuthMe && !onAuthPage) {
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
