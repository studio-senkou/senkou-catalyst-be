-- migrate:up
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_phone_key;

-- migrate:down
ALTER TABLE users
    ADD CONSTRAINT users_phone_key UNIQUE (phone);
