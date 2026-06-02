-- Очистка всех прикладных данных CarKeeper.
-- Справочники из schema.sql (роли, права, статусы заказов) сохраняются.
--
-- После сброса: psql -f database/seed.sql
-- Или: ./scripts/db-reseed.sh

BEGIN;

TRUNCATE TABLE
    documents,
    service_appointment_types,
    service_appointments,
    configuration_options,
    orders,
    configurations,
    user_cars,
    news,
    trim_options,
    trims,
    generations,
    models,
    brands,
    service_types,
    branches,
    users,
    colors,
    options,
    engine_types,
    transmissions,
    drive_types
RESTART IDENTITY CASCADE;

COMMIT;
