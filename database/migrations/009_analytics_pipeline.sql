-- Analytics Pipeline Migration
-- Adds tables for analytics events, aggregated data, and reporting

-- Analytics events table for raw event tracking
CREATE TABLE IF NOT EXISTS analytics_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,
    user_id TEXT NOT NULL,
    event_id TEXT,
    data TEXT, -- JSON data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Engagement summary table for aggregated metrics
CREATE TABLE IF NOT EXISTS engagement_summary (
    event_id TEXT PRIMARY KEY,
    total_attendees INTEGER DEFAULT 0,
    content_pieces INTEGER DEFAULT 0,
    engagement_rate REAL DEFAULT 0,
    average_rating REAL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Attendance summary table for daily aggregates
CREATE TABLE IF NOT EXISTS attendance_summary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id TEXT NOT NULL,
    date TEXT NOT NULL,
    attendees INTEGER DEFAULT 0,
    check_ins INTEGER DEFAULT 0,
    avg_duration REAL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    UNIQUE(event_id, date)
);

-- Content summary table for content performance
CREATE TABLE IF NOT EXISTS content_summary (
    event_id TEXT PRIMARY KEY,
    total_content INTEGER DEFAULT 0,
    total_views INTEGER DEFAULT 0,
    total_shares INTEGER DEFAULT 0,
    avg_sentiment REAL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Engagement metrics table for detailed tracking
CREATE TABLE IF NOT EXISTS engagement_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    event_id TEXT,
    action TEXT NOT NULL,
    value REAL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Real-time stats table for current activity
CREATE TABLE IF NOT EXISTS realtime_stats (
    event_id TEXT PRIMARY KEY,
    active_users INTEGER DEFAULT 0,
    recent_content INTEGER DEFAULT 0,
    current_engagement REAL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_analytics_events_type ON analytics_events(type);
CREATE INDEX IF NOT EXISTS idx_analytics_events_user ON analytics_events(user_id);
CREATE INDEX IF NOT EXISTS idx_analytics_events_event ON analytics_events(event_id);
CREATE INDEX IF NOT EXISTS idx_analytics_events_created ON analytics_events(created_at);
CREATE INDEX IF NOT EXISTS idx_engagement_metrics_user ON engagement_metrics(user_id);
CREATE INDEX IF NOT EXISTS idx_engagement_metrics_event ON engagement_metrics(event_id);
CREATE INDEX IF NOT EXISTS idx_engagement_metrics_action ON engagement_metrics(action);
CREATE INDEX IF NOT EXISTS idx_attendance_summary_date ON attendance_summary(date);

-- Insert sample analytics data
INSERT OR IGNORE INTO engagement_summary (event_id, total_attendees, content_pieces, engagement_rate) VALUES
('event_1', 1247, 89, 23.5);

INSERT OR IGNORE INTO attendance_summary (event_id, date, attendees, check_ins) VALUES
('event_1', '2024-01-15', 120, 145),
('event_1', '2024-01-16', 180, 210),
('event_1', '2024-01-17', 250, 290);

INSERT OR IGNORE INTO content_summary (event_id, total_content, total_views, total_shares) VALUES
('event_1', 89, 15420, 1250);