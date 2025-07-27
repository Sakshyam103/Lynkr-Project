-- Feedback System Migration
-- Adds tables for polls, surveys, sentiment analysis, and gamification

-- Polls table for interactive polls
CREATE TABLE IF NOT EXISTS polls (
    id TEXT PRIMARY KEY,
    question TEXT NOT NULL,
    event_id TEXT,
    brand_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Poll options table
CREATE TABLE IF NOT EXISTS poll_options (
    id TEXT PRIMARY KEY,
    poll_id TEXT NOT NULL,
    text TEXT NOT NULL,
    order_index INTEGER DEFAULT 0,
    FOREIGN KEY (poll_id) REFERENCES polls(id)
);

-- Poll votes table
CREATE TABLE IF NOT EXISTS poll_votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    poll_id TEXT NOT NULL,
    option_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (poll_id) REFERENCES polls(id),
    FOREIGN KEY (option_id) REFERENCES poll_options(id),
    UNIQUE(user_id, poll_id)
);

-- Slider feedback table
CREATE TABLE IF NOT EXISTS slider_feedback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    slider_id TEXT NOT NULL,
    value REAL NOT NULL CHECK (value >= 0 AND value <= 10),
    event_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Quick feedback table
CREATE TABLE IF NOT EXISTS quick_feedback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    widget_id TEXT NOT NULL,
    option_id TEXT NOT NULL,
    event_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Sentiment analysis table
CREATE TABLE IF NOT EXISTS sentiment_analysis (
    id TEXT PRIMARY KEY,
    content_id TEXT NOT NULL,
    text TEXT NOT NULL,
    result TEXT NOT NULL, -- JSON with score, magnitude, label
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- User badges table for gamification
CREATE TABLE IF NOT EXISTS user_badges (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    badge_id TEXT NOT NULL,
    earned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_new BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(user_id, badge_id)
);

-- Badge definitions table
CREATE TABLE IF NOT EXISTS badge_definitions (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    icon TEXT NOT NULL,
    rarity TEXT NOT NULL CHECK (rarity IN ('common', 'rare', 'epic', 'legendary')),
    points INTEGER DEFAULT 0,
    criteria TEXT -- JSON with earning criteria
);

-- User points table for gamification
CREATE TABLE IF NOT EXISTS user_points (
    user_id TEXT PRIMARY KEY,
    total_points INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_poll_votes_user ON poll_votes(user_id);
CREATE INDEX IF NOT EXISTS idx_poll_votes_poll ON poll_votes(poll_id);
CREATE INDEX IF NOT EXISTS idx_slider_feedback_user ON slider_feedback(user_id);
CREATE INDEX IF NOT EXISTS idx_slider_feedback_event ON slider_feedback(event_id);
CREATE INDEX IF NOT EXISTS idx_sentiment_content ON sentiment_analysis(content_id);
CREATE INDEX IF NOT EXISTS idx_user_badges_user ON user_badges(user_id);

-- Insert sample badge definitions
INSERT OR IGNORE INTO badge_definitions (id, title, description, icon, rarity, points) VALUES
('feedback_champion', 'Feedback Champion', 'Provided feedback on 10 products', '🏆', 'rare', 100),
('event_explorer', 'Event Explorer', 'Attended 5 different events', '🗺️', 'common', 50),
('content_creator', 'Content Creator', 'Shared 20 pieces of content', '📸', 'rare', 150),
('poll_participant', 'Poll Participant', 'Voted in 5 polls', '🗳️', 'common', 25),
('sentiment_positive', 'Positive Vibes', 'Maintained positive sentiment in 10 posts', '😊', 'epic', 200),
('early_adopter', 'Early Adopter', 'One of the first 100 users', '🚀', 'legendary', 500);

-- Insert sample polls
INSERT OR IGNORE INTO polls (id, question, event_id) VALUES
('poll_1', 'How would you rate this product demo?', 'event_1'),
('poll_2', 'What feature interests you most?', 'event_1');

INSERT OR IGNORE INTO poll_options (id, poll_id, text, order_index) VALUES
('opt_1_1', 'poll_1', 'Excellent', 1),
('opt_1_2', 'poll_1', 'Good', 2),
('opt_1_3', 'poll_1', 'Average', 3),
('opt_1_4', 'poll_1', 'Poor', 4),
('opt_2_1', 'poll_2', 'AI Integration', 1),
('opt_2_2', 'poll_2', 'Mobile App', 2),
('opt_2_3', 'poll_2', 'Analytics Dashboard', 3),
('opt_2_4', 'poll_2', 'API Access', 4);