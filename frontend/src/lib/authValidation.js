/**
 * Client-side validation aligned with backend rules (see backend handler auth).
 * Does not replace server validation.
 */

const EMAIL_MAX = 255;
const NAME_MAX = 100;
const PASSWORD_MIN = 6;
const PASSWORD_MAX = 128;
/** Matches DB `users.phone` varchar(30) */
const PHONE_MAX = 30;

/** Max input lengths aligned with PostgreSQL column sizes / handler checks */
export const FIELD_LIMITS = {
  EMAIL: EMAIL_MAX,
  NAME: NAME_MAX,
  PHONE: PHONE_MAX,
  PASSWORD: PASSWORD_MAX,
};

/** UI-only caps (search) or safe upper bounds for unbounded TEXT columns */
export const UI_LIMITS = {
  CATALOG_SEARCH: 200,
  /** service_appointments.description is TEXT; limit payload size */
  APPOINTMENT_DESCRIPTION: 10000,
};

export function normalizeEmail(email) {
  return String(email ?? '')
    .trim()
    .toLowerCase();
}

export function validateEmail(email) {
  const t = normalizeEmail(email);
  if (!t) {
    return 'Укажите email';
  }
  if (t.length > EMAIL_MAX) {
    return 'Email слишком длинный';
  }
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!re.test(t)) {
    return 'Введите корректный email';
  }
  return null;
}

export function validatePassword(password) {
  if (password == null || password.length === 0) {
    return 'Введите пароль';
  }
  if (password.length < PASSWORD_MIN) {
    return `Пароль — не менее ${PASSWORD_MIN} символов`;
  }
  if (password.length > PASSWORD_MAX) {
    return 'Пароль слишком длинный';
  }
  return null;
}

export function validatePersonName(value, fieldLabel) {
  const t = String(value ?? '').trim();
  if (!t) {
    return `${fieldLabel} обязательно`;
  }
  if (t.length > NAME_MAX) {
    return `Не более ${NAME_MAX} символов`;
  }
  if (!/^[\p{L}\s\-'.]+$/u.test(t)) {
    return 'Допустимы буквы, пробелы, дефис и точка';
  }
  return null;
}

const PHONE_CHARS = /^[+\d\s().-]+$/;
const MIN_PHONE_DIGITS = 10;
const MAX_PHONE_DIGITS = 15;

/** Compact storage form aligned with backend PhonePtr normalization */
export function normalizePhoneOptional(phone) {
  const raw = String(phone ?? '').trim();
  if (!raw) return null;
  const digits = raw.replace(/\D/g, '');
  if (!digits) return null;
  if (/\+/.test(raw)) return `+${digits}`;
  return digits;
}

/** Optional phone: E.164-ish, 10–15 digits when provided */
export function validatePhoneOptional(phone) {
  const raw = String(phone ?? '').trim();
  if (!raw) {
    return null;
  }
  if (raw.length > PHONE_MAX) {
    return `Номер не длиннее ${PHONE_MAX} символов`;
  }
  if (!PHONE_CHARS.test(raw)) {
    return 'Номер может содержать только цифры и символы + ( ) - .';
  }
  const digits = raw.replace(/\D/g, '');
  if (digits.length < MIN_PHONE_DIGITS || digits.length > MAX_PHONE_DIGITS) {
    return 'Введите корректный номер (10–15 цифр)';
  }
  return null;
}

export function validatePasswordConfirm(password, confirm) {
  if (confirm == null || confirm.length === 0) {
    return 'Подтвердите пароль';
  }
  if (password !== confirm) {
    return 'Пароли не совпадают';
  }
  return null;
}
