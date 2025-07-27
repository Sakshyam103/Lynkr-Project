-- Missing tables for analytics service

-- Content analytics table
CREATE TABLE IF NOT EXISTS content_analytics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    action TEXT NOT NULL, -- 'view', 'like', 'share', 'comment'
    metadata TEXT, -- JSON metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Analytics events table
CREATE TABLE IF NOT EXISTS analytics_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL, -- 'page_view', 'click', 'interaction'
    user_id TEXT,
    event_id TEXT,
    data TEXT, -- JSON data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Sentiment analysis table
CREATE TABLE IF NOT EXISTS sentiment_analysis (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    result TEXT NOT NULL, -- JSON result with score, sentiment
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Attendances table (if not exists)
CREATE TABLE IF NOT EXISTS attendances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    checkin_time DATETIME,
    checkout_time DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    latitude REAL,
    longitude REAL,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Update content table to include view_count and share_count
-- ALTER TABLE content ADD COLUMN view_count INTEGER DEFAULT 0;
-- ALTER TABLE content ADD COLUMN share_count INTEGER DEFAULT 0;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_content_analytics_content_id ON content_analytics(content_id);
CREATE INDEX IF NOT EXISTS idx_content_analytics_action ON content_analytics(action);
CREATE INDEX IF NOT EXISTS idx_analytics_events_event_id ON analytics_events(event_id);
CREATE INDEX IF NOT EXISTS idx_analytics_events_type ON analytics_events(type);
CREATE INDEX IF NOT EXISTS idx_sentiment_analysis_content_id ON sentiment_analysis(content_id);
CREATE INDEX IF NOT EXISTS idx_attendances_event_id ON attendances(event_id);
CREATE INDEX IF NOT EXISTS idx_attendances_user_id ON attendances(user_id);

-- Discount codes tables

-- Discount codes table
CREATE TABLE IF NOT EXISTS discount_codes (
    id TEXT PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    event_id TEXT NOT NULL,
    brand_id TEXT NOT NULL,
    discount_pct DECIMAL(5,2) NOT NULL, -- e.g., 20.50 for 20.5%
    max_uses INTEGER NOT NULL,
    used_count INTEGER DEFAULT 0,
    expires_at DATETIME NOT NULL,
    status TEXT DEFAULT 'active', -- 'active', 'inactive', 'expired'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Discount code usage tracking
CREATE TABLE IF NOT EXISTS discount_code_usage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    discount_code_id TEXT NOT NULL,
    user_id TEXT,
    order_id TEXT,
    used_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (discount_code_id) REFERENCES discount_codes(id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_discount_codes_code ON discount_codes(code);
CREATE INDEX IF NOT EXISTS idx_discount_codes_event_id ON discount_codes(event_id);
CREATE INDEX IF NOT EXISTS idx_discount_codes_brand_id ON discount_codes(brand_id);
CREATE INDEX IF NOT EXISTS idx_discount_code_usage_code_id ON discount_code_usage(discount_code_id);
CREATE INDEX IF NOT EXISTS idx_discount_code_usage_user_id ON discount_code_usage(user_id);


-- Pixel tracking tables

-- Pixel events table
CREATE TABLE IF NOT EXISTS pixel_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id TEXT,
    brand_id TEXT,
    campaign_id TEXT,
    user_id TEXT,
    event_type TEXT NOT NULL, -- 'email_open', 'page_view', 'ad_view', 'click'
    user_agent TEXT,
    ip_address TEXT,
    referrer TEXT,
    metadata TEXT, -- JSON metadata
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Pixel campaigns table
CREATE TABLE IF NOT EXISTS pixel_campaigns (
    id TEXT PRIMARY KEY,
    brand_id TEXT NOT NULL,
    event_id TEXT,
    campaign_name TEXT NOT NULL,
    pixel_url TEXT NOT NULL,
    status TEXT DEFAULT 'active', -- 'active', 'paused', 'completed'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_pixel_events_event_id ON pixel_events(event_id);
CREATE INDEX IF NOT EXISTS idx_pixel_events_brand_id ON pixel_events(brand_id);
CREATE INDEX IF NOT EXISTS idx_pixel_events_user_id ON pixel_events(user_id);
CREATE INDEX IF NOT EXISTS idx_pixel_events_event_type ON pixel_events(event_type);
CREATE INDEX IF NOT EXISTS idx_pixel_campaigns_brand_id ON pixel_campaigns(brand_id);
CREATE INDEX IF NOT EXISTS idx_pixel_campaigns_event_id ON pixel_campaigns(event_id);


CREATE TABLE IF NOT EXISTS ai_tagging_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    products TEXT, -- JSON array of detected products
    tags TEXT, -- JSON array of generated tags
    confidence_score DECIMAL(3,2), -- 0.00 to 1.00
    processed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Detected products table (normalized)
CREATE TABLE IF NOT EXISTS detected_products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    product_name TEXT NOT NULL,
    brand_name TEXT,
    confidence DECIMAL(3,2), -- 0.00 to 1.00
    bounding_box TEXT, -- JSON coordinates
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Generated tags table (normalized)
CREATE TABLE IF NOT EXISTS generated_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    tag_name TEXT NOT NULL,
    tag_type TEXT, -- 'object', 'brand', 'emotion', 'activity'
    confidence DECIMAL(3,2), -- 0.00 to 1.00
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_ai_tagging_results_content_id ON ai_tagging_results(content_id);
CREATE INDEX IF NOT EXISTS idx_detected_products_content_id ON detected_products(content_id);
CREATE INDEX IF NOT EXISTS idx_detected_products_brand_name ON detected_products(brand_name);
CREATE INDEX IF NOT EXISTS idx_generated_tags_content_id ON generated_tags(content_id);
CREATE INDEX IF NOT EXISTS idx_generated_tags_tag_name ON generated_tags(tag_name);
CREATE INDEX IF NOT EXISTS idx_generated_tags_tag_type ON generated_tags(tag_type);

-- Survey and pulse survey tables

-- Pulse surveys table
CREATE TABLE IF NOT EXISTS pulse_surveys (
                                             id TEXT PRIMARY KEY,
                                             event_id TEXT,
                                             brand_id TEXT,
                                             title TEXT NOT NULL,
                                             description TEXT,
                                             questions TEXT NOT NULL, -- JSON array of questions
                                             reward_points INTEGER DEFAULT 0,
                                             status TEXT DEFAULT 'active', -- 'active', 'inactive', 'completed'
                                             created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                             expires_at DATETIME,
                                             FOREIGN KEY (event_id) REFERENCES events(id)
    );

-- Pulse survey responses table
CREATE TABLE IF NOT EXISTS pulse_survey_responses (
                                                      id INTEGER PRIMARY KEY AUTOINCREMENT,
                                                      survey_id TEXT NOT NULL,
                                                      user_id TEXT NOT NULL,
                                                      responses TEXT NOT NULL, -- JSON responses
                                                      completed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                                      FOREIGN KEY (survey_id) REFERENCES pulse_surveys(id)
    );

-- Survey analytics table
CREATE TABLE IF NOT EXISTS survey_analytics (
                                                id INTEGER PRIMARY KEY AUTOINCREMENT,
                                                survey_id TEXT NOT NULL,
                                                total_responses INTEGER DEFAULT 0,
                                                completion_rate DECIMAL(5,2) DEFAULT 0.00,
    avg_completion_time INTEGER DEFAULT 0, -- in seconds
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (survey_id) REFERENCES pulse_surveys(id)
    );

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_pulse_surveys_event_id ON pulse_surveys(event_id);
CREATE INDEX IF NOT EXISTS idx_pulse_surveys_brand_id ON pulse_surveys(brand_id);
CREATE INDEX IF NOT EXISTS idx_pulse_survey_responses_survey_id ON pulse_survey_responses(survey_id);
CREATE INDEX IF NOT EXISTS idx_pulse_survey_responses_user_id ON pulse_survey_responses(user_id);
CREATE INDEX IF NOT EXISTS idx_survey_analytics_survey_id ON survey_analytics(survey_id);

CREATE TABLE export_requests (
                                id TEXT PRIMARY KEY,
                                export_request_id TEXT NOT NULL,
                                brand_id TEXT NOT NULL,
                                event_id TEXT NOT NULL,
                                data_type TEXT NOT NULL,
                                format TEXT NOT NULL,
                                file_path TEXT,
                                file_size INTEGER,
                                download_url TEXT,
                                status TEXT NOT NULL DEFAULT 'processing',
                                error_message TEXT,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                completed_at TIMESTAMP,
                                downloaded_at TIMESTAMP,
                                FOREIGN KEY (export_request_id) REFERENCES export_requests(id),
                                FOREIGN KEY (brand_id) REFERENCES brands(id),
                                FOREIGN KEY (event_id) REFERENCES events(id)
);
