-- E-commerce Integration Migration
-- Adds tables for e-commerce platform integrations and purchase tracking

-- E-commerce integrations table
CREATE TABLE IF NOT EXISTS ecommerce_integrations (
    id TEXT PRIMARY KEY,
    brand_id TEXT NOT NULL,
    platform_type TEXT NOT NULL CHECK (platform_type IN ('shopify', 'woocommerce', 'magento', 'bigcommerce')),
    api_key TEXT NOT NULL,
    api_secret TEXT,
    store_url TEXT NOT NULL,
    webhook_url TEXT,
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'error')),
    last_sync DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id),
    UNIQUE(brand_id, platform_type)
);

-- Products table for synced products
CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY,
    brand_id TEXT NOT NULL,
    integration_id TEXT NOT NULL,
    external_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL,
    image_url TEXT,
    category TEXT,
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id),
    FOREIGN KEY (integration_id) REFERENCES ecommerce_integrations(id),
    UNIQUE(integration_id, external_id)
);

-- Purchases table for tracking purchases
CREATE TABLE IF NOT EXISTS purchases (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    event_id TEXT,
    external_order_id TEXT,
    amount REAL NOT NULL,
    currency TEXT DEFAULT 'USD',
    status TEXT DEFAULT 'completed' CHECK (status IN ('pending', 'completed', 'cancelled', 'refunded')),
    attribution_data TEXT, -- JSON with attribution details
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Purchase attribution table for detailed tracking
CREATE TABLE IF NOT EXISTS purchase_attribution (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    attribution_type TEXT NOT NULL CHECK (attribution_type IN ('direct', 'event', 'content', 'discount_code')),
    source_id TEXT, -- content_id, discount_code_id, etc.
    attribution_value REAL DEFAULT 1.0, -- attribution weight
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (purchase_id) REFERENCES purchases(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Webhook events table for tracking webhook deliveries
CREATE TABLE IF NOT EXISTS webhook_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    integration_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload TEXT NOT NULL, -- JSON payload
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'processed', 'failed')),
    retry_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,
    FOREIGN KEY (integration_id) REFERENCES ecommerce_integrations(id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_ecommerce_integrations_brand ON ecommerce_integrations(brand_id);
CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand_id);
CREATE INDEX IF NOT EXISTS idx_products_integration ON products(integration_id);
CREATE INDEX IF NOT EXISTS idx_purchases_user ON purchases(user_id);
CREATE INDEX IF NOT EXISTS idx_purchases_event ON purchases(event_id);
CREATE INDEX IF NOT EXISTS idx_purchases_created ON purchases(created_at);
CREATE INDEX IF NOT EXISTS idx_purchase_attribution_purchase ON purchase_attribution(purchase_id);
CREATE INDEX IF NOT EXISTS idx_purchase_attribution_event ON purchase_attribution(event_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_integration ON webhook_events(integration_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_status ON webhook_events(status);

-- Insert sample data
INSERT OR IGNORE INTO ecommerce_integrations (id, brand_id, platform_type, api_key, store_url, webhook_url) VALUES
('integration_1', 'brand_1', 'shopify', 'sk_test_123', 'https://test-store.myshopify.com', 'https://api.lynkr.com/webhooks/integration_1');

INSERT OR IGNORE INTO products (id, brand_id, integration_id, external_id, name, description, price, image_url) VALUES
('prod_1', 'brand_1', 'integration_1', 'shopify_123', 'Event Special T-Shirt', 'Limited edition t-shirt from the event', 29.99, 'https://example.com/tshirt.jpg'),
('prod_2', 'brand_1', 'integration_1', 'shopify_124', 'Tech Gadget', 'Latest tech gadget showcased at the event', 199.99, 'https://example.com/gadget.jpg');

INSERT OR IGNORE INTO purchases (id, user_id, product_id, event_id, amount, status) VALUES
('purchase_1', 'user_1', 'prod_1', 'event_1', 29.99, 'completed'),
('purchase_2', 'user_2', 'prod_2', 'event_1', 199.99, 'completed');

INSERT OR IGNORE INTO purchase_attribution (purchase_id, event_id, user_id, attribution_type, attribution_value) VALUES
('purchase_1', 'event_1', 'user_1', 'event', 1.0),
('purchase_2', 'event_1', 'user_2', 'event', 1.0);