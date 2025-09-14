-- migrate:up
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    merchant_id CHAR(16) NOT NULL,
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
            WHERE conname = 'fk_categories_merchant'
        ) THEN
            ALTER TABLE categories
                ADD CONSTRAINT fk_categories_merchant
                FOREIGN KEY (merchant_id)
                REFERENCES merchants(id)
                ON DELETE CASCADE;
        END IF;

        -- Verify merchant index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_categories_merchant_id'
        ) THEN
            CREATE INDEX idx_categories_merchant_id ON categories(merchant_id);
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE categories
    DROP CONSTRAINT IF EXISTS fk_categories_merchant;

DROP INDEX IF EXISTS idx_categories_merchant_id;

DROP TABLE IF EXISTS categories;
