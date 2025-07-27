/**
 * Conversion Funnel Service
 * Handles conversion tracking and ROI calculation
 */

package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type FunnelStage struct {
	Stage       string  `json:"stage"`
	Users       int     `json:"users"`
	Conversions int     `json:"conversions"`
	Rate        float64 `json:"rate"`
}

type ConversionFunnel struct {
	EventID    string        `json:"eventId"`
	BrandID    string        `json:"brandId"`
	Stages     []FunnelStage `json:"stages"`
	TotalUsers int           `json:"totalUsers"`
	Revenue    float64       `json:"revenue"`
	ROI        float64       `json:"roi"`
}

type ConversionFunnelService struct {
	db *sql.DB
}

func NewConversionFunnelService(db *sql.DB) *ConversionFunnelService {
	return &ConversionFunnelService{db: db}
}

func (cfs *ConversionFunnelService) GetConversionFunnel(eventID, brandID string) (*ConversionFunnel, error) {
	stages := []FunnelStage{
		{Stage: "attendance", Users: 0, Conversions: 0},
		{Stage: "content_view", Users: 0, Conversions: 0},
		{Stage: "engagement", Users: 0, Conversions: 0},
		{Stage: "website_visit", Users: 0, Conversions: 0},
		{Stage: "purchase", Users: 0, Conversions: 0},
	}

	// Get attendance
	attendanceQuery := `SELECT COUNT(DISTINCT user_id) FROM attendances WHERE event_id = ?`
	cfs.db.QueryRow(attendanceQuery, eventID).Scan(&stages[0].Users)
	stages[0].Conversions = stages[0].Users

	// Get content views
	contentQuery := `SELECT COUNT(DISTINCT user_id) FROM content_analytics WHERE content_id IN (SELECT id FROM content WHERE event_id = ?)`
	cfs.db.QueryRow(contentQuery, eventID).Scan(&stages[1].Users)
	stages[1].Conversions = stages[1].Users

	// Get engagement (feedback, polls, etc.)
	engagementQuery := `SELECT COUNT(DISTINCT user_id) FROM engagement_metrics WHERE event_id = ?`
	cfs.db.QueryRow(engagementQuery, eventID).Scan(&stages[2].Users)
	stages[2].Conversions = stages[2].Users

	// Get website visits
	visitQuery := `SELECT COUNT(DISTINCT user_id) FROM pixel_events WHERE event_id = ? AND event_type = 'website_visit'`
	cfs.db.QueryRow(visitQuery, eventID).Scan(&stages[3].Users)
	stages[3].Conversions = stages[3].Users

	// Get purchases
	purchaseQuery := `SELECT COUNT(DISTINCT user_id), COALESCE(SUM(amount), 0) FROM purchases WHERE event_id = ?`
	var revenue float64
	cfs.db.QueryRow(purchaseQuery, eventID).Scan(&stages[4].Users, &revenue)
	stages[4].Conversions = stages[4].Users

	// Calculate conversion rates
	totalUsers := stages[0].Users
	for i := range stages {
		if totalUsers > 0 {
			stages[i].Rate = float64(stages[i].Users) / float64(totalUsers) * 100
		}
	}

	// Calculate ROI (simplified)
	roi := 0.0
	if revenue > 0 {
		// Assume event cost of $10,000 for ROI calculation
		eventCost := 10000.0
		roi = ((revenue - eventCost) / eventCost) * 100
	}

	return &ConversionFunnel{
		EventID:    eventID,
		BrandID:    brandID,
		Stages:     stages,
		TotalUsers: totalUsers,
		Revenue:    revenue,
		ROI:        roi,
	}, nil
}

func (cfs *ConversionFunnelService) GetAttributionReport(eventID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			pa.attribution_type,
			COUNT(*) as count,
			SUM(p.amount) as revenue
		FROM purchase_attribution pa
		JOIN purchases p ON pa.purchase_id = p.id
		WHERE pa.event_id = ?
		GROUP BY pa.attribution_type
	`

	rows, err := cfs.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribution report: %w", err)
	}
	defer rows.Close()

	attribution := make(map[string]map[string]interface{})
	totalRevenue := 0.0
	totalPurchases := 0

	for rows.Next() {
		var attrType string
		var count int
		var revenue float64

		err := rows.Scan(&attrType, &count, &revenue)
		if err != nil {
			continue
		}

		attribution[attrType] = map[string]interface{}{
			"purchases": count,
			"revenue":   revenue,
		}

		totalRevenue += revenue
		totalPurchases += count
	}

	// Calculate percentages
	// for _, data := range attribution {
	// 	if totalRevenue > 0 {
	// 		data.(map[string]interface{})["revenuePercent"] = (data.(map[string]interface{})["revenue"].(float64) / totalRevenue) * 100
	// 	}
	// 	if totalPurchases > 0 {
	// 		data.(map[string]interface{})["purchasePercent"] = (float64(data.(map[string]interface{})["purchases"].(int)) / float64(totalPurchases)) * 100
	// 	}
	// }

	return map[string]interface{}{
		"attribution":    attribution,
		"totalRevenue":   totalRevenue,
		"totalPurchases": totalPurchases,
	}, nil
}

func (cfs *ConversionFunnelService) TrackConversion(userID, eventID, stage string, metadata map[string]interface{}) error {
	metadataJSON := "{}"
	if len(metadata) > 0 {
		data, _ := json.Marshal(metadata)
		metadataJSON = string(data)
	}

	query := `
		INSERT INTO conversion_tracking (user_id, event_id, stage, metadata, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := cfs.db.Exec(query, userID, eventID, stage, metadataJSON, time.Now())
	if err != nil {
		return fmt.Errorf("failed to track conversion: %w", err)
	}

	return nil
}
