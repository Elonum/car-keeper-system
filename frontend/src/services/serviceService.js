import apiClient from '@/api/client';

export const serviceService = {
  getUserCars: async () => {
    return await apiClient.get('/service/user-cars');
  },

  getBranches: async (params = {}) => {
    return await apiClient.get('/service/branches', { params });
  },

  createAppointment: async (appointmentData) => {
    return await apiClient.post('/service/appointments', appointmentData);
  },

  getAppointments: async (params = {}) => {
    return await apiClient.get('/service/appointments', { params });
  },

  getAppointment: async (appointmentId) => {
    return await apiClient.get(`/service/appointments/${appointmentId}`);
  },

  updateAppointment: async (appointmentId, appointmentData) => {
    return await apiClient.put(`/service/appointments/${appointmentId}`, appointmentData);
  },

  cancelAppointment: async (appointmentId) => {
    return await apiClient.patch(`/service/appointments/${appointmentId}/cancel`);
  },

  getServiceTypes: async (params = {}) => {
    return await apiClient.get('/service/types', { params });
  },
};

