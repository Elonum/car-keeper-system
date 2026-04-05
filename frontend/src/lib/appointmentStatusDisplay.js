/**
 * Подписи статусов записи на ТО. Должны соответствовать CHECK в БД (service_appointments.status).
 * Отдельного публичного API справочника нет — один модуль вместо размазанных строк в UI.
 */

const VISUAL_BY_CODE = {
  scheduled: 'bg-sky-100 text-sky-800 border-sky-200',
  completed: 'bg-green-100 text-green-800 border-green-200',
  cancelled: 'bg-red-100 text-red-800 border-red-200',
};

/** @type {Record<string, string>} */
export const APPOINTMENT_STATUS_LABEL_RU = {
  scheduled: 'Запланировано',
  completed: 'Завершено',
  cancelled: 'Отменено',
};

export const APPOINTMENT_STATUS_CODES = Object.keys(APPOINTMENT_STATUS_LABEL_RU);

export function appointmentStatusClassName(code) {
  const c = String(code ?? '').toLowerCase();
  return VISUAL_BY_CODE[c] ?? 'bg-slate-100 text-slate-700 border-slate-200';
}

export function resolveAppointmentStatusLabel(code) {
  const c = String(code ?? '');
  if (APPOINTMENT_STATUS_LABEL_RU[c]) return APPOINTMENT_STATUS_LABEL_RU[c];
  return c || '—';
}
