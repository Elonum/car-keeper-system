import { API_BASE_ORIGIN } from '@/api/client';

export function resolveApiAssetUrl(rawUrl) {
  const raw = String(rawUrl || '').trim();
  if (!raw) return '';
  if (raw.startsWith('http://') || raw.startsWith('https://')) return raw;
  if (raw.startsWith('/')) return `${API_BASE_ORIGIN}${raw}`;
  return `${API_BASE_ORIGIN}/${raw}`;
}

/** Stable React key for model images — changes when backend image_key changes. */
export function modelImageCacheKey(rawUrl, fallbackId = '') {
  const raw = String(rawUrl || '').trim();
  if (!raw) return String(fallbackId || '');
  const query = raw.split('?')[1] || '';
  const version = new URLSearchParams(query).get('v');
  return version || raw;
}
