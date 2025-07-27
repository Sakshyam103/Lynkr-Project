-- Rewards System Migration
-- Adds tables for incentive system and pulse surveys

-- Rewards table for tracking user rewards
CREATE TABLE IF NOT EXISTS rewards (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('quality_content', 'event_participation', 'feedback_completion', 'social_sharing', 'referral')),
    points INTEGER NOT NULL,
    description TEXT NOT NULL,
    event_id TEXT,
    content_id TEXT,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'delivered', 'expired')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    delivered_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Content quality scores table
CREATE TABLE IF NOT EXISTS content_quality_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content_id TEXT NOT NULL,
    quality_score INTEGER NOT NULL,
    engagement_score INTEGER NOT NULL,
    viral_score INTEGER NOT NULL,
    overall_score INTEGER NOT NULL,
    evaluated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (content_id) REFERENCES content(id)
);

-- Reward criteria table for defining reward rules
CREATE TABLE IF NOT EXISTS reward_criteria (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    points INTEGER NOT NULL,
    criteria TEXT NOT NULL, -- JSON with criteria rules
    active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- User notification preferences for surveys
CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id TEXT PRIMARY KEY,
    pulse_surveys BOOLEAN DEFAULT TRUE,
    reward_notifications BOOLEAN DEFAULT TRUE,
    survey_reminders BOOLEAN DEFAULT TRUE,
    optimal_time TEXT DEFAULT '18:00', -- Preferred notification time
    timezone TEXT DEFAULT 'UTC',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Survey scheduling table for managing survey delivery
CREATE TABLE IF NOT EXISTS survey_schedule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    survey_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    scheduled_for DATETIME NOT NULL,
    delivered_at DATETIME,
    status TEXT DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'delivered', 'expired', 'skipped')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (survey_id) REFERENCES pulse_surveys(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_rewards_user ON rewards(user_id);
CREATE INDEX IF NOT EXISTS idx_rewards_type ON rewards(type);
CREATE INDEX IF NOT EXISTS idx_rewards_status ON rewards(status);
CREATE INDEX IF NOT EXISTS idx_rewards_event ON rewards(event_id);
CREATE INDEX IF NOT EXISTS idx_rewards_content ON rewards(content_id);
CREATE INDEX IF NOT EXISTS idx_rewards_created ON rewards(created_at);
CREATE INDEX IF NOT EXISTS idx_content_quality_content ON content_quality_scores(content_id);
CREATE INDEX IF NOT EXISTS idx_content_quality_score ON content_quality_scores(overall_score);
CREATE INDEX IF NOT EXISTS idx_survey_schedule_user ON survey_schedule(user_id);
CREATE INDEX IF NOT EXISTS idx_survey_schedule_scheduled ON survey_schedule(scheduled_for);
CREATE INDEX IF NOT EXISTS idx_survey_schedule_status ON survey_schedule(status);

-- Insert sample reward criteria
INSERT OR IGNORE INTO reward_criteria (id, type, name, description, points, criteria) VALUES
('quality_content_high', 'quality_content', 'High Quality Content', 'Content with high engagement and quality score', 25, '{"min_quality_score": 20, "min_engagement": 10}'),
('event_checkin', 'event_participation', 'Event Check-in', 'Bonus for checking into an event', 10, '{"action": "checkin"}'),
('poll_participation', 'feedback_completion', 'Poll Participation', 'Completing a poll or survey', 5, '{"action": "poll_vote"}'),
('content_share', 'social_sharing', 'Content Sharing', 'Sharing content on social media', 15, '{"action": "share"}'),
('friend_referral', 'referral', 'Friend Referral', 'Referring a friend who joins an event', 50, '{"action": "referral", "min_events": 1}');

-- Insert sample rewards
INSERT OR IGNORE INTO rewards (id, user_id, type, points, description, event_id, status) VALUES
('reward_1', 'user_1', 'quality_content', 25, 'High-quality content reward', 'event_1', 'delivered'),
('reward_2', 'user_1', 'event_participation', 10, 'Event check-in bonus', 'event_1', 'delivered'),
('reward_3', 'user_2', 'feedback_completion', 5, 'Poll participation reward', 'event_1', 'delivered');

-- Insert sample content quality scores
INSERT OR IGNORE INTO content_quality_scores (content_id, quality_score, engagement_score, viral_score, overall_score) VALUES
('content_1', 25, 15, 10, 50),
('content_2', 20, 12, 8, 40);

-- Insert sample notification preferences
INSERT OR IGNORE INTO notification_preferences (user_id, pulse_surveys, reward_notifications, optimal_time) VALUES
('user_1', TRUE, TRUE, '18:00'),
('user_2', TRUE, FALSE, '20:00');