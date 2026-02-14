import apiClient from '@/api/client';

export const newsService = {
  getNews: async (params = {}) => {
    return await apiClient.get('/news', { params });
  },

  getNewsById: async (newsId) => {
    return await apiClient.get(`/news/${newsId}`);
  },
};

