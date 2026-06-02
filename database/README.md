# База данных CarKeeper

## Источник правды

Файл **`schema.sql`** — полная целевая схема (таблицы, RBAC, статусы заказов, триггеры).

При изменениях:

1. Выполните SQL на своей БД.
2. Внесите те же изменения в **`schema.sql`**, чтобы новые окружения не расходились с продакшеном.

Миграционного раннера нет — только ручные правки и актуальный `schema.sql`.

## Первичная установка

```bash
psql -U postgres -d carkeeper -f schema.sql
psql -U postgres -d carkeeper -f catalog_view.sql
psql -U postgres -d carkeeper -f seed.sql
```

## Очистка данных (без пересоздания схемы)

Сохраняются справочники из `schema.sql` (роли, права, статусы заказов). Удаляются пользователи, каталог, заказы, записи и т.д.

```bash
psql -U postgres -d carkeeper -f reset_data.sql
psql -U postgres -d carkeeper -f seed.sql
```

Из корня репозитория:

```bash
# Linux / macOS
./scripts/db-reseed.sh

# Windows PowerShell
.\scripts\db-reseed.ps1
```

В Docker:

```bash
docker compose exec postgres psql -U postgres -d carkeeper -f /scripts/reset_data.sql
docker compose exec postgres psql -U postgres -d carkeeper -f /scripts/seed.sql
```

## Демо-данные (`seed.sql`)

| Сущность | Что показано |
|----------|----------------|
| Пользователи | admin, manager, service_advisor, 2 клиента |
| Филиалы | 2 активных + 1 неактивный (реконструкция) |
| Каталог | 5 брендов, комплектации, опции, недоступный trim |
| Конфигурации | черновик, подтверждена, в заказе |
| Заказы | pending, approved, paid |
| Гараж | 3 автомобиля с VIN и пробегом |
| Записи на ТО | завершена, 2 будущие, отменённая |
| Документы | метаданные к заказу и ТО (файлы — через UI) |
| Новости | 3 опубликованные, 1 черновик |

Пароль всех пользователей seed: **`password123`**.

## Инкрементальные правки (существующая БД)

### Колонки изображений моделей

```sql
ALTER TABLE models ADD COLUMN IF NOT EXISTS image_key varchar(128);
ALTER TABLE models ADD COLUMN IF NOT EXISTS image_mime varchar(100);
```

### Демо-авто клиента (если seed не перезаливали целиком)

```sql
INSERT INTO user_cars (user_car_id, user_id, trim_id, color_id, vin, year, current_mileage)
VALUES (
  'd0000000-0000-0000-0000-000000000001',
  '00000000-0000-0000-0000-000000000004',
  '80000000-0000-0000-0000-000000000001',
  '90000000-0000-0000-0000-000000000001',
  'JTDBT923503012345', 2021, 42800
) ON CONFLICT (user_car_id) DO NOTHING;
```

## Документы

Метаданные в таблице `documents`, байты — в `DOCUMENT_STORAGE_ROOT` (см. `backend/.env.example`).

## Замечания

- Перед нетривиальным `ALTER` на проде — бэкап.
- После изменения `schema.sql` проверьте модели и репозитории в Go.
