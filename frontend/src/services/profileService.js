import apiClient from '@/api/client';

export const profileService = {
  updateProfile: async (payload) => {
    return apiClient.patch('/profile/me', payload);
  },

  changePassword: async (payload) => {
    return apiClient.post('/profile/me/password', payload);
  },

  getUserCars: async () => {
    return await apiClient.get('/profile/cars');
  },

  createUserCar: async (carData) => {
    return await apiClient.post('/profile/cars', carData);
  },

  deleteUserCar: async (carId) => {
    return apiClient.delete(`/profile/cars/${carId}`);
  },

  getUserCar: async (carId) => {
    return await apiClient.get(`/profile/cars/${carId}`);
  },

  getConfigurations: async (params = {}) => {
    return await apiClient.get('/profile/configurations', { params });
  },
};

