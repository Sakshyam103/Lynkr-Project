-- UX Improvements Migration
-- Adds tables for usability testing and user experience monitoring

-- User sessions table for usability testing
CREATE TABLE IF NOT EXISTS user_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    session_start DATETIME NOT NULL,
    session_end DATETIME,
    duration INTEGER, -- in seconds
    device_type TEXT,
    app_version TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- User actions table for tracking user interactions
CREATE TABLE IF NOT EXISTS user_actions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('tap', 'swipe', 'scroll', 'navigation', 'input', 'long_press')),
    screen TEXT NOT NULL,
    element TEXT,
    duration INTEGER DEFAULT 0, -- in milliseconds
    coordinates TEXT, -- JSON with x, y coordinates
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id)
);

-- User errors table for tracking UX issues
CREATE TABLE IF NOT EXISTS user_errors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('crash', 'network_error', 'validation_error', 'ui_error', 'timeout')),
    message TEXT NOT NULL,
    screen TEXT NOT NULL,
    element TEXT,
    recoverable BOOLEAN DEFAULT TRUE,
    stack_trace TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id)
);

-- App feedback table for user satisfaction ratings
CREATE TABLE IF NOT EXISTS app_feedback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('satisfaction', 'feature_request', 'bug_report', 'general')),
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    screen TEXT,
    category TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Performance metrics table
CREATE TABLE IF NOT EXISTS app_performance_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    metric_name TEXT NOT NULL,
    metric_value REAL NOT NULL,
    metric_unit TEXT,
    screen TEXT,
    recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id)
);

-- A/B test experiments table
CREATE TABLE IF NOT EXISTS ab_experiments (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    variant_a TEXT NOT NULL, -- JSON configuration
    variant_b TEXT NOT NULL, -- JSON configuration
    start_date DATETIME NOT NULL,
    end_date DATETIME,
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'paused', 'completed')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- A/B test assignments table
CREATE TABLE IF NOT EXISTS ab_assignments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    experiment_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    variant TEXT NOT NULL CHECK (variant IN ('A', 'B')),
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (experiment_id) REFERENCES ab_experiments(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(experiment_id, user_id)
);

-- Accessibility settings table
CREATE TABLE IF NOT EXISTS accessibility_settings (
    user_id TEXT PRIMARY KEY,
    screen_reader_enabled BOOLEAN DEFAULT FALSE,
    reduce_motion BOOLEAN DEFAULT FALSE,
    high_contrast BOOLEAN DEFAULT FALSE,
    large_text BOOLEAN DEFAULT FALSE,
    voice_over_enabled BOOLEAN DEFAULT FALSE,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Screen analytics table
CREATE TABLE IF NOT EXISTS screen_analytics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    screen_name TEXT NOT NULL,
    enter_time DATETIME NOT NULL,
    exit_time DATETIME,
    duration INTEGER, -- in seconds
    bounce BOOLEAN DEFAULT FALSE, -- user left immediately
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_sessions_user ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_start ON user_sessions(session_start);
CREATE INDEX IF NOT EXISTS idx_user_actions_session ON user_actions(session_id);
CREATE INDEX IF NOT EXISTS idx_user_actions_type_screen ON user_actions(type, screen);
CREATE INDEX IF NOT EXISTS idx_user_actions_created ON user_actions(created_at);
CREATE INDEX IF NOT EXISTS idx_user_errors_session ON user_errors(session_id);
CREATE INDEX IF NOT EXISTS idx_user_errors_type ON user_errors(type);
CREATE INDEX IF NOT EXISTS idx_user_errors_screen ON user_errors(screen);
CREATE INDEX IF NOT EXISTS idx_app_feedback_user ON app_feedback(user_id);
CREATE INDEX IF NOT EXISTS idx_app_feedback_type ON app_feedback(type);
CREATE INDEX IF NOT EXISTS idx_app_feedback_rating ON app_feedback(rating);
CREATE INDEX IF NOT EXISTS idx_performance_metrics_session ON app_performance_metrics(session_id);
CREATE INDEX IF NOT EXISTS idx_performance_metrics_name ON app_performance_metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_ab_assignments_experiment ON ab_assignments(experiment_id);
CREATE INDEX IF NOT EXISTS idx_ab_assignments_user ON ab_assignments(user_id);
CREATE INDEX IF NOT EXISTS idx_screen_analytics_session ON screen_analytics(session_id);
CREATE INDEX IF NOT EXISTS idx_screen_analytics_screen ON screen_analytics(screen_name);

-- Views for UX analytics
CREATE VIEW IF NOT EXISTS ux_metrics_summary AS
SELECT 
    DATE(us.session_start) as date,
    COUNT(DISTINCT us.id) as total_sessions,
    COUNT(DISTINCT us.user_id) as unique_users,
    AVG(us.duration) as avg_session_duration,
    COUNT(CASE WHEN us.session_end IS NOT NULL THEN 1 END) * 100.0 / COUNT(*) as completion_rate,
    COUNT(ue.id) * 100.0 / COUNT(us.id) as error_rate
FROM user_sessions us
LEFT JOIN user_errors ue ON us.id = ue.session_id
GROUP BY DATE(us.session_start);

CREATE VIEW IF NOT EXISTS screen_performance AS
SELECT 
    sa.screen_name,
    COUNT(*) as visits,
    AVG(sa.duration) as avg_duration,
    COUNT(CASE WHEN sa.bounce = 1 THEN 1 END) * 100.0 / COUNT(*) as bounce_rate,
    COUNT(ue.id) as error_count
FROM screen_analytics sa
LEFT JOIN user_errors ue ON sa.session_id = ue.session_id AND sa.screen_name = ue.screen
GROUP BY sa.screen_name;

CREATE VIEW IF NOT EXISTS user_satisfaction_trends AS
SELECT 
    DATE(created_at) as date,
    type,
    AVG(rating) as avg_rating,
    COUNT(*) as feedback_count
FROM app_feedback
WHERE rating IS NOT NULL
GROUP BY DATE(created_at), type;

-- Insert sample data for testing
INSERT OR IGNORE INTO user_sessions (id, user_id, session_start, session_end, duration, device_type, app_version) VALUES
('session_1', 'user_1', datetime('now', '-2 hours'), datetime('now', '-1 hour'), 3600, 'iOS', '1.0.0'),
('session_2', 'user_2', datetime('now', '-1 hour'), datetime('now', '-30 minutes'), 1800, 'Android', '1.0.0');

INSERT OR IGNORE INTO user_actions (session_id, type, screen, element, duration) VALUES
('session_1', 'tap', 'EventList', 'event_card_1', 250),
('session_1', 'navigation', 'EventList', 'event_details', 500),
('session_1', 'tap', 'EventDetails', 'checkin_button', 300),
('session_2', 'swipe', 'EventList', 'event_list', 150),
('session_2', 'tap', 'EventList', 'event_card_2', 200);

INSERT OR IGNORE INTO user_errors (session_id, type, message, screen, recoverable) VALUES
('session_1', 'network_error', 'Failed to load event details', 'EventDetails', TRUE),
('session_2', 'ui_error', 'Button not responding', 'EventList', TRUE);

INSERT OR IGNORE INTO app_feedback (user_id, type, rating, comment, screen) VALUES
('user_1', 'satisfaction', 4, 'Great app, easy to use', 'EventList'),
('user_2', 'satisfaction', 5, 'Love the check-in feature', 'EventDetails'),
('user_1', 'bug_report', 2, 'App crashes when uploading photos', 'ContentCreation');

INSERT OR IGNORE INTO app_performance_metrics (session_id, metric_name, metric_value, metric_unit, screen) VALUES
('session_1', 'load_time', 1.2, 'seconds', 'EventList'),
('session_1', 'memory_usage', 45.6, 'MB', 'EventDetails'),
('session_2', 'load_time', 0.8, 'seconds', 'EventList');

INSERT OR IGNORE INTO ab_experiments (id, name, description, variant_a, variant_b, start_date) VALUES
('exp_1', 'Check-in Button Color', 'Test blue vs green check-in button', '{"color": "blue"}', '{"color": "green"}', datetime('now', '-7 days'));

INSERT OR IGNORE INTO ab_assignments (experiment_id, user_id, variant) VALUES
('exp_1', 'user_1', 'A'),
('exp_1', 'user_2', 'B');

INSERT OR IGNORE INTO accessibility_settings (user_id, screen_reader_enabled, large_text) VALUES
('user_1', FALSE, TRUE),
('user_2', TRUE, FALSE);

INSERT OR IGNORE INTO screen_analytics (session_id, screen_name, enter_time, exit_time, duration) VALUES
('session_1', 'EventList', datetime('now', '-2 hours'), datetime('now', '-2 hours', '+5 minutes'), 300),
('session_1', 'EventDetails', datetime('now', '-2 hours', '+5 minutes'), datetime('now', '-2 hours', '+15 minutes'), 600),
('session_2', 'EventList', datetime('now', '-1 hour'), datetime('now', '-1 hour', '+3 minutes'), 180);