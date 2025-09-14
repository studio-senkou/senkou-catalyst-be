-- migrate:up
CREATE TABLE IF NOT EXISTS products (
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

DO $$
    BEGIN

        -- Verify merchant foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_products_merchant'
        ) THEN
            ALTER TABLE products
                ADD CONSTRAINT fk_products_merchant
                FOREIGN KEY (merchant_id) REFERENCES merchants(id)
                ON DELETE CASCADE;
        END IF;

        -- Verify category foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_products_category'
        ) THEN
            ALTER TABLE products
                ADD CONSTRAINT fk_products_category
                FOREIGN KEY (category_id) REFERENCES categories(id)
                ON DELETE SET NULL;
        END IF;

        -- Verify merchant index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_products_merchant_id'
        ) THEN
            CREATE INDEX idx_products_merchant_id ON products(merchant_id);
        END IF;

        -- Verify category index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_products_category_id'
        ) THEN
            CREATE INDEX idx_products_category_id ON products(category_id);
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE products
    DROP CONSTRAINT IF EXISTS fk_products_merchant;

ALTER TABLE products
    DROP CONSTRAINT IF EXISTS fk_products_category;

DROP INDEX IF EXISTS idx_products_merchant_id;
DROP INDEX IF EXISTS idx_products_category_id;

DROP TABLE IF EXISTS products;

