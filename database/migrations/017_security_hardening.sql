-- Security Hardening Migration
-- Adds security audit tables and enhanced privacy controls

-- Vulnerability reports table
CREATE TABLE IF NOT EXISTS vulnerability_reports (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL CHECK (type IN ('sql_injection', 'xss', 'csrf', 'weak_authentication', 'data_exposure', 'data_validation')),
    severity TEXT NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    description TEXT NOT NULL,
    endpoint TEXT,
    status TEXT DEFAULT 'open' CHECK (status IN ('open', 'investigating', 'fixed', 'false_positive')),
    assigned_to TEXT,
    fixed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Security events table
CREATE TABLE IF NOT EXISTS security_events (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL CHECK (type IN ('login_attempt', 'failed_login', 'suspicious_activity', 'rate_limit_exceeded', 'data_access', 'privilege_escalation')),
    user_id TEXT,
    ip_address TEXT,
    user_agent TEXT,
    details TEXT,
    severity TEXT DEFAULT 'medium' CHECK (severity IN ('low', 'medium', 'high')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Enhanced consent types table
CREATE TABLE IF NOT EXISTS consent_types (
    type TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    required BOOLEAN DEFAULT FALSE,
    version TEXT DEFAULT '1.0',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- User consent records table
CREATE TABLE IF NOT EXISTS user_consent_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    consent_type TEXT NOT NULL,
    granted BOOLEAN NOT NULL,
    version TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (consent_type) REFERENCES consent_types(type)
);

-- Data retention policies table
CREATE TABLE IF NOT EXISTS data_retention_policies (
    data_type TEXT PRIMARY KEY,
    retention_days INTEGER NOT NULL,
    auto_delete BOOLEAN DEFAULT TRUE,
    last_cleanup DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Scheduled deletions table
CREATE TABLE IF NOT EXISTS scheduled_deletions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    scheduled_for DATETIME NOT NULL,
    status TEXT DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'processing', 'completed', 'cancelled')),
    reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Privacy audit log table
CREATE TABLE IF NOT EXISTS privacy_audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL CHECK (action IN ('consent_updated', 'data_exported', 'data_anonymized', 'data_deleted', 'deletion_requested', 'data_cleanup')),
    user_id TEXT,
    details TEXT,
    ip_address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Failed login attempts table
CREATE TABLE IF NOT EXISTS failed_login_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT,
    ip_address TEXT NOT NULL,
    user_agent TEXT,
    attempt_count INTEGER DEFAULT 1,
    last_attempt DATETIME DEFAULT CURRENT_TIMESTAMP,
    blocked_until DATETIME
);

-- API rate limiting table
CREATE TABLE IF NOT EXISTS rate_limit_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    identifier TEXT NOT NULL, -- user_id or ip_address
    endpoint TEXT NOT NULL,
    request_count INTEGER DEFAULT 1,
    window_start DATETIME DEFAULT CURRENT_TIMESTAMP,
    blocked_until DATETIME
);

-- Security configuration table
CREATE TABLE IF NOT EXISTS security_config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Add security columns to existing users table
ALTER TABLE users ADD COLUMN deletion_requested BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN deletion_requested_at DATETIME;
ALTER TABLE users ADD COLUMN last_password_change DATETIME;
ALTER TABLE users ADD COLUMN failed_login_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN account_locked_until DATETIME;
ALTER TABLE users ADD COLUMN two_factor_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN two_factor_secret TEXT;

-- Indexes for security tables
CREATE INDEX IF NOT EXISTS idx_vulnerability_reports_type ON vulnerability_reports(type);
CREATE INDEX IF NOT EXISTS idx_vulnerability_reports_severity ON vulnerability_reports(severity);
CREATE INDEX IF NOT EXISTS idx_vulnerability_reports_status ON vulnerability_reports(status);
CREATE INDEX IF NOT EXISTS idx_security_events_type ON security_events(type);
CREATE INDEX IF NOT EXISTS idx_security_events_user ON security_events(user_id);
CREATE INDEX IF NOT EXISTS idx_security_events_ip ON security_events(ip_address);
CREATE INDEX IF NOT EXISTS idx_security_events_created ON security_events(created_at);
CREATE INDEX IF NOT EXISTS idx_user_consent_user_type ON user_consent_records(user_id, consent_type);
CREATE INDEX IF NOT EXISTS idx_scheduled_deletions_scheduled ON scheduled_deletions(scheduled_for);
CREATE INDEX IF NOT EXISTS idx_scheduled_deletions_status ON scheduled_deletions(status);
CREATE INDEX IF NOT EXISTS idx_privacy_audit_user ON privacy_audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_privacy_audit_action ON privacy_audit_log(action);
CREATE INDEX IF NOT EXISTS idx_failed_login_email ON failed_login_attempts(email);
CREATE INDEX IF NOT EXISTS idx_failed_login_ip ON failed_login_attempts(ip_address);
CREATE INDEX IF NOT EXISTS idx_rate_limit_identifier ON rate_limit_records(identifier, endpoint);

-- Security views for monitoring
CREATE VIEW IF NOT EXISTS security_dashboard AS
SELECT 
    'vulnerabilities' as metric,
    COUNT(*) as count,
    severity
FROM vulnerability_reports 
WHERE status = 'open'
GROUP BY severity
UNION ALL
SELECT 
    'security_events' as metric,
    COUNT(*) as count,
    type as severity
FROM security_events 
WHERE created_at > datetime('now', '-24 hours')
GROUP BY type;

CREATE VIEW IF NOT EXISTS privacy_compliance AS
SELECT 
    u.id as user_id,
    u.email,
    COUNT(ucr.id) as consent_records,
    MAX(ucr.created_at) as last_consent_update,
    u.deletion_requested,
    u.deletion_requested_at
FROM users u
LEFT JOIN user_consent_records ucr ON u.id = ucr.user_id
GROUP BY u.id;

-- Insert default consent types
INSERT OR IGNORE INTO consent_types (type, description, required, version) VALUES
('analytics_tracking', 'Allow collection of usage analytics to improve the app', FALSE, '2.0'),
('marketing_communications', 'Receive marketing emails and promotional content', FALSE, '2.0'),
('data_sharing_partners', 'Share anonymized data with trusted partners', FALSE, '2.0'),
('location_tracking', 'Track location for event check-ins and geofencing', TRUE, '2.0'),
('content_usage_rights', 'Allow brands to use your content for marketing', FALSE, '2.0'),
('personalized_recommendations', 'Receive personalized event and product recommendations', FALSE, '2.0');

-- Insert default data retention policies
INSERT OR IGNORE INTO data_retention_policies (data_type, retention_days, auto_delete) VALUES
('analytics_events', 365, TRUE),
('security_events', 90, TRUE),
('export_requests', 30, TRUE),
('query_performance_logs', 7, TRUE),
('failed_login_attempts', 30, TRUE),
('rate_limit_records', 1, TRUE);

-- Insert default security configuration
INSERT OR IGNORE INTO security_config (key, value, description) VALUES
('max_login_attempts', '5', 'Maximum failed login attempts before account lockout'),
('lockout_duration_minutes', '30', 'Account lockout duration in minutes'),
('password_min_length', '8', 'Minimum password length requirement'),
('session_timeout_minutes', '60', 'Session timeout in minutes'),
('rate_limit_requests_per_minute', '100', 'API rate limit per minute per user'),
('require_2fa_for_brands', 'false', 'Require two-factor authentication for brand accounts');

-- Sample security events for testing
INSERT OR IGNORE INTO security_events (id, type, user_id, ip_address, details) VALUES
('sec_event_1', 'failed_login', 'user_1', '192.168.1.100', 'Failed login attempt with incorrect password'),
('sec_event_2', 'suspicious_activity', 'user_2', '10.0.0.50', 'Multiple rapid API requests detected'),
('sec_event_3', 'data_access', 'user_1', '192.168.1.100', 'User accessed personal data export');

-- Sample vulnerability reports
INSERT OR IGNORE INTO vulnerability_reports (id, type, severity, description, endpoint) VALUES
('vuln_1', 'data_validation', 'low', 'Input validation could be improved for user registration', '/api/v1/users/register'),
('vuln_2', 'weak_authentication', 'medium', 'Some users have weak password hashing', '/api/v1/auth/login');