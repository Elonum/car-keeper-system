import apiClient from '@/api/client';

export const roleService = {
  listRoles: async () => {
    return await apiClient.get('/admin/roles');
  },
};

