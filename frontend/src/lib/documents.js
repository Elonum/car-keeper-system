/** Client-side document upload policy (keep in sync with backend internal/upload/mime.go). */

export const MAX_DOCUMENT_BYTES = 15 * 1024 * 1024;

export const SUPPORTED_DOCUMENT_EXTENSIONS = [
  '.pdf',
  '.jpg',
  '.jpeg',
  '.png',
  '.webp',
  '.doc',
  '.docx',
  '.xls',
  '.xlsx',
];

export const DOCUMENT_ACCEPT = SUPPORTED_DOCUMENT_EXTENSIONS.join(',');

export const SUPPORTED_FORMATS_LABEL =
  'PDF, JPEG, PNG, WebP, Word (.doc/.docx), Excel (.xls/.xlsx)';

function extensionOf(name) {
  if (!name || typeof name !== 'string') return '';
  const i = name.lastIndexOf('.');
  if (i < 0) return '';
  return name.slice(i).toLowerCase();
}

export function formatFileSize(bytes) {
  if (bytes == null || Number.isNaN(Number(bytes))) return '';
  const n = Number(bytes);
  if (n < 1024) return `${n} Б`;
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} КБ`;
  return `${(n / (1024 * 1024)).toFixed(1)} МБ`;
}

/**
 * @param {File|null|undefined} file
 * @returns {{ ok: true } | { ok: false, message: string }}
 */
export function validateDocumentFile(file) {
  if (!file) {
    return { ok: false, message: 'Выберите файл' };
  }
  const ext = extensionOf(file.name);
  if (!ext || !SUPPORTED_DOCUMENT_EXTENSIONS.includes(ext)) {
    return {
      ok: false,
      message: `Неподдерживаемый формат. Допустимо: ${SUPPORTED_FORMATS_LABEL}`,
    };
  }
  if (file.size <= 0) {
    return { ok: false, message: 'Файл пустой' };
  }
  if (file.size > MAX_DOCUMENT_BYTES) {
    return {
      ok: false,
      message: `Файл слишком большой (макс. ${formatFileSize(MAX_DOCUMENT_BYTES)})`,
    };
  }
  return { ok: true };
}
