-- Discount Codes Migration
-- Adds tables for unique discount codes and redemption tracking

-- Discount codes table
CREATE TABLE IF NOT EXISTS discount_codes (
    id TEXT PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    event_id TEXT NOT NULL,
    brand_id TEXT NOT NULL,
    discount_pct REAL NOT NULL CHECK (discount_pct > 0 AND discount_pct <= 100),
    discount_amount REAL DEFAULT 0,
    max_uses INTEGER DEFAULT 0, -- 0 = unlimited
    used_count INTEGER DEFAULT 0,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Code redemptions table
CREATE TABLE IF NOT EXISTS code_redemptions (
    id TEXT PRIMARY KEY,
    code_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    order_id TEXT,
    amount REAL NOT NULL,
    discount_applied REAL NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (code_id) REFERENCES discount_codes(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(code_id, order_id)
);

-- Code analytics table for tracking performance
CREATE TABLE IF NOT EXISTS code_analytics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code_id TEXT NOT NULL,
    event_type TEXT NOT NULL CHECK (event_type IN ('view', 'copy', 'attempt', 'success', 'failure')),
    user_id TEXT,
    metadata TEXT, -- JSON data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (code_id) REFERENCES discount_codes(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Attribution tracking for discount codes
CREATE TABLE IF NOT EXISTS discount_attribution (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    redemption_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    content_id TEXT,
    attribution_source TEXT NOT NULL CHECK (attribution_source IN ('event', 'content', 'social', 'direct')),
    attribution_weight REAL DEFAULT 1.0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (redemption_id) REFERENCES code_redemptions(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_discount_codes_event ON discount_codes(event_id);
CREATE INDEX IF NOT EXISTS idx_discount_codes_brand ON discount_codes(brand_id);
CREATE INDEX IF NOT EXISTS idx_discount_codes_code ON discount_codes(code);
CREATE INDEX IF NOT EXISTS idx_discount_codes_expires ON discount_codes(expires_at);
CREATE INDEX IF NOT EXISTS idx_code_redemptions_code ON code_redemptions(code_id);
CREATE INDEX IF NOT EXISTS idx_code_redemptions_user ON code_redemptions(user_id);
CREATE INDEX IF NOT EXISTS idx_code_redemptions_created ON code_redemptions(created_at);
CREATE INDEX IF NOT EXISTS idx_code_analytics_code ON code_analytics(code_id);
CREATE INDEX IF NOT EXISTS idx_code_analytics_event ON code_analytics(event_type);
CREATE INDEX IF NOT EXISTS idx_discount_attribution_redemption ON discount_attribution(redemption_id);

-- Insert sample discount codes
INSERT OR IGNORE INTO discount_codes (id, code, event_id, brand_id, discount_pct, max_uses, expires_at) VALUES
('discount_1', 'EVENT20', 'event_1', 'brand_1', 20.0, 100, datetime('now', '+30 days')),
('discount_2', 'SPECIAL15', 'event_1', 'brand_1', 15.0, 50, datetime('now', '+15 days')),
('discount_3', 'WELCOME10', 'event_1', 'brand_1', 10.0, 0, datetime('now', '+60 days'));

-- Insert sample redemptions
INSERT OR IGNORE INTO code_redemptions (id, code_id, user_id, order_id, amount, discount_applied) VALUES
('redemption_1', 'discount_1', 'user_1', 'order_123', 99.99, 19.99),
('redemption_2', 'discount_1', 'user_2', 'order_124', 149.99, 29.99),
('redemption_3', 'discount_2', 'user_3', 'order_125', 79.99, 11.99);

-- Insert sample analytics
INSERT OR IGNORE INTO code_analytics (code_id, event_type, user_id) VALUES
('discount_1', 'view', 'user_1'),
('discount_1', 'copy', 'user_1'),
('discount_1', 'success', 'user_1'),
('discount_2', 'view', 'user_2'),
('discount_2', 'attempt', 'user_2'),
('discount_2', 'success', 'user_2');