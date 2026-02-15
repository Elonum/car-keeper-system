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
    // Log the request payload for debugging
    console.log('[configuratorService] updateConfiguration called:', {
      configId,
      configData,
      hasStatus: 'status' in configData,
      statusValue: configData?.status,
      keys: Object.keys(configData || {}),
    });
    
    // Explicitly remove status if it's empty or undefined (for full updates)
    const payload = { ...configData };
    if (payload.status === '' || payload.status === undefined || payload.status === null) {
      delete payload.status;
      console.log('[configuratorService] Removed empty/undefined status from payload');
    }
    
    console.log('[configuratorService] Sending payload:', {
      payload,
      keys: Object.keys(payload),
    });
    
    return await apiClient.put(`/configurator/configurations/${configId}`, payload);
  },

  getConfiguration: async (configId) => {
    return await apiClient.get(`/configurator/configurations/${configId}`);
  },

  deleteConfiguration: async (configId) => {
    return await apiClient.delete(`/configurator/configurations/${configId}`);
  },
};

