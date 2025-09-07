-- migrate:up
CREATE TABLE subscription_plans (
    id SERIAL PRIMARY KEY,
    sub_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE subscription_plans
    ADD CONSTRAINT fk_subscription_plans_subscription
    FOREIGN KEY (sub_id) REFERENCES subscriptions(id)
    ON DELETE CASCADE;

CREATE INDEX idx_subscription_plans_sub_id ON subscription_plans(sub_id);

-- migrate:down
ALTER TABLE subscription_plans
    DROP CONSTRAINT IF EXISTS fk_subscription_plans_subscription;

ALTER TABLE subscription_plans
    DROP INDEX IF EXISTS idx_subscription_plans_sub_id;

DROP TABLE IF EXISTS subscription_plans;