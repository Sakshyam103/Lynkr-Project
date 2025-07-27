/**
 * Analytics Service
 * Core analytics service with engagement and performance metrics
 */

package services

import (
	"database/sql"
	"fmt"
	"time"
)

type AnalyticsService struct {
	db *sql.DB
}

type EngagementMetrics struct {
	Views          int     `json:"views"`
	Likes          int     `json:"likes"`
	Shares         int     `json:"shares"`
	Comments       int     `json:"comments"`
	EngagementRate float64 `json:"engagementRate"`
}

type ContentPerformance struct {
	ContentID      string  `json:"contentId"`
	Views          int     `json:"views"`
	Engagement     int     `json:"engagement"`
	SentimentScore float64 `json:"sentimentScore"`
	ShareRate      float64 `json:"shareRate"`
}

func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

func (as *AnalyticsService) GetEngagementMetrics(eventID string) (*EngagementMetrics, error) {
	query := `
		SELECT 
			SUM(view_count) as views,
			COUNT(CASE WHEN action = 'like' THEN 1 END) as likes,
			SUM(share_count) as shares,
			COUNT(CASE WHEN action = 'comment' THEN 1 END) as comments
		FROM content c
		LEFT JOIN content_analytics ca ON c.id = ca.content_id
		WHERE c.event_id = ?
	`

	var metrics EngagementMetrics
	err := as.db.QueryRow(query, eventID).Scan(
		&metrics.Views, &metrics.Likes, &metrics.Shares, &metrics.Comments,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get engagement metrics: %w", err)
	}

	totalEngagement := metrics.Likes + metrics.Shares + metrics.Comments
	if metrics.Views > 0 {
		metrics.EngagementRate = float64(totalEngagement) / float64(metrics.Views) * 100
	}

	return &metrics, nil
}

func (as *AnalyticsService) GetAttendanceAnalytics(eventID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(DISTINCT user_id) as unique_attendees,
			COUNT(*) as total_checkins,
			AVG(JULIANDAY(check_out_time) - JULIANDAY(check_in_time)) * 24 as avg_duration_hours
		FROM attendances 
		WHERE event_id = ? AND check_in_time IS NOT NULL
	`

	var uniqueAttendees, totalCheckins int
	var avgDuration sql.NullFloat64

	err := as.db.QueryRow(query, eventID).Scan(&uniqueAttendees, &totalCheckins, &avgDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance analytics: %w", err)
	}

	duration := 0.0
	if avgDuration.Valid {
		duration = avgDuration.Float64
	}

	return map[string]interface{}{
		"uniqueAttendees": uniqueAttendees,
		"totalCheckins":   totalCheckins,
		"avgDuration":     duration,
		"returnRate":      float64(totalCheckins) / float64(uniqueAttendees),
	}, nil
}

func (as *AnalyticsService) GetContentPerformance(eventID string) ([]ContentPerformance, error) {
	query := `
		SELECT 
			c.id,
			c.view_count,
			COUNT(ca.id) as engagement_count,
			COALESCE(AVG(CASE WHEN sa.result LIKE '%score%' THEN 
				CAST(json_extract(sa.result, '$.score') AS REAL) 
			END), 0) as sentiment_score,
			CAST(c.share_count AS REAL) / NULLIF(c.view_count, 0) as share_rate
		FROM content c
		LEFT JOIN content_analytics ca ON c.id = ca.content_id
		LEFT JOIN sentiment_analysis sa ON c.id = sa.content_id
		WHERE c.event_id = ?
		GROUP BY c.id
		ORDER BY c.view_count DESC
		LIMIT 20
	`

	rows, err := as.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content performance: %w", err)
	}
	defer rows.Close()

	var performance []ContentPerformance
	for rows.Next() {
		var cp ContentPerformance
		var shareRate sql.NullFloat64

		err := rows.Scan(
			&cp.ContentID, &cp.Views, &cp.Engagement,
			&cp.SentimentScore, &shareRate,
		)
		if err != nil {
			continue
		}

		if shareRate.Valid {
			cp.ShareRate = shareRate.Float64
		}

		performance = append(performance, cp)
	}

	return performance, nil
}

func (as *AnalyticsService) TrackEvent(eventType, userID, eventID string, data map[string]interface{}) error {
	query := `
		INSERT INTO analytics_events (type, user_id, event_id, data, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	dataJSON := "{}"
	if len(data) > 0 {
		// Simple JSON serialization
		dataJSON = fmt.Sprintf(`{"action":"%v","value":%v}`, data["action"], data["value"])
	}

	_, err := as.db.Exec(query, eventType, userID, eventID, dataJSON, time.Now())
	if err != nil {
		return fmt.Errorf("failed to track event: %w", err)
	}

	return nil
}

func (as *AnalyticsService) GetRealtimeStats(eventID string) (map[string]interface{}, error) {
	// Get current active users (checked in within last hour)
	activeUsersQuery := `
		SELECT COUNT(DISTINCT user_id)
		FROM attendances 
		WHERE event_id = ? AND checkin_time > datetime('now', '-1 hour')
		AND (checkout_time IS NULL OR checkout_time > datetime('now', '-1 hour'))
	`

	var activeUsers int
	as.db.QueryRow(activeUsersQuery, eventID).Scan(&activeUsers)

	// Get recent content count (last hour)
	recentContentQuery := `
		SELECT COUNT(*)
		FROM content 
		WHERE event_id = ? AND created_at > datetime('now', '-1 hour')
	`

	var recentContent int
	as.db.QueryRow(recentContentQuery, eventID).Scan(&recentContent)

	return map[string]interface{}{
		"activeUsers":   activeUsers,
		"recentContent": recentContent,
		"timestamp":     time.Now(),
	}, nil
}
