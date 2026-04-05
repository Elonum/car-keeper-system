import apiClient from '@/api/client';

export const getDocumentsApiBaseUrl = () =>
  import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

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
    const token = localStorage.getItem('auth_token');
    const url = `${getDocumentsApiBaseUrl()}/documents/${documentId}/file`;
    const res = await fetch(url, {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    if (!res.ok) {
      let msg = `Download failed (${res.status})`;
      try {
        const data = await res.json();
        if (data?.error) msg = data.error;
      } catch {
        /* ignore */
      }
      throw new Error(msg);
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
    a.click();
    URL.revokeObjectURL(href);
  },
};

export const DOCUMENT_TYPE_LABELS = {
  commercial_offer: 'Коммерческое предложение',
  order_contract: 'Договор',
  service_order: 'Заказ-наряд',
  service_act: 'Акт',
};
