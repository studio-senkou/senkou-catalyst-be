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

ALTER TABLE user_subscriptions
    ADD CONSTRAINT fk_user_subscriptions_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE;

ALTER TABLE user_subscriptions
    ADD CONSTRAINT fk_user_subscriptions_subscription
    FOREIGN KEY (sub_id) REFERENCES subscriptions(id)
    ON DELETE CASCADE;

CREATE INDEX idx_user_subscriptions_user_id ON user_subscriptions(user_id);

CREATE INDEX idx_user_subscriptions_sub_id ON user_subscriptions(sub_id);

-- migrate:down
ALTER TABLE user_subscriptions
    DROP CONSTRAINT IF EXISTS fk_user_subscriptions_subscription;

ALTER TABLE user_subscriptions
    DROP CONSTRAINT IF EXISTS fk_user_subscriptions_user;

DROP INDEX IF EXISTS idx_user_subscriptions_sub_id;

DROP INDEX IF EXISTS idx_user_subscriptions_user_id;

DROP TABLE IF EXISTS user_subscriptions;
