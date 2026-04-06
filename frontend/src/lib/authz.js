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

export function hasPermission(role, permission) {
  if (!role || !permission) return false;
  const perms = DEFAULT_ROLE_PERMISSIONS[role] || [];
  return perms.includes(permission);
}

