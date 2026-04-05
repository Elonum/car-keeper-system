# Схема БД CarKeeper

## Источник правды

Файл **`schema.sql`** описывает целевую схему. При изменениях в продакшене или разработке:

1. Выполните нужные SQL-команды в своей базе (через `psql`, GUI или скрипт).
2. Внесите те же изменения в **`schema.sql`** в репозитории, чтобы новые окружения и документация не расходились с реальностью.

Первичная установка пустой БД:

```bash
psql -U postgres -d carkeeper -f schema.sql
```

При необходимости затем: `seed.sql`, `catalog_view.sql`.

### Документы (вложения)

Таблица **`documents`** хранит метаданные; сами байты лежат на диске в каталоге **`DOCUMENT_STORAGE_ROOT`** (переменная окружения API, по умолчанию `./data/documents`). В поле `file_path` — только внутренний ключ (имя файла в этом каталоге), без абсолютного пути.

---

## Примеры ручных изменений

Все примеры иллюстративны: подставьте свои имена, типы и ограничения.

### Добавить nullable-колонку

```sql
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS middle_name varchar(100);

COMMENT ON COLUMN users.middle_name IS 'Отчество (опционально)';
```

После этого добавьте тот же `ADD COLUMN` в `schema.sql` в определение таблицы `users`.

### Добавить колонку с значением по умолчанию

```sql
ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS priority smallint NOT NULL DEFAULT 0;

ALTER TABLE orders
  ADD CONSTRAINT chk_orders_priority CHECK (priority BETWEEN 0 AND 3);
```

### Индекс под частый запрос

```sql
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_status_created
  ON orders (status, created_at DESC);
```

`CONCURRENTLY` удобен на живой базе (без долгой блокировки записи); в `schema.sql` обычно указывают обычный `CREATE INDEX` для новых развёртываний.

### Новая таблица со связью

```sql
CREATE TABLE promotion_codes (
    code_id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code         varchar(32) NOT NULL UNIQUE,
    discount_pct numeric(5,2) NOT NULL CHECK (discount_pct > 0 AND discount_pct <= 100),
    valid_until  timestamptz,
    created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_promotion_codes_valid_until ON promotion_codes(valid_until);
```

Скопируйте полный блок в `schema.sql` (и при необходимости добавьте триггер `updated_at` по аналогии с другими таблицами).

### Изменить тип или ограничение

```sql
ALTER TABLE news
  ALTER COLUMN title TYPE varchar(500);
```

### Удалить колонку (осторожно: потеря данных)

```sql
ALTER TABLE users DROP COLUMN IF EXISTS legacy_field;
```

### Переименовать колонку

```sql
ALTER TABLE branches RENAME COLUMN phone TO phone_main;
```

---

## Замечания

- Делайте бэкап перед нетривиальными `ALTER` на продакшене.
- После изменения `schema.sql` проверьте, что код приложения (модели, репозитории) согласован с новыми полями и типами.
