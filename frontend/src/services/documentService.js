import apiClient, { API_BASE_URL, getApiAuthHeaders } from '@/api/client';
import { authService } from '@/services/authService';
import { formatBackendErrorMessage } from '@/lib/apiErrors';

export const getDocumentsApiBaseUrl = () => API_BASE_URL;

export const DOCUMENT_TYPE_LABELS = {
  commercial_offer: 'Коммерческое предложение',
  order_contract: 'Договор',
  service_order: 'Заказ-наряд',
  service_act: 'Акт',
};

export const documentService = {
  list: async (params = {}) => {
    return await apiClient.get('/documents', { params });
  },

  get: async (documentId) => {
    return await apiClient.get(`/documents/${documentId}`);
  },

  upload: async ({ file, documentType, orderId, serviceAppointmentId }) => {
    const form = new FormData();
    form.append('file', file);
    form.append('document_type', documentType);
    if (orderId) form.append('order_id', orderId);
    if (serviceAppointmentId) form.append('service_appointment_id', serviceAppointmentId);
    return await apiClient.post('/documents', form);
  },

  delete: async (documentId) => {
    return await apiClient.delete(`/documents/${documentId}`);
  },

  /** Download file via fetch (binary, not JSON envelope). */
  download: async (documentId, fallbackName = 'document') => {
    const url = `${getDocumentsApiBaseUrl()}/documents/${documentId}/file`;
    const res = await fetch(url, {
      credentials: 'include',
      headers: getApiAuthHeaders(),
    });
    if (!res.ok) {
      if (res.status === 401) {
        authService.clearSession();
        const onAuthPage =
          window.location.pathname === '/Login' ||
          window.location.pathname === '/Register';
        if (!onAuthPage) {
          window.location.href = '/Login';
        }
      }
      let raw = '';
      try {
        const data = await res.json();
        raw = data?.error || data?.message || '';
      } catch {
        /* not JSON */
      }
      const message =
        formatBackendErrorMessage(raw) ||
        raw ||
        (res.status === 401
          ? 'Требуется вход'
          : res.status === 404
            ? 'Файл недоступен на сервере'
            : 'Не удалось скачать файл');
      const err = new Error(message);
      err.status = res.status;
      throw err;
    }
    const blob = await res.blob();
    const cd = res.headers.get('Content-Disposition');
    let name = fallbackName;
    if (cd) {
      const m = /filename\*?=(?:UTF-8'')?["']?([^"';]+)/i.exec(cd);
      if (m) name = decodeURIComponent(m[1].replace(/["']/g, ''));
    }
    const href = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = href;
    a.download = name;
    a.rel = 'noopener';
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(href);
  },
};
