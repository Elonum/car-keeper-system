/**
 * Format an ISO instant in a specific IANA timezone (branch schedule).
 */
export function formatAppointmentInTimeZone(isoString, timeZone, options) {
  const tz = timeZone || 'Europe/Moscow';
  return new Intl.DateTimeFormat('ru-RU', { timeZone: tz, ...options }).format(new Date(isoString));
}

export function appointmentServiceTypeIds(appt) {
  const types = appt?.service_types;
  if (!Array.isArray(types) || types.length === 0) return [];
  return types.map((t) => t.service_type_id || t.id).filter(Boolean);
}
