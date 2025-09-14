-- migrate:up
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    sub_id INT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    payment_status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

DO $$
    BEGIN

        -- Verify user foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_user_subscriptions_user'
        ) THEN
            ALTER TABLE user_subscriptions
                ADD CONSTRAINT fk_user_subscriptions_user
                FOREIGN KEY (user_id) REFERENCES users(id)
                ON DELETE CASCADE;
        END IF;

        -- Verify subscription foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_user_subscriptions_subscription'
        ) THEN
            ALTER TABLE user_subscriptions
                ADD CONSTRAINT fk_user_subscriptions_subscription
                FOREIGN KEY (sub_id) REFERENCES subscriptions(id)
                ON DELETE CASCADE;
        END IF;

        -- Verify user index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_user_subscriptions_user_id'
        ) THEN
            CREATE INDEX idx_user_subscriptions_user_id ON user_subscriptions(user_id);
        END IF;

        -- Verify subscription index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_user_subscriptions_sub_id'
        ) THEN
            CREATE INDEX idx_user_subscriptions_sub_id ON user_subscriptions(sub_id);
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE user_subscriptions
    DROP CONSTRAINT IF EXISTS fk_user_subscriptions_subscription;

ALTER TABLE user_subscriptions
    DROP CONSTRAINT IF EXISTS fk_user_subscriptions_user;

DROP INDEX IF EXISTS idx_user_subscriptions_sub_id;

DROP INDEX IF EXISTS idx_user_subscriptions_user_id;

DROP TABLE IF EXISTS user_subscriptions;
