CREATE TABLE IF NOT EXISTS categories (
    id         SERIAL PRIMARY KEY,
    code       VARCHAR(32) UNIQUE NOT NULL,
    name       VARCHAR(256)       NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE products
    ADD COLUMN category_id INTEGER NULL REFERENCES categories (id) ON DELETE CASCADE;

CREATE INDEX idx_product_category_id ON products (category_id);