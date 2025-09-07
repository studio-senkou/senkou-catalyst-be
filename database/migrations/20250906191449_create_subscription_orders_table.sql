-- migrate:up
CREATE TABLE subscription_orders (
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

ALTER TABLE subscription_orders
    ADD CONSTRAINT fk_subscription_orders_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE;

ALTER TABLE subscription_orders
    ADD CONSTRAINT fk_subscription_orders_subscription
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id)
    ON DELETE CASCADE;

ALTER TABLE subscription_orders
    ADD CONSTRAINT fk_subscription_orders_payment_transaction
    FOREIGN KEY (payment_transaction_id) REFERENCES payment_transactions(id)
    ON DELETE CASCADE;

CREATE INDEX idx_subscription_orders_user_id ON subscription_orders(user_id);

CREATE INDEX idx_subscription_orders_subscription_id ON subscription_orders(subscription_id);

CREATE INDEX idx_subscription_orders_payment_transaction_id ON subscription_orders(payment_transaction_id);

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
