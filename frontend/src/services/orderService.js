import apiClient from '@/api/client';

export const orderService = {
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

