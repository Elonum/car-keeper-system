/** Catalog model image policy (aligned with backend upload/image.go). */

export const MODEL_IMAGE_MAX_BYTES = 5 * 1024 * 1024;
export const MODEL_IMAGE_ACCEPT = 'image/jpeg,image/png,image/webp,.jpg,.jpeg,.png,.webp';
export const MODEL_IMAGE_MIME_TYPES = ['image/jpeg', 'image/png', 'image/webp'];
const MODEL_IMAGE_EXTENSIONS = ['jpg', 'jpeg', 'png', 'webp'];

function fileExtension(name) {
  const parts = String(name || '').toLowerCase().split('.');
  return parts.length > 1 ? parts.pop() : '';
}

export function isAllowedModelImageFile(file) {
  if (!file) return false;
  if (MODEL_IMAGE_MIME_TYPES.includes(file.type)) return true;
  return MODEL_IMAGE_EXTENSIONS.includes(fileExtension(file.name));
}

export function validateModelImageFile(file) {
  if (!file) return null;
  if (!isAllowedModelImageFile(file)) {
    return 'Допустимы JPG, PNG или WEBP';
  }
  if (file.size > MODEL_IMAGE_MAX_BYTES) {
    return 'Размер файла не более 5 МБ';
  }
  return null;
}
