import apiClient from '@/api/client';

export const orderService = {
  /** Active statuses for clients (code + customer_label_ru); matches GET /api/order-statuses */
  getPublicOrderStatuses: async () => {
    return await apiClient.get('/order-statuses');
  },

  getAdminOrderStatuses: async () => {
    return await apiClient.get('/admin/order-statuses');
  },

  createOrder: async (orderData) => {
    return await apiClient.post('/orders', orderData);
  },

  getOrders: async (params = {}) => {
    return await apiClient.get('/orders', { params });
  },

  /** Все заказы (staff с правом orders.view_any). */
  getStaffOrders: async () => {
    return await apiClient.get('/admin/orders');
  },

  getOrder: async (orderId) => {
    return await apiClient.get(`/orders/${orderId}`);
  },

  updateOrderStatus: async (orderId, status) => {
    return await apiClient.patch(`/orders/${orderId}/status`, { status });
  },

  createOrderStatus: async (payload) => {
    return await apiClient.post('/admin/order-statuses', payload);
  },

  patchOrderStatus: async (id, payload) => {
    return await apiClient.patch(`/admin/order-statuses/${id}`, payload);
  },

  deleteOrderStatus: async (id) => {
    return await apiClient.delete(`/admin/order-statuses/${id}`);
  },
};

