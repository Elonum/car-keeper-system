import apiClient from '@/api/client';

export const adminCatalogService = {
  createBrand: async (payload) => {
    return await apiClient.post('/admin/catalog/brands', payload);
  },
  updateBrand: async (id, payload) => {
    return await apiClient.patch(`/admin/catalog/brands/${id}`, payload);
  },
  deleteBrand: async (id) => {
    return await apiClient.delete(`/admin/catalog/brands/${id}`);
  },

  createServiceType: async (payload) => {
    return await apiClient.post('/admin/service/types', payload);
  },
  updateServiceType: async (id, payload) => {
    return await apiClient.patch(`/admin/service/types/${id}`, payload);
  },
  deleteServiceType: async (id) => {
    return await apiClient.delete(`/admin/service/types/${id}`);
  },

  updateBranch: async (id, payload) => {
    return await apiClient.patch(`/admin/branches/${id}`, payload);
  },
};
