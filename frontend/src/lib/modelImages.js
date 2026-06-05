/** Catalog model image policy (aligned with backend upload/image.go). */

export const MODEL_IMAGE_MAX_BYTES = 5 * 1024 * 1024;
export const MODEL_IMAGE_ACCEPT = 'image/jpeg,image/png,image/webp';
export const MODEL_IMAGE_MIME_TYPES = ['image/jpeg', 'image/png', 'image/webp'];

export function isAllowedModelImageFile(file) {
  if (!file) return false;
  return MODEL_IMAGE_MIME_TYPES.includes(file.type);
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
