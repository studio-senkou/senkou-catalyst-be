CREATE TABLE users (
    id TEXT PRIMARY KEY,
    merchant_id TEXT NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE merchants (
    id TEXT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    owner_id TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_owner FOREIGN KEY (owner_id) REFERENCES users(id),
    deleted_at TIMESTAMP
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    image_url VARCHAR(100),
    description VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, description),
    CONSTRAINT fk_user_category FOREIGN KEY (user_id) REFERENCES users(id),
    deleted_at TIMESTAMP
);

CREATE TABLE predefined_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    image_url VARCHAR(100),
    description VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE products (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price INT NOT NULL,
    affiliate_url TEXT NOT NULL,
    category_id BIGINT,
    image_urls TEXT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user_product FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_category_product FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    description TEXT,
    duration SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE subscription_plans (
    id BIGSERIAL PRIMARY KEY,
    sub_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_sub_plan FOREIGN KEY (sub_id) REFERENCES subscriptions(id)
);

CREATE TABLE user_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    sub_id BIGINT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    payment_status VARCHAR(20),
    CONSTRAINT fk_user_subscription FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_sub_subscription FOREIGN KEY (sub_id) REFERENCES subscriptions(id)
);
