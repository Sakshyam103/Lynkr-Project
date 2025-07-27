-- Performance Optimization Migration
-- Adds performance monitoring tables and optimized indexes

-- Query performance logs table
CREATE TABLE IF NOT EXISTS query_performance_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query_hash TEXT NOT NULL,
    query_text TEXT NOT NULL,
    execution_time REAL NOT NULL,
    rows_examined INTEGER DEFAULT 0,
    rows_returned INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- System performance metrics table
CREATE TABLE IF NOT EXISTS system_performance_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    metric_name TEXT NOT NULL,
    metric_value REAL NOT NULL,
    metric_unit TEXT,
    recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Cache performance table
CREATE TABLE IF NOT EXISTS cache_performance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cache_key TEXT NOT NULL,
    hit_count INTEGER DEFAULT 0,
    miss_count INTEGER DEFAULT 0,
    last_accessed DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Load test results table
CREATE TABLE IF NOT EXISTS load_test_results (
    id TEXT PRIMARY KEY,
    endpoint TEXT NOT NULL,
    concurrency INTEGER NOT NULL,
    duration INTEGER NOT NULL,
    total_requests INTEGER NOT NULL,
    successful_requests INTEGER NOT NULL,
    failed_requests INTEGER NOT NULL,
    avg_response_time REAL NOT NULL,
    min_response_time REAL NOT NULL,
    max_response_time REAL NOT NULL,
    requests_per_second REAL NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Performance-optimized indexes
CREATE INDEX IF NOT EXISTS idx_users_email_active ON users(email) WHERE active = 1;
CREATE INDEX IF NOT EXISTS idx_events_brand_date ON events(brand_id, start_date);
CREATE INDEX IF NOT EXISTS idx_attendances_event_user_time ON attendances(event_id, user_id, checkin_time);
CREATE INDEX IF NOT EXISTS idx_content_event_created_views ON content(event_id, created_at DESC, view_count DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_events_type_event_time ON analytics_events(type, event_id, created_at);
CREATE INDEX IF NOT EXISTS idx_engagement_metrics_event_user_time ON engagement_metrics(event_id, user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_purchases_event_user_time ON purchases(event_id, user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_rewards_user_status_time ON rewards(user_id, status, created_at);
CREATE INDEX IF NOT EXISTS idx_export_requests_brand_status ON export_requests(brand_id, status);
CREATE INDEX IF NOT EXISTS idx_crm_integrations_brand_status ON crm_integrations(brand_id, status);

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_content_analytics_composite ON content_analytics(content_id, action, created_at);
CREATE INDEX IF NOT EXISTS idx_poll_votes_composite ON poll_votes(poll_id, user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_sentiment_analysis_composite ON sentiment_analysis(content_id, created_at);
CREATE INDEX IF NOT EXISTS idx_pixel_events_composite ON pixel_events(event_id, brand_id, event_type, created_at);

-- Partial indexes for better performance on filtered queries
CREATE INDEX IF NOT EXISTS idx_content_published ON content(event_id, created_at) WHERE status = 'published';
CREATE INDEX IF NOT EXISTS idx_attendances_active ON attendances(event_id, user_id) WHERE checkout_time IS NULL;
CREATE INDEX IF NOT EXISTS idx_rewards_pending ON rewards(user_id, created_at) WHERE status = 'pending';

-- Covering indexes to avoid table lookups
CREATE INDEX IF NOT EXISTS idx_events_brand_covering ON events(brand_id) INCLUDE (name, start_date, end_date, status);
CREATE INDEX IF NOT EXISTS idx_users_auth_covering ON users(email) INCLUDE (id, name, active, created_at);

-- Performance monitoring views
CREATE VIEW IF NOT EXISTS slow_queries AS
SELECT 
    query_hash,
    query_text,
    AVG(execution_time) as avg_time,
    COUNT(*) as execution_count,
    MAX(execution_time) as max_time
FROM query_performance_logs 
WHERE created_at > datetime('now', '-1 day')
GROUP BY query_hash
HAVING avg_time > 100
ORDER BY avg_time DESC;

CREATE VIEW IF NOT EXISTS cache_hit_rates AS
SELECT 
    cache_key,
    hit_count,
    miss_count,
    ROUND((CAST(hit_count AS REAL) / (hit_count + miss_count)) * 100, 2) as hit_rate
FROM cache_performance
WHERE hit_count + miss_count > 0
ORDER BY hit_rate ASC;

CREATE VIEW IF NOT EXISTS performance_summary AS
SELECT 
    'database' as component,
    'query_performance' as metric,
    AVG(execution_time) as value,
    'ms' as unit,
    datetime('now') as recorded_at
FROM query_performance_logs 
WHERE created_at > datetime('now', '-1 hour')
UNION ALL
SELECT 
    'cache' as component,
    'hit_rate' as metric,
    AVG(CAST(hit_count AS REAL) / (hit_count + miss_count) * 100) as value,
    '%' as unit,
    datetime('now') as recorded_at
FROM cache_performance
WHERE hit_count + miss_count > 0;

-- Indexes for performance monitoring tables
CREATE INDEX IF NOT EXISTS idx_query_performance_hash_time ON query_performance_logs(query_hash, created_at);
CREATE INDEX IF NOT EXISTS idx_system_metrics_name_time ON system_performance_metrics(metric_name, recorded_at);
CREATE INDEX IF NOT EXISTS idx_cache_performance_key ON cache_performance(cache_key);
CREATE INDEX IF NOT EXISTS idx_load_test_endpoint_time ON load_test_results(endpoint, created_at);

-- Insert sample performance data
INSERT OR IGNORE INTO system_performance_metrics (metric_name, metric_value, metric_unit) VALUES
('cpu_usage', 45.2, '%'),
('memory_usage', 68.7, '%'),
('disk_usage', 23.1, '%'),
('response_time', 125.5, 'ms');

INSERT OR IGNORE INTO cache_performance (cache_key, hit_count, miss_count) VALUES
('event:event_1', 1250, 45),
('user:user_1', 890, 23),
('analytics:daily_stats', 567, 12);

-- Database optimization settings
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 10000;
PRAGMA temp_store = MEMORY;
PRAGMA mmap_size = 268435456;