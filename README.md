# CarKeeper

Информационная система автосалона: каталог и конфигуратор, заказы, сервисные записи, документы, новости и админ-управление справочниками.

## Стек

| Часть | Технологии |
|--------|------------|
| Backend | Go, Chi, PostgreSQL, JWT + HttpOnly cookie |
| Frontend | React, Vite, TanStack Query, Tailwind, shadcn/ui |
| БД | `database/schema.sql` — источник правды схемы |

## Быстрый старт (Docker)

Требуется [Docker](https://docs.docker.com/get-docker/) и Docker Compose.

```bash
git clone <repo-url> car-keeper
cd car-keeper
docker compose up --build
```

| Сервис | URL |
|--------|-----|
| Веб-интерфейс | http://localhost:8081 |
| API | http://localhost:8080 |
| Health | http://localhost:8080/health |
| PostgreSQL | `localhost:5432` (user/pass/db: `postgres` / `postgres` / `carkeeper`) |

При **первом** запуске Postgres автоматически применяет `schema.sql`, `catalog_view.sql` и `seed.sql`.

### Демо-аккаунты (пароль у всех: `password123`)

| Роль | Email |
|------|--------|
| Администратор | admin@carkeeper.ru |
| Менеджер | manager@carkeeper.ru |
| Мастер-приёмщик | service_advisor@carkeeper.ru |
| Клиент | customer@carkeeper.ru |
| Клиент | dmitry@carkeeper.ru |

Остановка: `docker compose down`. Полное удаление данных БД: `docker compose down -v`.

### Сброс и перезагрузка демо-данных в Docker

```bash
docker compose exec postgres psql -U postgres -d carkeeper -f /scripts/reset_data.sql
docker compose exec postgres psql -U postgres -d carkeeper -f /scripts/seed.sql
```

Или локально (см. `scripts/db-reseed.ps1` / `scripts/db-reseed.sh`).

## Локальная разработка без Docker

### База данных

```bash
# Новая БД
psql -U postgres -c "CREATE DATABASE carkeeper;"
psql -U postgres -d carkeeper -f database/schema.sql
psql -U postgres -d carkeeper -f database/catalog_view.sql
psql -U postgres -d carkeeper -f database/seed.sql
```

Подробнее: [database/README.md](database/README.md).

### Backend

```bash
cd backend
cp .env.example .env
go run .
```

API: http://localhost:8080

### Frontend

```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

UI: http://localhost:3000 (Vite проксирует `/api` на backend).

## Тесты

```bash
cd backend
go test ./...
```

Integration-тесты требуют PostgreSQL с применёнными `schema.sql` + `seed.sql`.

## Структура репозитория

```
backend/          # Go API
frontend/         # React SPA
database/         # schema.sql, seed.sql, reset_data.sql
scripts/          # db-reseed.sh / .ps1
docker-compose.yml
```

## Production (кратко)

- `ENV=production`, уникальный `JWT_SECRET`, `JWT_COOKIE_SECURE=true`
- `DB_SSLMODE=require`, `CORS_ALLOWED_ORIGINS` = URL фронта
- **Не** запускать `seed.sql` на боевой БД
- Файлы документов: каталог `DOCUMENT_STORAGE_ROOT`, бэкап вместе с БД

См. `backend/.env.example` и `frontend/.env.example`.
