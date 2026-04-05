/** Client-side rules aligned with backend `validate` + DB checks for user_cars. */

const VIN_RE = /^[A-HJ-NPR-Z0-9]{17}$/;

export function normalizeVIN(value) {
  return String(value ?? '')
    .toUpperCase()
    .replace(/\s/g, '');
}

export function validateVIN(value) {
  const v = normalizeVIN(value);
  if (v.length !== 17) {
    return 'VIN — ровно 17 символов';
  }
  if (!VIN_RE.test(v)) {
    return 'Допустимы буквы A–Z (кроме I, O, Q) и цифры';
  }
  return null;
}

export function validateVehicleYear(yearNum) {
  const y = Number(yearNum);
  const maxY = new Date().getFullYear() + 1;
  if (!Number.isInteger(y) || y < 1900 || y > maxY) {
    return `Год — от 1900 до ${maxY}`;
  }
  return null;
}

export function validateMileage(value) {
  const n = Number(value);
  if (value === '' || value == null) {
    return null;
  }
  if (!Number.isFinite(n) || n < 0 || !Number.isInteger(n)) {
    return 'Пробег — целое число от 0';
  }
  return null;
}
