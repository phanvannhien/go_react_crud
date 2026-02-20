-- migrate:up
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    stock INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes
CREATE INDEX idx_products_name ON products (name);
CREATE INDEX idx_products_price ON products (price);
CREATE INDEX idx_products_category_id ON products (category_id);
CREATE INDEX idx_products_is_active ON products (is_active);
CREATE INDEX idx_products_stock ON products (stock);
CREATE INDEX idx_products_created_at ON products (created_at);

-- Composite index for common filter combination
CREATE INDEX idx_products_category_active ON products (category_id, is_active);

-- migrate:down
DROP TABLE IF EXISTS products;
