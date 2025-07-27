-- Pixel Tracking Migration
-- Adds tables for pixel tracking and post-event engagement monitoring

-- Pixel events table for tracking brand-related searches and visits
CREATE TABLE IF NOT EXISTS pixel_events (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    event_id TEXT NOT NULL,
    brand_id TEXT NOT NULL,
    event_type TEXT NOT NULL CHECK (event_type IN ('page_view', 'search', 'website_visit', 'qr_scan', 'follow', 'save', 'share', 'delayed_24h', 'delayed_72h', 'delayed_7d')),
    url TEXT,
    referrer TEXT,
    user_agent TEXT,
    ip_address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- QR codes table for brand-specific tracking codes
CREATE TABLE IF NOT EXISTS qr_codes (
    id TEXT PRIMARY KEY,
    event_id TEXT NOT NULL,
    brand_id TEXT NOT NULL,
    code_data TEXT NOT NULL,
    target_url TEXT NOT NULL,
    scan_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Post-event engagement tracking
CREATE TABLE IF NOT EXISTS post_event_engagement (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    brand_id TEXT NOT NULL,
    engagement_type TEXT NOT NULL CHECK (engagement_type IN ('search', 'website_visit', 'social_follow', 'content_save', 'content_share')),
    engagement_data TEXT, -- JSON with additional data
    days_after_event INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Delayed pulse surveys table
CREATE TABLE IF NOT EXISTS pulse_surveys (
    id TEXT PRIMARY KEY,
    event_id TEXT NOT NULL,
    brand_id TEXT NOT NULL,
    survey_type TEXT NOT NULL CHECK (survey_type IN ('24h', '72h', '7d')),
    questions TEXT NOT NULL, -- JSON array of questions
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Pulse survey responses table
CREATE TABLE IF NOT EXISTS pulse_survey_responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    survey_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    responses TEXT NOT NULL, -- JSON responses
    completed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (survey_id) REFERENCES pulse_surveys(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(survey_id, user_id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_pixel_events_user ON pixel_events(user_id);
CREATE INDEX IF NOT EXISTS idx_pixel_events_event ON pixel_events(event_id);
CREATE INDEX IF NOT EXISTS idx_pixel_events_brand ON pixel_events(brand_id);
CREATE INDEX IF NOT EXISTS idx_pixel_events_type ON pixel_events(event_type);
CREATE INDEX IF NOT EXISTS idx_pixel_events_created ON pixel_events(created_at);
CREATE INDEX IF NOT EXISTS idx_qr_codes_event ON qr_codes(event_id);
CREATE INDEX IF NOT EXISTS idx_qr_codes_brand ON qr_codes(brand_id);
CREATE INDEX IF NOT EXISTS idx_post_event_engagement_user ON post_event_engagement(user_id);
CREATE INDEX IF NOT EXISTS idx_post_event_engagement_event ON post_event_engagement(event_id);
CREATE INDEX IF NOT EXISTS idx_post_event_engagement_days ON post_event_engagement(days_after_event);
CREATE INDEX IF NOT EXISTS idx_pulse_surveys_event ON pulse_surveys(event_id);
CREATE INDEX IF NOT EXISTS idx_pulse_survey_responses_survey ON pulse_survey_responses(survey_id);

-- Insert sample QR codes
INSERT OR IGNORE INTO qr_codes (id, event_id, brand_id, code_data, target_url) VALUES
('qr_1', 'event_1', 'brand_1', 'LYNKR_EVENT1_BRAND1', 'https://brand1.com/event-special'),
('qr_2', 'event_1', 'brand_1', 'LYNKR_EVENT1_PROMO', 'https://brand1.com/promo');

-- Insert sample pulse surveys
INSERT OR IGNORE INTO pulse_surveys (id, event_id, brand_id, survey_type, questions, expires_at) VALUES
('survey_24h_1', 'event_1', 'brand_1', '24h', '[{"question": "How likely are you to recommend this brand?", "type": "scale", "min": 1, "max": 10}]', datetime('now', '+1 day')),
('survey_72h_1', 'event_1', 'brand_1', '72h', '[{"question": "Have you visited the brand website since the event?", "type": "boolean"}]', datetime('now', '+3 days')),
('survey_7d_1', 'event_1', 'brand_1', '7d', '[{"question": "How has your perception of the brand changed?", "type": "multiple_choice", "options": ["Much better", "Better", "Same", "Worse", "Much worse"]}]', datetime('now', '+7 days'));

-- Insert sample pixel events
INSERT OR IGNORE INTO pixel_events (id, user_id, event_id, brand_id, event_type, url) VALUES
('pixel_1', 'user_1', 'event_1', 'brand_1', 'search', 'https://google.com/search?q=brand1+products'),
('pixel_2', 'user_1', 'event_1', 'brand_1', 'website_visit', 'https://brand1.com'),
('pixel_3', 'user_2', 'event_1', 'brand_1', 'qr_scan', 'https://brand1.com/event-special'),
('pixel_4', 'user_2', 'event_1', 'brand_1', 'follow', 'https://instagram.com/brand1');