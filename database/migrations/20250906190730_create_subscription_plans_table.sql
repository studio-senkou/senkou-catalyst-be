-- migrate:up
CREATE TABLE IF NOT EXISTS subscription_plans (
    id SERIAL PRIMARY KEY,
    sub_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

DO $$
    BEGIN

        -- Verify subscription foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_subscription_plans_subscription'
        ) THEN
            ALTER TABLE subscription_plans
                ADD CONSTRAINT fk_subscription_plans_subscription
                FOREIGN KEY (sub_id) REFERENCES subscriptions(id)
                ON DELETE CASCADE;
        END IF;

        -- Verify subscription index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_subscription_plans_sub_id'
        ) THEN
            CREATE INDEX idx_subscription_plans_sub_id ON subscription_plans(sub_id);
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE subscription_plans
    DROP CONSTRAINT IF EXISTS fk_subscription_plans_subscription;

DROP INDEX IF EXISTS idx_subscription_plans_sub_id;

DROP TABLE IF EXISTS subscription_plans;