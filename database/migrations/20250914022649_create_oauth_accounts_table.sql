-- migrate:up
CREATE TABLE IF NOT EXISTS oauth_accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    provider VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    token_expiry TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
    BEGIN

        -- Verify user foreign key constraint is not exists
        -- If already exists, skip the migration to avoid errors
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_user_oauth_accounts'
        ) THEN
            ALTER TABLE oauth_accounts
                ADD CONSTRAINT fk_user_oauth_accounts
                FOREIGN KEY (user_id) REFERENCES users(id)
                ON DELETE CASCADE;
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE oauth_accounts
    DROP CONSTRAINT IF EXISTS fk_user_oauth_accounts;

DROP TABLE IF EXISTS oauth_accounts;