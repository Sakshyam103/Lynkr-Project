-- Export System Migration
-- Adds tables for data export requests and CRM integrations

-- Export requests table
CREATE TABLE IF NOT EXISTS export_requests (
    id TEXT PRIMARY KEY,
    brand_id TEXT NOT NULL,
    event_id TEXT,
    data_type TEXT NOT NULL CHECK (data_type IN ('attendance', 'content', 'analytics', 'feedback', 'all')),
    format TEXT NOT NULL CHECK (format IN ('csv', 'json', 'xlsx')),
    status TEXT DEFAULT 'processing' CHECK (status IN ('processing', 'completed', 'failed', 'expired')),
    file_url TEXT,
    file_size INTEGER DEFAULT 0,
    record_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (brand_id) REFERENCES brands(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- CRM integrations table
CREATE TABLE IF NOT EXISTS crm_integrations (
    id TEXT PRIMARY KEY,
    brand_id TEXT NOT NULL,
    crm_type TEXT NOT NULL CHECK (crm_type IN ('salesforce', 'hubspot', 'mailchimp', 'pipedrive', 'zoho')),
    api_key TEXT NOT NULL,
    api_secret TEXT,
    webhook_url TEXT,
    sync_interval INTEGER DEFAULT 60, -- minutes
    last_sync DATETIME,
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'error')),
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id),
    UNIQUE(brand_id, crm_type)
);

-- CRM sync logs table
CREATE TABLE IF NOT EXISTS crm_sync_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    integration_id TEXT NOT NULL,
    event_id TEXT,
    sync_type TEXT NOT NULL CHECK (sync_type IN ('manual', 'scheduled', 'webhook')),
    records_synced INTEGER DEFAULT 0,
    status TEXT NOT NULL CHECK (status IN ('success', 'partial', 'failed')),
    error_message TEXT,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (integration_id) REFERENCES crm_integrations(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Scheduled exports table
CREATE TABLE IF NOT EXISTS scheduled_exports (
    id TEXT PRIMARY KEY,
    brand_id TEXT NOT NULL,
    name TEXT NOT NULL,
    data_type TEXT NOT NULL,
    format TEXT NOT NULL,
    schedule_cron TEXT NOT NULL, -- Cron expression
    event_filter TEXT, -- JSON filter criteria
    last_run DATETIME,
    next_run DATETIME,
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'paused', 'error')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Export templates table
CREATE TABLE IF NOT EXISTS export_templates (
    id TEXT PRIMARY KEY,
    brand_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    data_type TEXT NOT NULL,
    format TEXT NOT NULL,
    columns TEXT NOT NULL, -- JSON array of column definitions
    filters TEXT, -- JSON filter criteria
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

-- Data access logs table for compliance
CREATE TABLE IF NOT EXISTS data_access_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    brand_id TEXT NOT NULL,
    user_id TEXT,
    action TEXT NOT NULL CHECK (action IN ('export', 'view', 'download', 'sync')),
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_export_requests_brand ON export_requests(brand_id);
CREATE INDEX IF NOT EXISTS idx_export_requests_status ON export_requests(status);
CREATE INDEX IF NOT EXISTS idx_export_requests_created ON export_requests(created_at);
CREATE INDEX IF NOT EXISTS idx_export_requests_expires ON export_requests(expires_at);
CREATE INDEX IF NOT EXISTS idx_crm_integrations_brand ON crm_integrations(brand_id);
CREATE INDEX IF NOT EXISTS idx_crm_integrations_type ON crm_integrations(crm_type);
CREATE INDEX IF NOT EXISTS idx_crm_sync_logs_integration ON crm_sync_logs(integration_id);
CREATE INDEX IF NOT EXISTS idx_crm_sync_logs_started ON crm_sync_logs(started_at);
CREATE INDEX IF NOT EXISTS idx_scheduled_exports_brand ON scheduled_exports(brand_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_exports_next_run ON scheduled_exports(next_run);
CREATE INDEX IF NOT EXISTS idx_export_templates_brand ON export_templates(brand_id);
CREATE INDEX IF NOT EXISTS idx_data_access_logs_brand ON data_access_logs(brand_id);
CREATE INDEX IF NOT EXISTS idx_data_access_logs_created ON data_access_logs(created_at);

-- Insert sample export templates
INSERT OR IGNORE INTO export_templates (id, brand_id, name, description, data_type, format, columns) VALUES
('template_attendance', 'brand_1', 'Basic Attendance Export', 'Standard attendance data export', 'attendance', 'csv', '["user_id", "email", "checkin_time", "checkout_time", "duration"]'),
('template_content', 'brand_1', 'Content Analytics Export', 'Content performance and engagement data', 'content', 'csv', '["content_id", "user_id", "media_type", "caption", "view_count", "share_count", "created_at"]'),
('template_analytics', 'brand_1', 'Event Analytics Summary', 'High-level event analytics and KPIs', 'analytics', 'json', '["metric", "value", "event_id", "generated_at"]');

-- Insert sample scheduled exports
INSERT OR IGNORE INTO scheduled_exports (id, brand_id, name, data_type, format, schedule_cron, next_run) VALUES
('scheduled_1', 'brand_1', 'Weekly Attendance Report', 'attendance', 'csv', '0 9 * * 1', datetime('now', '+7 days')),
('scheduled_2', 'brand_1', 'Monthly Analytics Summary', 'analytics', 'json', '0 9 1 * *', datetime('now', '+1 month'));

-- Insert sample export requests
INSERT OR IGNORE INTO export_requests (id, brand_id, event_id, data_type, format, status, file_url, record_count, expires_at) VALUES
('export_1', 'brand_1', 'event_1', 'attendance', 'csv', 'completed', 'https://exports.lynkr.com/export_1.csv', 1247, datetime('now', '+7 days')),
('export_2', 'brand_1', 'event_1', 'content', 'json', 'completed', 'https://exports.lynkr.com/export_2.json', 89, datetime('now', '+7 days'));

-- Insert sample CRM integrations
INSERT OR IGNORE INTO crm_integrations (id, brand_id, crm_type, api_key, sync_interval, status) VALUES
('crm_1', 'brand_1', 'salesforce', 'sf_api_key_123', 60, 'active'),
('crm_2', 'brand_1', 'hubspot', 'hs_api_key_456', 120, 'active');