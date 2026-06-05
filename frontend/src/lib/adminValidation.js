/**
 * Client-side validation for admin/catalog/news forms.
 * Limits align with backend validate package and PostgreSQL schema.
 */

export const ADMIN_LIMITS = {
  BRAND_NAME: 150,
  BRAND_COUNTRY: 100,
  MODEL_NAME: 150,
  MODEL_SEGMENT: 100,
  MODEL_DESCRIPTION: 2000,
  SERVICE_NAME: 150,
  SERVICE_DESCRIPTION: 2000,
  SERVICE_PRICE_MAX: 99_999_999.99,
  SERVICE_DURATION_MAX: 1440,
  NEWS_TITLE: 255,
  NEWS_CONTENT: 50_000,
  ORDER_STATUS_CODE: 32,
  ORDER_STATUS_LABEL: 120,
  ORDER_STATUS_DESCRIPTION: 2000,
  ORDER_STATUS_SORT_MIN: -10_000,
  ORDER_STATUS_SORT_MAX: 100_000,
};

export const SERVICE_CATEGORIES = [
  'maintenance',
  'repair',
  'diagnostics',
  'detailing',
  'tires',
];

const CONTROL_CHARS = /[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]/;

function hasControlChars(value, allowNewlines = false) {
  const text = String(value ?? '');
  if (!allowNewlines && /[\r\n\t]/.test(text)) return true;
  return CONTROL_CHARS.test(text);
}

export function validateRequiredText(value, label, maxLen, { multiline = false } = {}) {
  const t = String(value ?? '').trim();
  if (!t) return `${label} обязательно`;
  if (t.length > maxLen) return `${label}: не более ${maxLen} символов`;
  if (hasControlChars(t, multiline)) return `${label} содержит недопустимые символы`;
  return null;
}

export function validateOptionalText(value, label, maxLen, { multiline = false } = {}) {
  const t = String(value ?? '').trim();
  if (!t) return null;
  if (t.length > maxLen) return `${label}: не более ${maxLen} символов`;
  if (hasControlChars(t, multiline)) return `${label} содержит недопустимые символы`;
  return null;
}

export function validateBrandForm(form) {
  const errors = {};
  const name = validateRequiredText(form.name, 'Название бренда', ADMIN_LIMITS.BRAND_NAME);
  const country = validateRequiredText(form.country, 'Страна', ADMIN_LIMITS.BRAND_COUNTRY);
  if (name) errors.name = name;
  if (country) errors.country = country;
  return errors;
}

export function validateModelForm(form) {
  const errors = {};
  const brandId = String(form.brand_id ?? '').trim();
  const name = String(form.name ?? '').trim();
  if (!brandId) errors.brand_id = 'Выберите бренд';
  const nameErr = validateRequiredText(name, 'Название модели', ADMIN_LIMITS.MODEL_NAME);
  if (nameErr) errors.name = nameErr;
  const segmentErr = validateOptionalText(form.segment, 'Сегмент', ADMIN_LIMITS.MODEL_SEGMENT);
  if (segmentErr) errors.segment = segmentErr;
  const descErr = validateOptionalText(
    form.description,
    'Описание',
    ADMIN_LIMITS.MODEL_DESCRIPTION,
    { multiline: true }
  );
  if (descErr) errors.description = descErr;
  return errors;
}

export function validateServiceForm(form) {
  const errors = {};
  const nameErr = validateRequiredText(form.name, 'Название услуги', ADMIN_LIMITS.SERVICE_NAME);
  if (nameErr) errors.name = nameErr;
  if (!SERVICE_CATEGORIES.includes(form.category)) {
    errors.category = 'Выберите категорию из списка';
  }
  const price = Number(form.price);
  if (!Number.isFinite(price) || price < 0) {
    errors.price = 'Цена должна быть числом >= 0';
  } else if (price > ADMIN_LIMITS.SERVICE_PRICE_MAX) {
    errors.price = 'Цена слишком большая';
  }
  const durationRaw = form.duration_minutes;
  if (durationRaw !== '' && durationRaw != null) {
    const duration = Number(durationRaw);
    if (
      !Number.isFinite(duration) ||
      duration < 1 ||
      duration > ADMIN_LIMITS.SERVICE_DURATION_MAX
    ) {
      errors.duration_minutes = `Длительность: от 1 до ${ADMIN_LIMITS.SERVICE_DURATION_MAX} минут`;
    }
  }
  const descErr = validateOptionalText(
    form.description,
    'Описание',
    ADMIN_LIMITS.SERVICE_DESCRIPTION,
    { multiline: true }
  );
  if (descErr) errors.description = descErr;
  return errors;
}

export function validateOrderStatusForm(form) {
  const errors = {};
  const code = String(form.code ?? '').trim().toLowerCase();
  if (!code) {
    errors.code = 'Укажите код статуса';
  } else if (code.length > ADMIN_LIMITS.ORDER_STATUS_CODE) {
    errors.code = `Код: не более ${ADMIN_LIMITS.ORDER_STATUS_CODE} символов`;
  } else if (!/^[a-z][a-z0-9_]*$/.test(code)) {
    errors.code = 'Код: латиница, цифры и _ (начинается с буквы)';
  }
  const customerErr = validateRequiredText(
    form.customer_label_ru,
    'Подпись для клиента',
    ADMIN_LIMITS.ORDER_STATUS_LABEL
  );
  if (customerErr) errors.customer_label_ru = customerErr;
  const adminErr = validateOptionalText(
    form.admin_label_ru,
    'Подпись для админа',
    ADMIN_LIMITS.ORDER_STATUS_LABEL
  );
  if (adminErr) errors.admin_label_ru = adminErr;
  if (!Number.isFinite(form.sort_order)) {
    errors.sort_order = 'Укажите корректный порядок сортировки';
  } else if (
    form.sort_order < ADMIN_LIMITS.ORDER_STATUS_SORT_MIN ||
    form.sort_order > ADMIN_LIMITS.ORDER_STATUS_SORT_MAX
  ) {
    errors.sort_order = 'Порядок сортировки вне допустимого диапазона';
  }
  const descErr = validateOptionalText(
    form.description,
    'Описание',
    ADMIN_LIMITS.ORDER_STATUS_DESCRIPTION,
    { multiline: true }
  );
  if (descErr) errors.description = descErr;
  return errors;
}

export function validateNewsDraft(title, content) {
  const errors = {};
  const titleErr = validateRequiredText(title, 'Заголовок', ADMIN_LIMITS.NEWS_TITLE);
  const contentErr = validateRequiredText(
    content,
    'Текст новости',
    ADMIN_LIMITS.NEWS_CONTENT,
    { multiline: true }
  );
  if (titleErr) errors.title = titleErr;
  if (contentErr) errors.content = contentErr;
  return errors;
}
