-- migrate:up
CREATE TABLE IF NOT EXISTS payment_transactions (
    id UUID PRIMARY KEY,
    payment_type VARCHAR(50) NOT NULL,
    payment_channel VARCHAR(50) NOT NULL,
    fraud_status VARCHAR(20),
    amount NUMERIC(15, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL,
    transaction_id VARCHAR(100) UNIQUE,
    transaction_time TIMESTAMP,
    signature_key VARCHAR(255),
    expired_at TIMESTAMP,
    settled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS payment_transactions;
