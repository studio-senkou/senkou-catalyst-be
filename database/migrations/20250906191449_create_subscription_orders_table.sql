-- migrate:up
CREATE TABLE IF NOT EXISTS subscription_orders (
    id UUID PRIMARY KEY,
    user_id INT NOT NULL,
    subscription_id INT NOT NULL,
    payment_transaction_id UUID NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
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
            WHERE conname = 'fk_subscription_orders_user'
        ) THEN
            ALTER TABLE subscription_orders
                ADD CONSTRAINT fk_subscription_orders_user
                FOREIGN KEY (user_id) REFERENCES users(id)
                ON DELETE CASCADE;
        END IF;

        -- Verify subscription foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_subscription_orders_subscription'
        ) THEN
            ALTER TABLE subscription_orders
                ADD CONSTRAINT fk_subscription_orders_subscription
                FOREIGN KEY (subscription_id) REFERENCES subscriptions(id)
                ON DELETE CASCADE;
        END IF;

        -- Verify payment transaction foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_subscription_orders_payment_transaction'
        ) THEN
            ALTER TABLE subscription_orders
                ADD CONSTRAINT fk_subscription_orders_payment_transaction
                FOREIGN KEY (payment_transaction_id) REFERENCES payment_transactions(id)
                ON DELETE CASCADE;
        END IF;

        -- Verify user index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_subscription_orders_user_id'
        ) THEN
            CREATE INDEX idx_subscription_orders_user_id ON subscription_orders(user_id);
        END IF;

        -- Verify subscription index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_subscription_orders_subscription_id'
        ) THEN
            CREATE INDEX idx_subscription_orders_subscription_id ON subscription_orders(subscription_id);
        END IF;

        -- Verify payment transaction index is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_subscription_orders_payment_transaction_id'
        ) THEN
            CREATE INDEX idx_subscription_orders_payment_transaction_id ON subscription_orders(payment_transaction_id);
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE subscription_orders
    DROP CONSTRAINT IF EXISTS fk_subscription_orders_payment_transaction;

ALTER TABLE subscription_orders
    DROP CONSTRAINT IF EXISTS fk_subscription_orders_subscription;

ALTER TABLE subscription_orders
    DROP CONSTRAINT IF EXISTS fk_subscription_orders_user;

DROP INDEX IF EXISTS idx_subscription_orders_payment_transaction_id ON subscription_orders;

DROP INDEX IF EXISTS idx_subscription_orders_subscription_id ON subscription_orders;

DROP INDEX IF EXISTS idx_subscription_orders_user_id ON subscription_orders;

DROP TABLE IF EXISTS subscription_orders;
