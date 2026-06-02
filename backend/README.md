# CarKeeper — Backend

REST API на Go (Chi), PostgreSQL, JWT + HttpOnly cookie `carkeeper_session`.

## Требования

- Go 1.22+
- PostgreSQL 16+

## Настройка

```bash
cp .env.example .env
```

Запуск:

```bash
go run .
```

Health: http://localhost:8080/health

## База данных

Перед первым запуском:

```bash
psql -U postgres -d carkeeper -f ../database/schema.sql
psql -U postgres -d carkeeper -f ../database/catalog_view.sql
psql -U postgres -d carkeeper -f ../database/seed.sql
```

Сброс демо-данных: `../database/reset_data.sql` + `seed.sql` или `../scripts/db-reseed.sh`.

## Тесты

```bash
go test ./...
go test ./internal/integration/... -v   # нужна БД с seed
```

## Docker

```bash
docker compose up --build api postgres
```

API: http://localhost:8080

## Документы

Файлы хранятся в `DOCUMENT_STORAGE_ROOT` (по умолчанию `./data/documents`). В БД — только ключ в `documents.file_path`.

## Полезные команды

```bash
go run ./cmd/generate-password   # хеш пароля для ручного INSERT
```
