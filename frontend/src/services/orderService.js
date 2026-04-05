import apiClient from '@/api/client';

export const orderService = {
  /** Active statuses for clients (code + customer_label_ru); matches GET /api/order-statuses */
  getPublicOrderStatuses: async () => {
    return await apiClient.get('/order-statuses');
  },

  createOrder: async (orderData) => {
    return await apiClient.post('/orders', orderData);
  },

  getOrders: async (params = {}) => {
    return await apiClient.get('/orders', { params });
  },

  getOrder: async (orderId) => {
    return await apiClient.get(`/orders/${orderId}`);
  },

  updateOrderStatus: async (orderId, status) => {
    return await apiClient.patch(`/orders/${orderId}/status`, { status });
  },
};

