-- Database schema for Brand Activations

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Users table
DROP TABLE  IF EXISTS users;
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    privacy_settings TEXT NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- -- Drop and recreate users table
-- DROP TABLE IF EXISTS users;
-- CREATE TABLE users (
--                        id TEXT PRIMARY KEY,
--                        email TEXT UNIQUE NOT NULL,
--                        password TEXT NOT NULL,
--                        name TEXT NOT NULL,
--                        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
-- );

-- Create index on users
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Brands table
CREATE TABLE IF NOT EXISTS brands (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    logo_url TEXT,
    contact_info TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on brands
CREATE INDEX IF NOT EXISTS idx_brands_name ON brands(name);

-- Events table
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    location TEXT NOT NULL,
    geofence_data TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    brand_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Create index on events
CREATE INDEX IF NOT EXISTS idx_events_brand_id ON events(brand_id);
CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);

-- Attendances table
CREATE TABLE IF NOT EXISTS attendances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    check_in_time TIMESTAMP NOT NULL,
    check_out_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Create index on attendances
CREATE INDEX IF NOT EXISTS idx_attendances_user_id ON attendances(user_id);
CREATE INDEX IF NOT EXISTS idx_attendances_event_id ON attendances(event_id);
CREATE INDEX IF NOT EXISTS idx_attendances_check_in_time ON attendances(check_in_time);

-- Content table
CREATE TABLE IF NOT EXISTS content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    type TEXT NOT NULL, -- 'photo', 'video', 'text'
    url TEXT NOT NULL,
    caption TEXT,
    tags TEXT,
    permissions TEXT NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Create index on content
CREATE INDEX IF NOT EXISTS idx_content_user_id ON content(user_id);
CREATE INDEX IF NOT EXISTS idx_content_event_id ON content(event_id);
CREATE INDEX IF NOT EXISTS idx_content_type ON content(type);

-- Interactions table
CREATE TABLE IF NOT EXISTS interactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content_id INTEGER NOT NULL,
    type TEXT NOT NULL, -- 'like', 'comment', 'share', 'reaction'
    data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Create index on interactions
CREATE INDEX IF NOT EXISTS idx_interactions_user_id ON interactions(user_id);
CREATE INDEX IF NOT EXISTS idx_interactions_content_id ON interactions(content_id);
CREATE INDEX IF NOT EXISTS idx_interactions_type ON interactions(type);

-- Campaigns table
CREATE TABLE IF NOT EXISTS campaigns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    brand_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Create index on campaigns
CREATE INDEX IF NOT EXISTS idx_campaigns_brand_id ON campaigns(brand_id);
CREATE INDEX IF NOT EXISTS idx_campaigns_start_date ON campaigns(start_date);

-- Conversions table
CREATE TABLE IF NOT EXISTS conversions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    campaign_id INTEGER NOT NULL,
    type TEXT NOT NULL, -- 'view', 'click', 'purchase'
    value REAL,
    timestamp TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
);

-- Create index on conversions
CREATE INDEX IF NOT EXISTS idx_conversions_user_id ON conversions(user_id);
CREATE INDEX IF NOT EXISTS idx_conversions_campaign_id ON conversions(campaign_id);
CREATE INDEX IF NOT EXISTS idx_conversions_type ON conversions(type);
CREATE INDEX IF NOT EXISTS idx_conversions_timestamp ON conversions(timestamp);


-- Content Tags Table
-- Stores tags associated with content items for events and brands

CREATE TABLE IF NOT EXISTS content_tags (
                                            id TEXT PRIMARY KEY,
                                            name TEXT NOT NULL,
                                            type TEXT NOT NULL,
                                            brand_id TEXT,
                                            event_id TEXT,
                                            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- Foreign key constraints
                                            FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
                                            FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_content_tags_event_id ON content_tags(event_id);
CREATE INDEX IF NOT EXISTS idx_content_tags_brand_id ON content_tags(brand_id);
CREATE INDEX IF NOT EXISTS idx_content_tags_name_type ON content_tags(name, type);
CREATE INDEX IF NOT EXISTS idx_content_tags_type ON content_tags(type);

-- Trigger to update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_content_tags_timestamp
    AFTER UPDATE ON content_tags
    FOR EACH ROW
BEGIN
    UPDATE content_tags SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;


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
