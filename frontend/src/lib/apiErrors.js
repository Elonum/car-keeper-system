/**
 * Maps backend `error` strings (English, from apperr / handlers) to Russian for the UI.
 * Unknown strings are returned as-is (e.g. already-Russian validation text).
 */

const NETWORK_FALLBACK_RU =
  'Не удалось связаться с сервером. Проверьте подключение к интернету.';
const GENERIC_FALLBACK_RU = 'Операция не выполнена. Попробуйте позже.';

/**
 * Order: more specific patterns first. Case-insensitive flags on each regex.
 */
const BACKEND_TO_RU = [
  // JSON / HTTP limits (handler)
  [/^\s*request body is too large\s*$/i, 'Слишком большой запрос. Уменьшите объём данных.'],
  [/^\s*invalid json\s*$/i, 'Некорректный формат данных (JSON).'],
  [/^\s*request body is required\s*$/i, 'Тело запроса обязательно.'],

  // Generic API (apperr / HandleError)
  [
    /^\s*an unexpected error occurred\. please try again later\.\s*$/i,
    'Произошла ошибка сервера. Попробуйте позже.',
  ],
  [/^\s*resource not found\s*$/i, 'Данные не найдены.'],
  [/^\s*access denied\s*$/i, 'Доступ запрещён.'],

  // Auth
  [/^\s*invalid email or password\s*$/i, 'Неверный email или пароль.'],
  [/^\s*this email is already registered\s*$/i, 'Пользователь с таким email уже зарегистрирован.'],
  [/^\s*user not found\s*$/i, 'Пользователь не найден.'],
  [/^\s*current password is incorrect\s*$/i, 'Неверный текущий пароль.'],

  // Orders & configuration (apperr)
  [/^\s*configuration not found\s*$/i, 'Конфигурация не найдена.'],
  [
    /^\s*configuration does not belong to your account\s*$/i,
    'Эта конфигурация не принадлежит вашему аккаунту.',
  ],
  [
    /^\s*configuration cannot be ordered in its current status\s*$/i,
    'Заказ недоступен для этой конфигурации в текущем статусе.',
  ],
  [/^\s*status is required\s*$/i, 'Укажите статус.'],
  [/^\s*unknown order status\s*$/i, 'Неизвестный статус заказа.'],
  [/^\s*this order status is not available\s*$/i, 'Этот статус заказа недоступен.'],
  [/^\s*invalid order status\s*$/i, 'Некорректный статус заказа.'],

  [/^\s*status cannot be empty\s*$/i, 'Статус не может быть пустым.'],
  [/^\s*invalid configuration status\s*$/i, 'Некорректный статус конфигурации.'],
  [/^\s*only draft configurations can be updated\s*$/i, 'Редактировать можно только черновик.'],
  [/^\s*only draft configurations can be deleted\s*$/i, 'Удалить можно только черновик.'],

  [
    /^\s*invalid trim or catalog data\s*$/i,
    'Некорректная комплектация или данные каталога.',
  ],
  [
    /^\s*invalid color or catalog data\s*$/i,
    'Некорректный цвет или данные каталога.',
  ],
  [/^\s*invalid options selection\s*$/i, 'Некорректный набор опций.'],

  // Admin order statuses (still useful if admin UI calls API later)
  [/^\s*code is required\s*$/i, 'Укажите код статуса.'],
  [/^\s*customer_label_ru is required\s*$/i, 'Укажите подпись для клиента (RU).'],
  [/^\s*customer_label_ru cannot be empty\s*$/i, 'Подпись для клиента не может быть пустой.'],
  [/^\s*order status code already exists\s*$/i, 'Такой код статуса заказа уже есть.'],
  [
    /^\s*cannot delete status: orders still use this code\s*$/i,
    'Нельзя удалить статус: есть заказы с этим кодом.',
  ],

  // Profile / garage
  [/^\s*this vin is already registered\s*$/i, 'Этот VIN уже зарегистрирован в системе.'],
  [/^\s*invalid vehicle year\s*$/i, 'Укажите корректный год выпуска.'],
  [/^\s*mileage must be non-negative\s*$/i, 'Пробег не может быть отрицательным.'],

  // Handler validation (English)
  [/^\s*invalid order id\s*$/i, 'Некорректный номер заказа.'],
  [/^\s*invalid configuration id\s*$/i, 'Некорректный идентификатор конфигурации.'],
  [/^\s*invalid trim_id\s*$/i, 'Некорректный trim_id.'],
  [/^\s*trim_id is required\s*$/i, 'Укажите комплектацию (trim_id).'],
  [/^\s*invalid trim_id format\s*$/i, 'Некорректный формат trim_id.'],
  [/^\s*invalid color_id format\s*$/i, 'Некорректный формат color_id.'],
  [/^\s*invalid option_id format\s*$/i, 'Некорректный формат option_id.'],
  [/^\s*trim_id is required for full update\s*$/i, 'Для полного обновления укажите trim_id.'],
  [/^\s*color_id is required for full update\s*$/i, 'Для полного обновления укажите color_id.'],
  [/^\s*invalid news id\s*$/i, 'Некорректный идентификатор новости.'],
  [/^\s*invalid appointment id\s*$/i, 'Некорректный идентификатор записи.'],
  [/^\s*invalid status id\s*$/i, 'Некорректный идентификатор статуса.'],
  [/^\s*invalid user car id\s*$/i, 'Некорректный идентификатор автомобиля.'],
  [/^\s*invalid document id\s*$/i, 'Некорректный идентификатор документа.'],
  [/^\s*invalid branch id\s*$/i, 'Некорректный идентификатор филиала.'],
  [/^\s*invalid query id\s*$/i, 'Некорректный параметр запроса (id).'],
  [/^\s*current_password is required\s*$/i, 'Введите текущий пароль.'],

  // Catalog
  [/^\s*brand_id is required\s*$/i, 'Укажите марку (brand_id).'],
  [/^\s*invalid brand_id\s*$/i, 'Некорректный brand_id.'],
  [/^\s*model_id is required\s*$/i, 'Укажите модель (model_id).'],
  [/^\s*invalid model_id\s*$/i, 'Некорректный model_id.'],
  [/^\s*invalid trim id\s*$/i, 'Некорректный идентификатор комплектации.'],

  // Service / appointments (service layer messages still exposed via HandleError)
  [/^\s*appointment_date is required\s*$/i, 'Укажите дату и время записи.'],
  [/^\s*date is required\s*$/i, 'Укажите дату.'],
  [/^\s*service_type_ids is required\s*$/i, 'Выберите хотя бы одну услугу.'],
  [/^\s*invalid service_type_ids\s*$/i, 'Некорректный список услуг.'],
  [/^\s*user car does not belong to user\s*$/i, 'Этот автомобиль недоступен для записи.'],
  [/^\s*branch is not active\s*$/i, 'Выбранный филиал недоступен.'],
  [/^\s*branch not found\s*$/i, 'Филиал не найден.'],
  [/^\s*at least one service type is required\s*$/i, 'Выберите хотя бы одну услугу.'],
  [
    /one or more service types are not available/i,
    'Одна или несколько услуг недоступны.',
  ],
  [/^\s*appointment must be in the future\s*$/i, 'Выберите дату и время в будущем.'],
  [
    /^\s*appointment date is too far in the future\s*$/i,
    'Дата записи слишком далеко вперёди (не более года).',
  ],
  [/^\s*description is too long\s*$/i, 'Описание слишком длинное.'],
  [/^\s*this time slot is no longer available\s*$/i, 'Это время уже занято. Выберите другой слот.'],
  [/^\s*branch is closed on this day\s*$/i, 'В этот день филиал не работает.'],
  [
    /^\s*appointment outside branch working hours\s*$/i,
    'Время вне графика работы филиала.',
  ],
  [/^\s*invalid appointment time slot\s*$/i, 'Некорректное время записи.'],
  [/^\s*invalid date format\s*$/i, 'Некорректная дата.'],
  [/^\s*invalid service duration\s*$/i, 'Некорректная длительность услуг.'],
  [/^\s*failed to get branch\s*$/i, 'Не удалось загрузить данные филиала.'],
  [
    /^\s*only scheduled appointments can be rescheduled\s*$/i,
    'Перенести можно только запланированную запись.',
  ],
  [
    /appointment not found or cannot be rescheduled/i,
    'Запись не найдена или её нельзя перенести.',
  ],
  [/^\s*appointment could not be updated\s*$/i, 'Не удалось обновить запись. Попробуйте ещё раз.'],
  [
    /^\s*not allowed to reschedule this appointment\s*$/i,
    'Нет прав на перенос этой записи.',
  ],
  [/^\s*appointment already cancelled\s*$/i, 'Запись уже отменена.'],
  [/^\s*cannot cancel completed appointment\s*$/i, 'Нельзя отменить завершённую запись.'],

  // Documents (handler)
  [/^\s*invalid multipart form\s*$/i, 'Некорректная форма загрузки файла.'],
  [/^\s*file field is required\s*$/i, 'Выберите файл для загрузки.'],
  [
    /^\s*invalid order_id or service_appointment_id\s*$/i,
    'Некорректный order_id или service_appointment_id.',
  ],
  [
    /^\s*could not determine file size/i,
    'Не удалось определить размер файла. Загрузите файл с известным размером.',
  ],
  [/^\s*file too large\s*$/i, 'Файл слишком большой.'],

  // News (handler)
  [/^\s*title is required \(max 255 characters\)\s*$/i, 'Заголовок обязателен (до 255 символов).'],
  [/^\s*content is required\s*$/i, 'Текст новости обязателен.'],

  // Legacy / validation phrases (still possible from validators or old paths)
  [/^\s*invalid credentials\s*$/i, 'Неверный email или пароль.'],
  [/^\s*email already exists\s*$/i, 'Пользователь с таким email уже зарегистрирован.'],
  [/^\s*invalid request body\s*$/i, 'Некорректные данные запроса.'],
  [/^\s*password must be at least 6 characters\s*$/i, 'Пароль: не менее 6 символов.'],
  [/^\s*password is too long\s*$/i, 'Пароль слишком длинный.'],
  [/^\s*first_name is required/i, 'Имя обязательно (до 100 символов).'],
  [/^\s*last_name is required/i, 'Фамилия обязательна (до 100 символов).'],
  [/^\s*email is required/i, 'Некорректный email.'],
  [/^\s*phone is too long/i, 'Телефон слишком длинный (до 30 символов).'],
  [
    /^\s*new password must be different from the current password\s*$/i,
    'Новый пароль должен отличаться от текущего.',
  ],
  [/vin must be exactly 17 characters/i, 'VIN — ровно 17 символов.'],
  [/vin must contain only letters/i, 'Допустимы буквы A–Z (кроме I, O, Q) и цифры.'],
  [/^\s*vin already exists\s*$/i, 'Этот VIN уже зарегистрирован в системе.'],
  [/^\s*user car not found\s*$/i, 'Автомобиль не найден.'],
  [
    /^\s*service type .* not found or not available/i,
    'Одна из услуг недоступна.',
  ],
];

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

const MAX_USER_MESSAGE_LEN = 280;

/**
 * @param {unknown} err
 * @param {string} [fallback]
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
  if (out.length > MAX_USER_MESSAGE_LEN) {
    return fallback;
  }
  return out;
}
