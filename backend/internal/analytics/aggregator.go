/**
 * Analytics Aggregator
 * Batch processing for analytics data aggregation
 */

package analytics

import (
	"database/sql"
	"fmt"
	"time"
)

type Aggregator struct {
	db *sql.DB
}

type EngagementMetrics struct {
	EventID         string  `json:"eventId"`
	TotalAttendees  int     `json:"totalAttendees"`
	ContentPieces   int     `json:"contentPieces"`
	EngagementRate  float64 `json:"engagementRate"`
	AverageRating   float64 `json:"averageRating"`
}

type AttendanceAnalytics struct {
	EventID    string    `json:"eventId"`
	Date       time.Time `json:"date"`
	Attendees  int       `json:"attendees"`
	CheckIns   int       `json:"checkIns"`
	Duration   float64   `json:"duration"`
}

func NewAggregator(db *sql.DB) *Aggregator {
	return &Aggregator{db: db}
}

func (a *Aggregator) StartBatchProcessing() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			a.aggregateEngagementMetrics()
			a.aggregateAttendanceData()
			a.aggregateContentMetrics()
		}
	}()
}

func (a *Aggregator) aggregateEngagementMetrics() {
	query := `
		INSERT OR REPLACE INTO engagement_summary (event_id, total_attendees, content_pieces, engagement_rate, updated_at)
		SELECT 
			e.id,
			COUNT(DISTINCT a.user_id) as total_attendees,
			COUNT(DISTINCT c.id) as content_pieces,
			COALESCE(AVG(sf.value), 0) as engagement_rate,
			?
		FROM events e
		LEFT JOIN attendances a ON e.id = a.event_id
		LEFT JOIN content c ON e.id = c.event_id
		LEFT JOIN slider_feedback sf ON e.id = sf.event_id
		GROUP BY e.id
	`
	a.db.Exec(query, time.Now())
}

func (a *Aggregator) aggregateAttendanceData() {
	query := `
		INSERT OR REPLACE INTO attendance_summary (event_id, date, attendees, check_ins, updated_at)
		SELECT 
			event_id,
			DATE(created_at) as date,
			COUNT(DISTINCT user_id) as attendees,
			COUNT(*) as check_ins,
			?
		FROM attendances
		WHERE created_at >= DATE('now', '-7 days')
		GROUP BY event_id, DATE(created_at)
	`
	a.db.Exec(query, time.Now())
}

func (a *Aggregator) aggregateContentMetrics() {
	query := `
		INSERT OR REPLACE INTO content_summary (event_id, total_content, total_views, total_shares, updated_at)
		SELECT 
			event_id,
			COUNT(*) as total_content,
			SUM(view_count) as total_views,
			SUM(share_count) as total_shares,
			?
		FROM content
		WHERE event_id IS NOT NULL
		GROUP BY event_id
	`
	a.db.Exec(query, time.Now())
}

func (a *Aggregator) GetEngagementMetrics(eventID string) (*EngagementMetrics, error) {
	query := `
		SELECT event_id, total_attendees, content_pieces, engagement_rate
		FROM engagement_summary WHERE event_id = ?
	`
	
	var metrics EngagementMetrics
	err := a.db.QueryRow(query, eventID).Scan(
		&metrics.EventID, &metrics.TotalAttendees,
		&metrics.ContentPieces, &metrics.EngagementRate,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement metrics: %w", err)
	}
	
	return &metrics, nil
}

func (a *Aggregator) GetAttendanceAnalytics(eventID string) ([]AttendanceAnalytics, error) {
	query := `
		SELECT event_id, date, attendees, check_ins
		FROM attendance_summary 
		WHERE event_id = ?
		ORDER BY date DESC
		LIMIT 30
	`
	
	rows, err := a.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance analytics: %w", err)
	}
	defer rows.Close()
	
	var analytics []AttendanceAnalytics
	for rows.Next() {
		var a AttendanceAnalytics
		var dateStr string
		
		err := rows.Scan(&a.EventID, &dateStr, &a.Attendees, &a.CheckIns)
		if err != nil {
			continue
		}
		
		a.Date, _ = time.Parse("2006-01-02", dateStr)
		analytics = append(analytics, a)
	}
	
	return analytics, nil
}