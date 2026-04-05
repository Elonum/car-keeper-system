import apiClient from '@/api/client';

export const configuratorService = {
  getColors: async () => {
    return await apiClient.get('/configurator/colors');
  },

  getOptions: async (trimId) => {
    return await apiClient.get('/configurator/options', { params: { trim_id: trimId } });
  },

  createConfiguration: async (configData) => {
    return await apiClient.post('/configurator/configurations', configData);
  },

  updateConfiguration: async (configId, configData) => {
    const payload = { ...configData };
    if (payload.status === '' || payload.status === undefined || payload.status === null) {
      delete payload.status;
    }
    return await apiClient.put(`/configurator/configurations/${configId}`, payload);
  },

  getConfiguration: async (configId) => {
    return await apiClient.get(`/configurator/configurations/${configId}`);
  },

  deleteConfiguration: async (configId) => {
    return await apiClient.delete(`/configurator/configurations/${configId}`);
  },
};
