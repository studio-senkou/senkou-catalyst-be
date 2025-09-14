-- migrate:up
CREATE TABLE IF NOT EXISTS user_has_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token VARCHAR(255) NOT NULL,
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
            WHERE conname = 'fk_user_has_tokens'
        ) THEN
            ALTER TABLE user_has_tokens 
                ADD CONSTRAINT fk_user_has_tokens
                    FOREIGN KEY (user_id)
                    REFERENCES users(id)
                    ON DELETE CASCADE;
        END IF;
    END;
$$;

-- migrate:down
ALTER TABLE user_has_tokens
    DROP CONSTRAINT IF EXISTS fk_user_has_tokens;

DROP TABLE IF EXISTS user_has_tokens;
