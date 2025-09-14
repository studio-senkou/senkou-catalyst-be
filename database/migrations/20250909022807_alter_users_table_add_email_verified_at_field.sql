-- migrate:up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMP DEFAULT NULL;

-- migrate:down
ALTER TABLE users
    DROP COLUMN email_verified_at;
