/**
 * Pixel Tracking Service
 * Handles pixel tracking for brand-related searches and website visits
 */

package services

import (
	"database/sql"
	"fmt"
	"time"
)

type PixelEvent struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	EventID   string    `json:"eventId"`
	BrandID   string    `json:"brandId"`
	EventType string    `json:"eventType"`
	URL       string    `json:"url"`
	Referrer  string    `json:"referrer"`
	UserAgent string    `json:"userAgent"`
	CreatedAt time.Time `json:"createdAt"`
}

type PixelService struct {
	db *sql.DB
}

func NewPixelService(db *sql.DB) *PixelService {
	return &PixelService{db: db}
}

func (ps *PixelService) TrackPixelEvent(userID, eventID, brandID, eventType, url, referrer, userAgent string) (*PixelEvent, error) {
	pixelID := fmt.Sprintf("pixel_%d", time.Now().UnixNano())

	query := `
		INSERT INTO pixel_events (id, user_id, event_id, brand_id, event_type, url, referrer, user_agent, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err := ps.db.Exec(query, pixelID, userID, eventID, brandID, eventType, url, referrer, userAgent, now)
	if err != nil {
		return nil, fmt.Errorf("failed to track pixel event: %w", err)
	}

	return &PixelEvent{
		ID:        pixelID,
		UserID:    userID,
		EventID:   eventID,
		BrandID:   brandID,
		EventType: eventType,
		URL:       url,
		Referrer:  referrer,
		UserAgent: userAgent,
		CreatedAt: now,
	}, nil
}

func (ps *PixelService) GetPixelAnalytics(eventID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			event_type,
			COUNT(*) as count,
			COUNT(DISTINCT user_id) as unique_users
		FROM pixel_events 
		WHERE event_id = ?
		GROUP BY event_type
	`

	rows, err := ps.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pixel analytics: %w", err)
	}
	defer rows.Close()

	analytics := make(map[string]interface{})
	totalEvents := 0
	// uniqueUsers := make(map[string]bool)

	for rows.Next() {
		var eventType string
		var count, users int

		err := rows.Scan(&eventType, &count, &users)
		if err != nil {
			continue
		}

		analytics[eventType] = map[string]int{
			"count":        count,
			"unique_users": users,
		}
		totalEvents += count
	}

	// Get total unique users across all events
	userQuery := `SELECT COUNT(DISTINCT user_id) FROM pixel_events WHERE event_id = ?`
	var totalUniqueUsers int
	ps.db.QueryRow(userQuery, eventID).Scan(&totalUniqueUsers)

	analytics["summary"] = map[string]int{
		"total_events": totalEvents,
		"unique_users": totalUniqueUsers,
	}

	return analytics, nil
}

func (ps *PixelService) GeneratePixelURL(eventID, brandID string) string {
	return fmt.Sprintf("https://api.lynkr.com/pixel/track?event=%s&brand=%s", eventID, brandID)
}
