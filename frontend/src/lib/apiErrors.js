/**
 * Normalizes API failures (axios, custom rejections, react-query) into a user-facing string.
 * Backend returns JSON: { success: false, error: "..." }.
 */

const NETWORK_FALLBACK_RU = 'Не удалось связаться с сервером. Проверьте подключение к интернету.';
const GENERIC_FALLBACK_RU = 'Операция не выполнена. Попробуйте позже.';

/** Maps known English API messages to Russian (backend handler strings). */
const BACKEND_TO_RU = [
  [/current password is incorrect/i, 'Неверный текущий пароль'],
  [/invalid credentials/i, 'Неверный email или пароль'],
  [/password must be at least 6 characters/i, 'Пароль: не менее 6 символов'],
  [/password is too long/i, 'Пароль слишком длинный'],
  [/first_name is required/i, 'Имя обязательно (до 100 символов)'],
  [/last_name is required/i, 'Фамилия обязательна (до 100 символов)'],
  [/email is required/i, 'Некорректный email'],
  [/phone is too long/i, 'Телефон слишком длинный (до 30 символов)'],
  [/invalid request body/i, 'Некорректные данные запроса'],
  [/user not found/i, 'Пользователь не найден'],
  [/email already exists/i, 'Пользователь с таким email уже зарегистрирован'],
  [
    /new password must be different from the current password/i,
    'Новый пароль должен отличаться от текущего',
  ],
  [/vin must be exactly 17 characters/i, 'VIN — ровно 17 символов'],
  [
    /vin must contain only letters/i,
    'Допустимы буквы A–Z (кроме I, O, Q) и цифры',
  ],
  [/vin already exists/i, 'Этот VIN уже зарегистрирован в системе'],
  [/invalid vehicle year/i, 'Укажите корректный год выпуска'],
  [/mileage must be non-negative/i, 'Пробег не может быть отрицательным'],
  [/appointment_date is required/i, 'Укажите дату и время записи'],
  [/appointment must be in the future/i, 'Выберите дату и время в будущем'],
  [
    /appointment date is too far in the future/i,
    'Дата записи слишком далеко вперёди (не более года)',
  ],
  [/description is too long/i, 'Описание слишком длинное'],
  [/at least one service type is required/i, 'Выберите хотя бы одну услугу'],
  [/branch is not active/i, 'Выбранный филиал недоступен'],
  [/branch not found/i, 'Филиал не найден'],
  [/user car not found/i, 'Автомобиль не найден'],
  [/user car does not belong to user/i, 'Этот автомобиль недоступен для записи'],
  [/service type .* not found or not available/i, 'Одна из услуг недоступна'],
  [/this time slot is no longer available/i, 'Это время уже занято. Выберите другой слот'],
  [/branch is closed on this day/i, 'В этот день филиал не работает'],
  [/appointment outside branch working hours/i, 'Время вне графика работы филиала'],
  [/invalid appointment time slot/i, 'Некорректное время записи'],
  [/invalid date format/i, 'Некорректная дата'],
  [/invalid service duration/i, 'Некорректная длительность услуг'],
  [/failed to get branch/i, 'Не удалось загрузить филиал'],
];

/** Turns a backend `error` string into Russian when we know the English phrase. */
export function formatBackendErrorMessage(text) {
  const t = String(text || '').trim();
  if (!t) return '';
  for (const [re, ru] of BACKEND_TO_RU) {
    if (re.test(t)) return ru;
  }
  return t;
}

function extractPayload(err) {
  if (!err || typeof err !== 'object') return null;
  if (err.data !== undefined) return err.data;
  if (err.response?.data !== undefined) return err.response.data;
  return null;
}

function rawMessageFromPayload(data) {
  if (!data) return '';
  if (typeof data === 'string') return data;
  if (typeof data.error === 'string' && data.error) return data.error;
  if (typeof data.message === 'string' && data.message) return data.message;
  return '';
}

/**
 * @param {unknown} err - Rejection from apiClient, axios, or Error
 * @param {string} [fallback] - Russian message if nothing usable is found
 */
export function getApiErrorMessage(err, fallback = GENERIC_FALLBACK_RU) {
  const data = extractPayload(err);
  let raw = rawMessageFromPayload(data);

  if (!raw && err && typeof err === 'object' && typeof err.message === 'string') {
    const m = err.message.trim();
    if (m && !/^Request failed with status code \d+$/i.test(m) && !/^Network Error$/i.test(m)) {
      raw = m;
    }
  }

  if (!raw) {
    if (err && typeof err === 'object' && err.status === 0) {
      return NETWORK_FALLBACK_RU;
    }
    return fallback;
  }

  const lower = raw.toLowerCase();
  if (lower.includes('network error') || lower.includes('connection')) {
    return NETWORK_FALLBACK_RU;
  }

  const out = formatBackendErrorMessage(raw);
  if (out.length > 280) {
    return fallback;
  }
  return out;
}
