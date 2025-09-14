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

DO $$
    BEGIN
        -- Verify user foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_activation_user'
        ) THEN
            ALTER TABLE email_activation_tokens
                ADD CONSTRAINT fk_activation_user
                FOREIGN KEY (user_id) REFERENCES users(id)
                ON DELETE CASCADE;
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE email_activation_tokens
    DROP CONSTRAINT IF EXISTS fk_activation_user;

DROP TABLE IF EXISTS email_activation_tokens;

