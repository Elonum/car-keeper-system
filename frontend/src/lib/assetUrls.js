import { API_BASE_ORIGIN } from '@/api/client';

export function resolveApiAssetUrl(rawUrl) {
  const raw = String(rawUrl || '').trim();
  if (!raw) return '';
  if (raw.startsWith('http://') || raw.startsWith('https://')) return raw;
  if (raw.startsWith('/')) return `${API_BASE_ORIGIN}${raw}`;
  return `${API_BASE_ORIGIN}/${raw}`;
}

