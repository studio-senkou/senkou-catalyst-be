-- migrate:up
CREATE TABLE products (
    id UUID PRIMARY KEY,
    merchant_id CHAR(16) NOT NULL,
    category_id INT,
    title VARCHAR(150) NOT NULL,
    price VARCHAR(30) NOT NULL,
    description TEXT,
    affiliate_url TEXT NOT NULL,
    photos JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE products
    ADD CONSTRAINT fk_products_merchant
    FOREIGN KEY (merchant_id) REFERENCES merchants(id)
    ON DELETE CASCADE;

ALTER TABLE products
    ADD CONSTRAINT fk_products_category
    FOREIGN KEY (category_id) REFERENCES categories(id)
    ON DELETE SET NULL;

CREATE INDEX idx_products_merchant_id ON products(merchant_id);

CREATE INDEX idx_products_category_id ON products(category_id);

-- migrate:down
ALTER TABLE products
    DROP CONSTRAINT IF EXISTS fk_products_merchant;

ALTER TABLE products
    DROP CONSTRAINT IF EXISTS fk_products_category;

ALTER TABLE products
    DROP INDEX IF EXISTS idx_products_merchant_id;

ALTER TABLE products
    DROP INDEX IF EXISTS idx_products_category_id;

DROP TABLE IF EXISTS products;

