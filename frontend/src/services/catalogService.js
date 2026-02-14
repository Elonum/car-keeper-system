import apiClient from '@/api/client';

export const catalogService = {
  getTrims: async (params = {}) => {
    return await apiClient.get('/catalog/trims', { params });
  },

  getTrim: async (trimId) => {
    return await apiClient.get(`/catalog/trims/${trimId}`);
  },

  getBrands: async () => {
    return await apiClient.get('/catalog/brands');
  },

  getModels: async (brandId) => {
    return await apiClient.get('/catalog/models', { params: { brand_id: brandId } });
  },

  getGenerations: async (modelId) => {
    return await apiClient.get('/catalog/generations', { params: { model_id: modelId } });
  },

  getEngineTypes: async () => {
    return await apiClient.get('/catalog/engine-types');
  },

  getTransmissions: async () => {
    return await apiClient.get('/catalog/transmissions');
  },

  getDriveTypes: async () => {
    return await apiClient.get('/catalog/drive-types');
  },
};

