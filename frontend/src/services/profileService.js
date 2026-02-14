import apiClient from '@/api/client';

export const profileService = {
  getUserCars: async () => {
    return await apiClient.get('/profile/cars');
  },

  createUserCar: async (carData) => {
    return await apiClient.post('/profile/cars', carData);
  },

  getUserCar: async (carId) => {
    return await apiClient.get(`/profile/cars/${carId}`);
  },

  getConfigurations: async (params = {}) => {
    return await apiClient.get('/profile/configurations', { params });
  },
};

