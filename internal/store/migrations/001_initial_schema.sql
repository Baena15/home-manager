-- ─── Initial schema for Home Manager ───────────────────────────────

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ─── users ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'partner',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- ─── products ──────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    unit VARCHAR(50) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);

-- ─── product_prices ────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS product_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    store VARCHAR(100) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_product_prices_product_id ON product_prices(product_id);
CREATE INDEX IF NOT EXISTS idx_product_prices_recorded_at ON product_prices(recorded_at);

-- ─── shopping_lists ────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS shopping_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_shopping_lists_status CHECK (status IN ('active', 'completed'))
);

CREATE INDEX IF NOT EXISTS idx_shopping_lists_status ON shopping_lists(status);
CREATE INDEX IF NOT EXISTS idx_shopping_lists_created_by ON shopping_lists(created_by);

-- ─── shopping_list_items ───────────────────────────────────────────
CREATE TABLE IF NOT EXISTS shopping_list_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID NOT NULL REFERENCES shopping_lists(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity DECIMAL(10, 3) NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    purchased BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_shopping_list_items_list_id ON shopping_list_items(list_id);
CREATE INDEX IF NOT EXISTS idx_shopping_list_items_product_id ON shopping_list_items(product_id);

-- ─── down migration ────────────────────────────────────────────────
-- DROP TABLE IF EXISTS shopping_list_items;
-- DROP TABLE IF EXISTS shopping_lists;
-- DROP TABLE IF EXISTS product_prices;
-- DROP TABLE IF EXISTS products;
-- DROP TABLE IF EXISTS users;
