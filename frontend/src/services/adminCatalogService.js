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
  listModels: async () => {
    return await apiClient.get('/admin/catalog/models');
  },
  createModel: async (payload) => {
    return await apiClient.post('/admin/catalog/models', payload);
  },
  updateModel: async (id, payload) => {
    return await apiClient.patch(`/admin/catalog/models/${id}`, payload);
  },
  deleteModel: async (id) => {
    return await apiClient.delete(`/admin/catalog/models/${id}`);
  },
  uploadModelImage: async (id, file) => {
    const formData = new FormData();
    formData.append('file', file);
    return await apiClient.post(`/admin/catalog/models/${id}/image`, formData);
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
