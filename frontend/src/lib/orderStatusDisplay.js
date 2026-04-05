/**
 * Order status presentation: labels come from API (order.status_label or GET /order-statuses).
 * Only Tailwind visual tokens are keyed by code here — not user-facing copy.
 */

const VISUAL_BY_CODE = {
  pending: 'bg-amber-100 text-amber-800 border-amber-200',
  approved: 'bg-blue-100 text-blue-800 border-blue-200',
  paid: 'bg-emerald-100 text-emerald-800 border-emerald-200',
  completed: 'bg-green-100 text-green-800 border-green-200',
  cancelled: 'bg-red-100 text-red-800 border-red-200',
};

export function orderStatusBadgeClassName(code) {
  const c = String(code ?? '').toLowerCase();
  return VISUAL_BY_CODE[c] ?? 'bg-slate-100 text-slate-700 border-slate-200';
}

/** @param {Array<{ code: string, customer_label_ru: string }>|undefined} statuses */
export function buildOrderStatusLabelMap(statuses) {
  const m = new Map();
  for (const s of statuses || []) {
    if (s?.code && s.customer_label_ru != null) {
      m.set(String(s.code), String(s.customer_label_ru));
    }
  }
  return m;
}

/**
 * @param {{ status?: string, status_label?: string }|undefined} order
 * @param {Map<string, string>} labelByCode from public order-statuses
 */
export function resolveOrderStatusLabel(order, labelByCode) {
  const code = order?.status;
  const direct = order?.status_label?.trim();
  if (direct) return direct;
  if (code != null && code !== '' && labelByCode?.has(code)) {
    return labelByCode.get(code);
  }
  if (code != null && code !== '') return String(code);
  return '—';
}
