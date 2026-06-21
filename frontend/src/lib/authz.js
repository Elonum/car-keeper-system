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
  CATALOG_MANAGE: 'catalog.manage',
  SERVICE_MANAGE: 'service.manage',
};

const ALL_PERMISSIONS = Object.values(PERMISSIONS);

/** Fallback matrix when API permissions are unavailable (e.g. logged-out). */
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
    PERMISSIONS.SERVICE_MANAGE,
  ],
  service_advisor: [
    PERMISSIONS.ORDERS_VIEW_ANY,
    PERMISSIONS.ORDERS_MANAGE_STATUS,
    PERMISSIONS.CONFIGURATIONS_VIEW_ANY,
    PERMISSIONS.APPOINTMENTS_VIEW_ANY,
    PERMISSIONS.GARAGE_VIEW_ANY,
    PERMISSIONS.DOCUMENTS_VIEW_ANY,
    PERMISSIONS.SERVICE_MANAGE,
  ],
  admin: ALL_PERMISSIONS,
};

/** Короткие подписи для карточки «ваши возможности» */
export const PERMISSION_LABEL_RU = {
  'orders.view_any': 'Все заказы клиентов',
  'orders.manage_status': 'Смена статуса заказов',
  'configurations.view_any': 'Конфигурации клиентов',
  'configurations.manage': 'Управление статусами конфигураций',
  'appointments.view_any': 'Все записи на ТО',
  'garage.view_any': 'Авто клиентов в гараже',
  'documents.view_any': 'Документы по чужим сделкам',
  'news.manage': 'Редакция новостей',
  'admin.order_statuses': 'Справочник статусов заказа',
  'admin.roles_view': 'Справочник ролей',
  'catalog.manage': 'Каталог автомобилей (бренды и др.)',
  'service.manage': 'Услуги ТО и филиалы',
};

export const ROLE_TITLE_RU = {
  customer: 'Клиент',
  manager: 'Менеджер',
  service_advisor: 'Мастер-приёмщик',
  admin: 'Администратор',
};

function permissionsFromSubject(subject) {
  if (subject && typeof subject === 'object' && Array.isArray(subject.permissions)) {
    return subject.permissions;
  }
  return null;
}

function roleFromSubject(subject) {
  if (subject && typeof subject === 'object') {
    return subject.role || '';
  }
  return typeof subject === 'string' ? subject : '';
}

/**
 * Checks permission using server-provided user.permissions when available,
 * otherwise falls back to the static role matrix.
 * @param {string|{role?: string, permissions?: string[]}} subject
 */
export function hasPermission(subject, permission) {
  if (!permission) return false;

  const live = permissionsFromSubject(subject);
  if (live) {
    return live.includes(permission);
  }

  const role = roleFromSubject(subject);
  if (!role) return false;
  const fallback = DEFAULT_ROLE_PERMISSIONS[role] || [];
  return fallback.includes(permission);
}

/** True when user is staff according to API or role fallback. */
export function isStaffUser(user) {
  if (!user) return false;
  if (typeof user.is_staff === 'boolean') return user.is_staff;
  const role = user.role || '';
  return role === 'admin' || role === 'manager' || role === 'service_advisor';
}
