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