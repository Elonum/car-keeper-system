/**
 * API list endpoints may return `null` when empty (Go nil slice).
 * React Query default `data = []` only applies when data is undefined.
 */
export function asArray(value) {
  return Array.isArray(value) ? value : [];
}
