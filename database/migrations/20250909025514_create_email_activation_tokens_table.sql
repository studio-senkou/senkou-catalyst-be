-- migrate:up
CREATE TABLE IF NOT EXISTS email_activation_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE email_activation_tokens
    ADD CONSTRAINT fk_activation_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE;

-- migrate:down
ALTER TABLE email_activation_tokens
    DROP CONSTRAINT IF EXISTS fk_activation_user;

DROP TABLE IF EXISTS email_activation_tokens;

