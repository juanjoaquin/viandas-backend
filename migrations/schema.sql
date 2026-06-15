CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(100)  NOT NULL,
    email         VARCHAR(255)  NOT NULL UNIQUE,
    password_hash TEXT          NOT NULL,
    role          VARCHAR(20)   NOT NULL CHECK (role IN ('ADMIN', 'EMPLOYEE')),
    active        BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP     NOT NULL DEFAULT NOW()
);

CREATE TABLE customers (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    type       VARCHAR(20)  NOT NULL CHECK (type IN ('COMPANY', 'PERSON')),
    phone      VARCHAR(100),
    address    TEXT,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE deliveries (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    phone      VARCHAR(100),
    active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE menu_types (
    id         UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100)  NOT NULL UNIQUE,
    price      NUMERIC(10,2),
    active     BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP     NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP     NOT NULL DEFAULT NOW()
);

CREATE TABLE dishes (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    menu_type_id UUID         NOT NULL REFERENCES menu_types(id),
    active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE extra_products (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    category   VARCHAR(30)  NOT NULL CHECK (category IN ('SALAD', 'SANDWICH')),
    active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE week_menus (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    week_start_date DATE      NOT NULL UNIQUE,
    week_end_date   DATE      NOT NULL,
    created_by      UUID      NOT NULL REFERENCES users(id),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE week_menu_items (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    week_menu_id UUID NOT NULL REFERENCES week_menus(id) ON DELETE CASCADE,
    menu_date    DATE NOT NULL,
    menu_type_id UUID NOT NULL REFERENCES menu_types(id),
    dish_id      UUID NOT NULL REFERENCES dishes(id),
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (week_menu_id, menu_date, menu_type_id)
);

CREATE TABLE daily_productions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    production_date  DATE        NOT NULL,
    customer_id      UUID        NOT NULL REFERENCES customers(id),
    fulfillment_type VARCHAR(20) NOT NULL DEFAULT 'PENDING'
    CHECK (fulfillment_type IN ('PENDING', 'DELIVERY', 'PICKUP')),
    delivery_id      UUID        REFERENCES deliveries(id),
    notes            TEXT,
    created_by       UUID        NOT NULL REFERENCES users(id),
    created_at       TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP   NOT NULL DEFAULT NOW(),
    UNIQUE (production_date, customer_id)
);

CREATE TABLE daily_production_lines (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    daily_production_id UUID    NOT NULL REFERENCES daily_productions(id) ON DELETE CASCADE,
    menu_type_id        UUID    NOT NULL REFERENCES menu_types(id),
    quantity            INTEGER NOT NULL DEFAULT 0,
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (daily_production_id, menu_type_id)
);

CREATE TABLE daily_production_extras (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    daily_production_id UUID     NOT NULL REFERENCES daily_productions(id) ON DELETE CASCADE,
    extra_product_id    UUID     NOT NULL REFERENCES extra_products(id),
    quantity            INTEGER  NOT NULL CHECK (quantity > 0),
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW()
);

