/**
 * Клиентские фильтры ЛК (без дублирования логики в компонентах).
 */

function norm(s) {
  return String(s ?? '')
    .toLowerCase()
    .trim();
}

/**
 * @param {Array<Record<string, unknown>>|undefined} orders
 * @param {{ search: string, status: string }} filters status: 'all' | order status code
 */
export function filterOrdersForCabinet(orders, { search, status }) {
  const list = Array.isArray(orders) ? orders : [];
  const q = norm(search);
  const st = status === 'all' || !status ? null : status;

  return list.filter((o) => {
    if (st && o.status !== st) return false;
    if (!q) return true;
    const id = norm(o.order_id ?? o.id ?? '');
    const shortId = id.replace(/-/g, '');
    const qCompact = q.replace(/-/g, '');
    const summary = norm(o.configuration_summary);
    const price = norm(o.final_price ?? o.total_price ?? '');
    return (
      id.includes(q) ||
      shortId.includes(qCompact) ||
      summary.includes(q) ||
      price.includes(q)
    );
  });
}

/**
 * @param {Array<Record<string, unknown>>|undefined} appointments
 * @param {{ search: string, status: string }} filters
 */
export function filterAppointmentsForCabinet(appointments, { search, status }) {
  const list = Array.isArray(appointments) ? appointments : [];
  const q = norm(search);
  const st = status === 'all' || !status ? null : status;

  return list.filter((a) => {
    if (st && a.status !== st) return false;
    if (!q) return true;
    const vin = norm(a.user_car_vin ?? a.car_display);
    const branch = norm(a.branch_name ?? a.branch_display);
    const addr = norm(a.branch_address);
    const desc = norm(a.description);
    const id = norm(a.service_appointment_id ?? a.id ?? '');
    return (
      vin.includes(q) ||
      branch.includes(q) ||
      addr.includes(q) ||
      desc.includes(q) ||
      id.includes(q) ||
      String(id).replace(/-/g, '').includes(q.replace(/-/g, ''))
    );
  });
}
