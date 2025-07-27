-- Advanced Analytics Migration
-- Adds tables for AI tagging results and conversion tracking

-- AI tagging results table
CREATE TABLE IF NOT EXISTS ai_tagging_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    products TEXT NOT NULL, -- JSON array of detected products
    tags TEXT NOT NULL, -- JSON array of generated tags
    confidence_score REAL DEFAULT 0,
    processed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Product detection table for detailed tracking
CREATE TABLE IF NOT EXISTS product_detections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    brand_id TEXT NOT NULL,
    confidence REAL NOT NULL,
    bounding_box TEXT, -- JSON with x, y, width, height
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Conversion tracking table
CREATE TABLE IF NOT EXISTS conversion_tracking (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    stage TEXT NOT NULL CHECK (stage IN ('attendance', 'content_view', 'engagement', 'website_visit', 'purchase')),
    metadata TEXT, -- JSON with additional data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- ROI calculations table
CREATE TABLE IF NOT EXISTS roi_calculations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id TEXT NOT NULL,
    brand_id TEXT NOT NULL,
    total_investment REAL NOT NULL,
    total_revenue REAL NOT NULL,
    roi_percentage REAL NOT NULL,
    calculation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Attribution weights table for multi-touch attribution
CREATE TABLE IF NOT EXISTS attribution_weights (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id TEXT NOT NULL,
    attribution_model TEXT NOT NULL CHECK (attribution_model IN ('first_touch', 'last_touch', 'linear', 'time_decay')),
    stage TEXT NOT NULL,
    weight REAL NOT NULL CHECK (weight >= 0 AND weight <= 1),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_ai_tagging_content ON ai_tagging_results(content_id);
CREATE INDEX IF NOT EXISTS idx_ai_tagging_processed ON ai_tagging_results(processed_at);
CREATE INDEX IF NOT EXISTS idx_product_detections_content ON product_detections(content_id);
CREATE INDEX IF NOT EXISTS idx_product_detections_product ON product_detections(product_id);
CREATE INDEX IF NOT EXISTS idx_product_detections_brand ON product_detections(brand_id);
CREATE INDEX IF NOT EXISTS idx_conversion_tracking_user ON conversion_tracking(user_id);
CREATE INDEX IF NOT EXISTS idx_conversion_tracking_event ON conversion_tracking(event_id);
CREATE INDEX IF NOT EXISTS idx_conversion_tracking_stage ON conversion_tracking(stage);
CREATE INDEX IF NOT EXISTS idx_roi_calculations_event ON roi_calculations(event_id);
CREATE INDEX IF NOT EXISTS idx_roi_calculations_brand ON roi_calculations(brand_id);

-- Insert sample AI tagging results
INSERT OR IGNORE INTO ai_tagging_results (content_id, products, tags, confidence_score) VALUES
('content_1', '[{"productId": "prod_1", "productName": "Sample Product", "brandId": "brand_1", "confidence": 0.85}]', '["sample product", "brand_brand_1"]', 0.85),
('content_2', '[{"productId": "prod_2", "productName": "Tech Gadget", "brandId": "brand_1", "confidence": 0.92}]', '["tech gadget", "brand_brand_1", "technology"]', 0.92);

-- Insert sample product detections
INSERT OR IGNORE INTO product_detections (content_id, product_id, brand_id, confidence, bounding_box) VALUES
('content_1', 'prod_1', 'brand_1', 0.85, '{"x": 100, "y": 150, "width": 200, "height": 250}'),
('content_2', 'prod_2', 'brand_1', 0.92, '{"x": 50, "y": 75, "width": 300, "height": 400}');

-- Insert sample conversion tracking
INSERT OR IGNORE INTO conversion_tracking (user_id, event_id, stage, metadata) VALUES
('user_1', 'event_1', 'attendance', '{"checkin_time": "2024-03-15T10:00:00Z"}'),
('user_1', 'event_1', 'content_view', '{"content_id": "content_1", "view_duration": 45}'),
('user_1', 'event_1', 'engagement', '{"action": "poll_vote", "poll_id": "poll_1"}'),
('user_1', 'event_1', 'website_visit', '{"url": "https://brand1.com", "referrer": "app"}'),
('user_1', 'event_1', 'purchase', '{"purchase_id": "purchase_1", "amount": 29.99}');

-- Insert sample ROI calculations
INSERT OR IGNORE INTO roi_calculations (event_id, brand_id, total_investment, total_revenue, roi_percentage) VALUES
('event_1', 'brand_1', 10000.0, 15000.0, 50.0);

-- Insert sample attribution weights for linear model
INSERT OR IGNORE INTO attribution_weights (event_id, attribution_model, stage, weight) VALUES
('event_1', 'linear', 'attendance', 0.2),
('event_1', 'linear', 'content_view', 0.2),
('event_1', 'linear', 'engagement', 0.2),
('event_1', 'linear', 'website_visit', 0.2),
('event_1', 'linear', 'purchase', 0.2);