import apiClient from '@/api/client';

export const newsService = {
  getNews: async (params = {}) => {
    return await apiClient.get('/news', { params });
  },

  getNewsById: async (newsId) => {
    return await apiClient.get(`/news/${newsId}`);
  },

  createNews: async (payload) => {
    return await apiClient.post('/news', payload);
  },

  updateNews: async (newsId, payload) => {
    return await apiClient.put(`/news/${newsId}`, payload);
  },

  publishNews: async (newsId) => {
    return await apiClient.patch(`/news/${newsId}/publish`);
  },

  unpublishNews: async (newsId) => {
    return await apiClient.patch(`/news/${newsId}/unpublish`);
  },

  deleteNews: async (newsId) => {
    return await apiClient.delete(`/news/${newsId}`);
  },
};

