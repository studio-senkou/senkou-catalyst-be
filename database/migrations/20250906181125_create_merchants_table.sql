-- migrate:up
CREATE TABLE IF NOT EXISTS merchants (
    id CHAR(16) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    owner_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

DO $$
    BEGIN 

        -- Verify owner foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_merchant_owner'
        ) THEN
            ALTER TABLE merchants 
                ADD CONSTRAINT fk_merchant_owner
                    FOREIGN KEY (owner_id)
                    REFERENCES users(id)
                    ON DELETE SET NULL;
        END IF;

        -- Verify owner index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_merchants_owner_id'
        ) THEN
            CREATE INDEX idx_merchants_owner_id ON merchants(owner_id);
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE merchants 
    DROP CONSTRAINT IF EXISTS fk_merchant_owner;

DROP INDEX IF EXISTS idx_merchants_owner_id;

DROP TABLE IF EXISTS merchants;
