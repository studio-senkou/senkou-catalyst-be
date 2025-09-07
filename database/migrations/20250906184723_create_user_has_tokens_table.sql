-- migrate:up
CREATE TABLE users_has_token (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE users_has_token
    ADD CONSTRAINT fk_users_has_token
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE;

-- migrate:down
ALTER TABLE users_has_token
    DROP CONSTRAINT IF EXISTS fk_users_has_token;

DROP TABLE IF EXISTS users_has_token;
