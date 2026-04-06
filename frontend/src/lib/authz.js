export const PERMISSIONS = {
  ORDERS_VIEW_ANY: 'orders.view_any',
  ORDERS_MANAGE_STATUS: 'orders.manage_status',
  CONFIGURATIONS_VIEW_ANY: 'configurations.view_any',
  CONFIGURATIONS_MANAGE: 'configurations.manage',
  APPOINTMENTS_VIEW_ANY: 'appointments.view_any',
  GARAGE_VIEW_ANY: 'garage.view_any',
  DOCUMENTS_VIEW_ANY: 'documents.view_any',
  NEWS_MANAGE: 'news.manage',
  ADMIN_ORDER_STATUSES: 'admin.order_statuses',
  ADMIN_ROLES_VIEW: 'admin.roles_view',
};

const ALL_PERMISSIONS = Object.values(PERMISSIONS);

const DEFAULT_ROLE_PERMISSIONS = {
  customer: [],
  manager: [
    PERMISSIONS.ORDERS_VIEW_ANY,
    PERMISSIONS.ORDERS_MANAGE_STATUS,
    PERMISSIONS.CONFIGURATIONS_VIEW_ANY,
    PERMISSIONS.CONFIGURATIONS_MANAGE,
    PERMISSIONS.APPOINTMENTS_VIEW_ANY,
    PERMISSIONS.GARAGE_VIEW_ANY,
    PERMISSIONS.DOCUMENTS_VIEW_ANY,
    PERMISSIONS.NEWS_MANAGE,
  ],
  service_advisor: [
    PERMISSIONS.ORDERS_VIEW_ANY,
    PERMISSIONS.ORDERS_MANAGE_STATUS,
    PERMISSIONS.CONFIGURATIONS_VIEW_ANY,
    PERMISSIONS.APPOINTMENTS_VIEW_ANY,
    PERMISSIONS.GARAGE_VIEW_ANY,
    PERMISSIONS.DOCUMENTS_VIEW_ANY,
  ],
  admin: ALL_PERMISSIONS,
};

export function hasPermission(role, permission) {
  if (!role || !permission) return false;
  const perms = DEFAULT_ROLE_PERMISSIONS[role] || [];
  return perms.includes(permission);
}

