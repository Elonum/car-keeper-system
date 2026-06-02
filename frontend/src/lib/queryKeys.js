/** Shared React Query keys — keep reads and invalidations aligned. */

export const queryKeys = {
  brands: () => ['brands'],
  trims: (params) => (params ? ['trims', params] : ['trims']),
  engineTypes: () => ['engineTypes'],
  transmissions: () => ['transmissions'],
  driveTypes: () => ['driveTypes'],
  serviceTypes: (filters) => ['serviceTypes', filters ?? {}],
  branches: () => ['branches'],
  branchAvailability: (branchId, dateKey, serviceKey) =>
    ['branch-availability', branchId, dateKey, serviceKey],
  branchAvailabilityReschedule: (branchId, dateKey, idKey) =>
    ['branch-availability', 'reschedule', branchId, dateKey, idKey],
  myConfigurations: () => ['my-configurations'],
  myOrders: () => ['my-orders'],
  staffOrders: () => ['staff-orders'],
  myAppointments: () => ['my-appointments'],
  staffAppointments: () => ['staff-appointments'],
  myDocuments: () => ['my-documents'],
  myCars: () => ['my-cars'],
  userCars: () => ['userCars'],
  catalogBrandsAdmin: () => ['catalog', 'brands'],
  catalogModelsAdmin: () => ['catalog', 'models', 'admin'],
  serviceTypesAdmin: () => ['service', 'types', 'admin'],
  serviceBranchesAdmin: () => ['service', 'branches'],
  news: (scope, canManage) => ['news', scope, canManage],
  newsItem: (id) => ['news', id],
};

/** Invalidate public catalog caches after admin catalog mutations. */
export function invalidatePublicCatalog(qc) {
  qc.invalidateQueries({ queryKey: queryKeys.brands() });
  qc.invalidateQueries({ queryKey: ['trims'] });
}

/** Invalidate garage and appointment car pickers. */
export function invalidateGarageRelated(qc) {
  qc.invalidateQueries({ queryKey: queryKeys.myCars() });
  qc.invalidateQueries({ queryKey: queryKeys.userCars() });
}

/** Invalidate branch availability (any date/service variant). */
export function invalidateBranchAvailability(qc) {
  qc.invalidateQueries({ queryKey: ['branch-availability'] });
}
