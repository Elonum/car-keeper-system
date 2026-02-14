BEGIN;

-- Extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Trigger function for auto updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Users table
CREATE TABLE users (
    user_id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name   varchar(100) NOT NULL,
    last_name    varchar(100) NOT NULL,
    email        varchar(255) NOT NULL UNIQUE,
    phone        varchar(30) UNIQUE,
    password_hash varchar(255) NOT NULL,
    role         varchar(30) NOT NULL DEFAULT 'customer',
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    CHECK (role IN ('customer','manager','admin'))
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_email ON users(email);

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Branches table (филиалы для сервисного обслуживания)
CREATE TABLE branches (
    branch_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(200) NOT NULL,
    address     text NOT NULL,
    phone       varchar(30),
    email       varchar(255),
    is_active   boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_branches_is_active ON branches(is_active);

CREATE TRIGGER trg_branches_updated_at
BEFORE UPDATE ON branches
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Brands table
CREATE TABLE brands (
    brand_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name     varchar(150) NOT NULL UNIQUE,
    country  varchar(100) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_brands_name ON brands(name);

-- Models table
CREATE TABLE models (
    model_id    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_id    uuid NOT NULL REFERENCES brands(brand_id) ON DELETE RESTRICT,
    name        varchar(150) NOT NULL,
    segment     varchar(100),
    description text,
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (brand_id, name)
);

CREATE INDEX idx_models_brand_id ON models(brand_id);
CREATE INDEX idx_models_name ON models(name);

-- Generations table
CREATE TABLE generations (
    generation_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id      uuid NOT NULL REFERENCES models(model_id) ON DELETE CASCADE,
    name          varchar(150) NOT NULL,
    year_from     integer NOT NULL CHECK (year_from >= 1900 AND year_from <= 2100),
    year_to       integer CHECK (year_to IS NULL OR (year_to >= 1900 AND year_to <= 2100)),
    created_at    timestamptz NOT NULL DEFAULT now(),
    CHECK (year_to IS NULL OR year_to >= year_from),
    UNIQUE (model_id, name, year_from)
);

CREATE INDEX idx_generations_model_id ON generations(model_id);

-- Technical dictionaries
CREATE TABLE engine_types (
    engine_type_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name           varchar(100) NOT NULL UNIQUE,
    created_at     timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE transmissions (
    transmission_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name            varchar(100) NOT NULL UNIQUE,
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE drive_types (
    drive_type_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name          varchar(100) NOT NULL UNIQUE,
    created_at   timestamptz NOT NULL DEFAULT now()
);

-- Trims table
CREATE TABLE trims (
    trim_id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    generation_id    uuid NOT NULL REFERENCES generations(generation_id) ON DELETE CASCADE,
    name             varchar(150) NOT NULL,
    base_price       numeric(12,2) NOT NULL CHECK (base_price >= 0),
    engine_type_id   uuid NOT NULL REFERENCES engine_types(engine_type_id) ON DELETE RESTRICT,
    transmission_id  uuid NOT NULL REFERENCES transmissions(transmission_id) ON DELETE RESTRICT,
    drive_type_id    uuid NOT NULL REFERENCES drive_types(drive_type_id) ON DELETE RESTRICT,
    is_available     boolean NOT NULL DEFAULT true,
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now(),
    UNIQUE (generation_id, name)
);

CREATE INDEX idx_trims_generation_id ON trims(generation_id);
CREATE INDEX idx_trims_engine_type_id ON trims(engine_type_id);
CREATE INDEX idx_trims_transmission_id ON trims(transmission_id);
CREATE INDEX idx_trims_drive_type_id ON trims(drive_type_id);
CREATE INDEX idx_trims_is_available ON trims(is_available);
CREATE INDEX idx_trims_base_price ON trims(base_price);
CREATE INDEX idx_trims_available_price ON trims(is_available, base_price) WHERE is_available = true;

CREATE TRIGGER trg_trims_updated_at
BEFORE UPDATE ON trims
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Colors table
CREATE TABLE colors (
    color_id    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(100) NOT NULL UNIQUE,
    hex_code    varchar(7) CHECK (hex_code IS NULL OR hex_code ~ '^#[0-9A-Fa-f]{6}$'),
    price_delta numeric(12,2) NOT NULL DEFAULT 0 CHECK (price_delta >= 0),
    is_available boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_colors_is_available ON colors(is_available);

-- Options table
CREATE TABLE options (
    option_id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        varchar(150) NOT NULL UNIQUE,
    description text,
    price       numeric(12,2) NOT NULL CHECK (price >= 0),
    is_available boolean NOT NULL DEFAULT true,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_options_is_available ON options(is_available);

-- Trim options junction table
CREATE TABLE trim_options (
    trim_id   uuid NOT NULL REFERENCES trims(trim_id) ON DELETE CASCADE,
    option_id uuid NOT NULL REFERENCES options(option_id) ON DELETE CASCADE,
    PRIMARY KEY (trim_id, option_id)
);

CREATE INDEX idx_trim_options_option_id ON trim_options(option_id);

-- Configurations table
CREATE TABLE configurations (
    configuration_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    trim_id          uuid NOT NULL REFERENCES trims(trim_id) ON DELETE RESTRICT,
    color_id         uuid NOT NULL REFERENCES colors(color_id) ON DELETE RESTRICT,
    status           varchar(30) NOT NULL DEFAULT 'draft',
    total_price      numeric(12,2) NOT NULL CHECK (total_price >= 0),
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now(),
    CHECK (status IN ('draft','confirmed','ordered','cancelled','purchased'))
);

CREATE INDEX idx_configurations_user_id ON configurations(user_id);
CREATE INDEX idx_configurations_trim_id ON configurations(trim_id);
CREATE INDEX idx_configurations_status ON configurations(status);
CREATE INDEX idx_configurations_created_at ON configurations(created_at);

CREATE TRIGGER trg_configurations_updated_at
BEFORE UPDATE ON configurations
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Configuration options junction table
CREATE TABLE configuration_options (
    configuration_id uuid NOT NULL REFERENCES configurations(configuration_id) ON DELETE CASCADE,
    option_id        uuid NOT NULL REFERENCES options(option_id) ON DELETE RESTRICT,
    PRIMARY KEY (configuration_id, option_id)
);

CREATE INDEX idx_configuration_options_option_id ON configuration_options(option_id);

-- Orders table
CREATE TABLE orders (
    order_id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          uuid NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    configuration_id uuid NOT NULL UNIQUE REFERENCES configurations(configuration_id) ON DELETE RESTRICT,
    manager_id       uuid REFERENCES users(user_id) ON DELETE SET NULL,
    status           varchar(30) NOT NULL DEFAULT 'pending',
    final_price      numeric(12,2) NOT NULL CHECK (final_price >= 0),
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now(),
    CHECK (status IN ('pending','approved','paid','completed','cancelled'))
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_configuration_id ON orders(configuration_id);
CREATE INDEX idx_orders_manager_id ON orders(manager_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

CREATE TRIGGER trg_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- User cars table
CREATE TABLE user_cars (
    user_car_id     uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    trim_id         uuid NOT NULL REFERENCES trims(trim_id) ON DELETE RESTRICT,
    color_id        uuid NOT NULL REFERENCES colors(color_id) ON DELETE RESTRICT,
    vin             varchar(17) NOT NULL UNIQUE CHECK (LENGTH(vin) = 17),
    year            integer NOT NULL CHECK (year >= 1900 AND year <= EXTRACT(YEAR FROM CURRENT_DATE) + 1),
    current_mileage integer NOT NULL DEFAULT 0 CHECK (current_mileage >= 0),
    purchase_date   date CHECK (purchase_date IS NULL OR purchase_date <= CURRENT_DATE),
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_cars_user_id ON user_cars(user_id);
CREATE INDEX idx_user_cars_trim_id ON user_cars(trim_id);
CREATE INDEX idx_user_cars_vin ON user_cars(vin);

-- Service types table (типы сервисных услуг)
CREATE TABLE service_types (
    service_type_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name           varchar(150) NOT NULL UNIQUE,
    category       varchar(50) NOT NULL,
    description    text,
    price          numeric(12,2) NOT NULL CHECK (price >= 0),
    duration_minutes integer CHECK (duration_minutes IS NULL OR duration_minutes > 0),
    is_available   boolean NOT NULL DEFAULT true,
    created_at     timestamptz NOT NULL DEFAULT now(),
    CHECK (category IN ('maintenance','repair','diagnostics','detailing','tires'))
);

CREATE INDEX idx_service_types_category ON service_types(category);
CREATE INDEX idx_service_types_is_available ON service_types(is_available);

-- Service appointments table
CREATE TABLE service_appointments (
    service_appointment_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_car_id            uuid NOT NULL REFERENCES user_cars(user_car_id) ON DELETE CASCADE,
    branch_id              uuid NOT NULL REFERENCES branches(branch_id) ON DELETE RESTRICT,
    manager_id             uuid REFERENCES users(user_id) ON DELETE SET NULL,
    appointment_date       timestamptz NOT NULL,
    status                 varchar(30) NOT NULL DEFAULT 'scheduled',
    description            text,
    created_at             timestamptz NOT NULL DEFAULT now(),
    updated_at             timestamptz NOT NULL DEFAULT now(),
    CHECK (status IN ('scheduled','completed','cancelled'))
);

CREATE INDEX idx_service_user_car_id ON service_appointments(user_car_id);
CREATE INDEX idx_service_branch_id ON service_appointments(branch_id);
CREATE INDEX idx_service_manager_id ON service_appointments(manager_id);
CREATE INDEX idx_service_status ON service_appointments(status);
CREATE INDEX idx_service_appointment_date ON service_appointments(appointment_date);

-- Service appointment types junction table (связь записей на ТО с типами услуг)
CREATE TABLE service_appointment_types (
    service_appointment_id uuid NOT NULL REFERENCES service_appointments(service_appointment_id) ON DELETE CASCADE,
    service_type_id        uuid NOT NULL REFERENCES service_types(service_type_id) ON DELETE RESTRICT,
    PRIMARY KEY (service_appointment_id, service_type_id)
);

CREATE INDEX idx_service_appointment_types_service_type_id ON service_appointment_types(service_type_id);

CREATE TRIGGER trg_service_appointments_updated_at
BEFORE UPDATE ON service_appointments
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Documents table
CREATE TABLE documents (
    document_id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    order_id            uuid REFERENCES orders(order_id) ON DELETE CASCADE,
    service_appointment_id uuid REFERENCES service_appointments(service_appointment_id) ON DELETE CASCADE,
    document_type       varchar(50) NOT NULL,
    file_path           text NOT NULL,
    file_name           varchar(255),
    file_size           bigint CHECK (file_size IS NULL OR file_size > 0),
    mime_type           varchar(100),
    created_at          timestamptz NOT NULL DEFAULT now(),
    CHECK (document_type IN ('commercial_offer','order_contract','service_order','service_act')),
    CHECK (order_id IS NOT NULL OR service_appointment_id IS NOT NULL)
);

CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_documents_order_id ON documents(order_id);
CREATE INDEX idx_documents_service_appointment_id ON documents(service_appointment_id);
CREATE INDEX idx_documents_document_type ON documents(document_type);

-- News table
CREATE TABLE news (
    news_id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title        varchar(255) NOT NULL,
    content      text NOT NULL,
    author_id    uuid REFERENCES users(user_id) ON DELETE SET NULL,
    published_at timestamptz,
    is_published boolean NOT NULL DEFAULT false,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_news_published_at ON news(published_at);
CREATE INDEX idx_news_is_published ON news(is_published);
CREATE INDEX idx_news_author_id ON news(author_id);
CREATE INDEX idx_news_published ON news(is_published, published_at DESC) WHERE is_published = true;

CREATE TRIGGER trg_news_updated_at
BEFORE UPDATE ON news
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

COMMIT;

