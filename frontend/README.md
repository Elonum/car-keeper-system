# CarKeeper — Frontend

SPA для каталога, конфигуратора, личного кабинета и админ-управления.

## Требования

- Node.js 18+
- npm

## Установка

```bash
npm install
cp .env.example .env.local
```

## Переменные окружения

| Переменная | Описание |
|------------|----------|
| `VITE_API_BASE_URL` | Базовый URL API. Локально с Vite: `/api` (прокси). Прямой доступ: `http://localhost:8080/api` |
| `VITE_API_ORIGIN` | Опционально: origin для картинок каталога, если API на другом хосте |

## Запуск

```bash
npm run dev
```

http://localhost:3000 — запросы `/api` проксируются на backend (см. `vite.config.js`).

Сборка:

```bash
npm run build
npm run preview
```

## Docker

Собирается образ из `Dockerfile`, в compose доступен на http://localhost:8081 с прокси `/api` → сервис `api`.

```bash
docker compose up --build web api postgres
```

## Демо-вход

Пароль seed-пользователей: `password123` (см. корневой README).

## Стек

React 18, Vite, React Router, TanStack Query, Axios, Tailwind CSS, shadcn/ui, Lucide.

Ошибки отображаются через `ErrorNotice` (без toast-уведомлений).
