-- migrate:up
CREATE TABLE IF NOT EXISTS user_has_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE user_has_tokens
    ADD CONSTRAINT fk_user_has_tokens
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE;

-- migrate:down
ALTER TABLE user_has_tokens
    DROP CONSTRAINT IF EXISTS fk_user_has_tokens;

DROP TABLE IF EXISTS user_has_tokens;
