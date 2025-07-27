-- E-commerce integration tables

-- E-commerce integrations table
CREATE TABLE IF NOT EXISTS ecommerce_integrations (
    id TEXT PRIMARY KEY,
    brand_id TEXT NOT NULL,
    platform_type TEXT NOT NULL, -- 'shopify', 'woocommerce', 'magento', etc.
    api_key TEXT NOT NULL,
    store_url TEXT NOT NULL,
    webhook_url TEXT NOT NULL,
    status TEXT DEFAULT 'active', -- 'active', 'inactive', 'error'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Purchase tracking table
CREATE TABLE IF NOT EXISTS purchases (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL,
    user_id TEXT,
    event_id TEXT,
    brand_id TEXT,
    integration_id TEXT,
    amount DECIMAL(10,2) NOT NULL,
    currency TEXT DEFAULT 'USD',
    products TEXT, -- JSON array of products
    status TEXT DEFAULT 'completed', -- 'pending', 'completed', 'refunded'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (integration_id) REFERENCES ecommerce_integrations(id)
);

-- Purchase products table (normalized)
CREATE TABLE IF NOT EXISTS purchase_products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    product_name TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    quantity INTEGER DEFAULT 1,
    FOREIGN KEY (purchase_id) REFERENCES purchases(id)
);

-- Webhook logs table
CREATE TABLE IF NOT EXISTS webhook_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    integration_id TEXT NOT NULL,
    webhook_type TEXT NOT NULL, -- 'order_created', 'order_updated', etc.
    payload TEXT NOT NULL, -- JSON payload
    status TEXT DEFAULT 'received', -- 'received', 'processed', 'failed'
    processed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (integration_id) REFERENCES ecommerce_integrations(id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_ecommerce_integrations_brand_id ON ecommerce_integrations(brand_id);
CREATE INDEX IF NOT EXISTS idx_purchases_event_id ON purchases(event_id);
CREATE INDEX IF NOT EXISTS idx_purchases_user_id ON purchases(user_id);
CREATE INDEX IF NOT EXISTS idx_purchases_brand_id ON purchases(brand_id);
CREATE INDEX IF NOT EXISTS idx_purchase_products_purchase_id ON purchase_products(purchase_id);
CREATE INDEX IF NOT EXISTS idx_webhook_logs_integration_id ON webhook_logs(integration_id);