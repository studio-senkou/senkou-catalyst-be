-- migrate:up
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    merchant_id CHAR(16) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE categories
    ADD CONSTRAINT fk_categories_merchant
    FOREIGN KEY (merchant_id)
    REFERENCES merchants(id)
    ON DELETE CASCADE;

CREATE INDEX idx_categories_merchant_id ON categories(merchant_id);

-- migrate:down
ALTER TABLE categories
    DROP CONSTRAINT IF EXISTS fk_categories_merchant;

ALTER TABLE categories
    DROP INDEX IF EXISTS idx_categories_merchant_id;

DROP TABLE IF EXISTS categories;
