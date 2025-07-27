-- Content Management Migration
-- Adds tables for enhanced content management, tagging, and analytics

-- Content tags table for organizing content
CREATE TABLE IF NOT EXISTS content_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('brand', 'product', 'event', 'location', 'custom')),
    brand_id TEXT,
    event_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Content analytics table for tracking performance
CREATE TABLE IF NOT EXISTS content_analytics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    action TEXT NOT NULL,
    metadata TEXT, -- JSON data
    user_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Content moderation table for managing inappropriate content
CREATE TABLE IF NOT EXISTS content_moderation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected', 'flagged')),
    reason TEXT,
    moderator_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Content rights table for managing usage permissions
CREATE TABLE IF NOT EXISTS content_rights (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    brand_id TEXT NOT NULL,
    permission_type TEXT NOT NULL CHECK (permission_type IN ('view', 'commercial', 'modify', 'share')),
    granted BOOLEAN DEFAULT FALSE,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Update content table to include additional fields
ALTER TABLE content ADD COLUMN file_size INTEGER DEFAULT 0;
ALTER TABLE content ADD COLUMN duration INTEGER; -- For videos
ALTER TABLE content ADD COLUMN moderation_status TEXT DEFAULT 'pending';
ALTER TABLE content ADD COLUMN view_count INTEGER DEFAULT 0;
ALTER TABLE content ADD COLUMN share_count INTEGER DEFAULT 0;

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_content_tags_name ON content_tags(name);
CREATE INDEX IF NOT EXISTS idx_content_tags_type ON content_tags(type);
CREATE INDEX IF NOT EXISTS idx_content_tags_event ON content_tags(event_id);
CREATE INDEX IF NOT EXISTS idx_content_analytics_content ON content_analytics(content_id);
CREATE INDEX IF NOT EXISTS idx_content_analytics_action ON content_analytics(action);
CREATE INDEX IF NOT EXISTS idx_content_moderation_status ON content_moderation(status);
CREATE INDEX IF NOT EXISTS idx_content_rights_content ON content_rights(content_id);
CREATE INDEX IF NOT EXISTS idx_content_rights_brand ON content_rights(brand_id);

-- Insert sample content tags
INSERT OR IGNORE INTO content_tags (name, type, event_id) VALUES
('Brand Activation', 'event', 'event_1'),
('Product Demo', 'product', 'event_1'),
('Sample Brand', 'brand', 'event_1'),
('Tech Conference', 'event', 'event_1'),
('Innovation', 'custom', NULL),
('Networking', 'custom', NULL),
('Product Launch', 'product', NULL),
('Sponsorship', 'brand', NULL);