-- migrate:up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_oauth BOOLEAN DEFAULT FALSE;

-- migrate:down
ALTER TABLE users
    DROP COLUMN IF EXISTS is_oauth;
