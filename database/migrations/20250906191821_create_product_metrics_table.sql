-- migrate:up
CREATE TABLE IF NOT EXISTS product_metrics (
    id SERIAL PRIMARY KEY,
    product_id UUID NOT NULL,
    origin VARCHAR(20) NOT NULL,
    ua_browser TEXT,
    ua_os TEXT,
    interaction TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
    BEGIN

        -- Verify product foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_product'
        ) THEN
            ALTER TABLE product_metrics
                ADD CONSTRAINT fk_product
                FOREIGN KEY (product_id) REFERENCES products(id)
                ON DELETE CASCADE;
        END IF;

        -- Verify product index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_product_metrics_product_id'
        ) THEN
            CREATE INDEX idx_product_metrics_product_id ON product_metrics(product_id);
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE product_metrics
    DROP CONSTRAINT IF EXISTS fk_product;

DROP INDEX IF EXISTS idx_product_metrics_product_id;

DROP TABLE IF EXISTS product_metrics;
