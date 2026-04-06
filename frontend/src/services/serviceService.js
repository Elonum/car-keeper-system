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

  /** Все записи на ТО (staff с правом appointments.view_any). */
  getStaffAppointments: async () => {
    return await apiClient.get('/admin/appointments');
  },

  getAppointment: async (appointmentId) => {
    return await apiClient.get(`/service/appointments/${appointmentId}`);
  },

  /** @param {{ date: string, service_type_ids: string[] }} params */
  getBranchAvailability: async (branchId, params) => {
    const { date, service_type_ids } = params;
    const ids = Array.isArray(service_type_ids) ? service_type_ids.join(',') : String(service_type_ids);
    return await apiClient.get(`/service/branches/${branchId}/availability`, {
      params: { date, service_type_ids: ids },
    });
  },

  cancelAppointment: async (appointmentId) => {
    return await apiClient.patch(`/service/appointments/${appointmentId}/cancel`);
  },

  rescheduleAppointment: async (appointmentId, { appointment_date }) => {
    return await apiClient.patch(`/service/appointments/${appointmentId}/reschedule`, {
      appointment_date,
    });
  },

  getServiceTypes: async (params = {}) => {
    return await apiClient.get('/service/types', { params });
  },
};

