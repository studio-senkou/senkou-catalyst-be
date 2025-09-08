-- migrate:up
CREATE TABLE IF NOT EXISTS product_metrics (
    id SERIAL PRIMARY KEY,
    product_id UUID NOT NULL,
    origin VARCHAR(20) NOT NULL,
    ua_browser TEXT,
    ua_os TEXT,
    interaction TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE product_metrics
    ADD CONSTRAINT fk_product
    FOREIGN KEY (product_id) REFERENCES products(id)
    ON DELETE CASCADE;

CREATE INDEX idx_product_metrics_product_id ON product_metrics(product_id);

-- migrate:down
ALTER TABLE product_metrics
    DROP CONSTRAINT IF EXISTS fk_product;

DROP INDEX IF EXISTS idx_product_metrics_product_id;

DROP TABLE IF EXISTS product_metrics;
