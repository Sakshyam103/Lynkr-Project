/**
 * Export Service
 * Handles data export functionality for brands
 */

package services

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ExportService struct {
	db *sql.DB
}

type ExportRequest struct {
	ID        string    `json:"id"`
	BrandID   string    `json:"brandId"`
	EventID   string    `json:"eventId"`
	DataType  string    `json:"dataType"`
	Format    string    `json:"format"`
	Status    string    `json:"status"`
	FileURL   string    `json:"fileUrl"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewExportService(db *sql.DB) *ExportService {
	return &ExportService{db: db}
}

func (es *ExportService) CreateExportRequest(brandID, eventID, dataType, format string) (*ExportRequest, error) {
	requestID := fmt.Sprintf("export_%d", time.Now().UnixNano())

	query := `
		INSERT INTO export_requests (id, export_request_id, brand_id, event_id, data_type, format, status, created_at)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // Expires in 7 days

	_, err := es.db.Exec(query, requestID, requestID, brandID, eventID, dataType, format, "processing", now)
	if err != nil {
		return nil, fmt.Errorf("failed to create export request: %w", err)
	}

	// Process export asynchronously
	go es.processExport(requestID, brandID, eventID, dataType, format)

	return &ExportRequest{
		ID:        requestID,
		BrandID:   brandID,
		EventID:   eventID,
		DataType:  dataType,
		Format:    format,
		Status:    "processing",
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}, nil
}

func (es *ExportService) processExport(requestID, brandID, eventID, dataType, format string) {
	var data []map[string]interface{}
	var err error

	switch dataType {
	case "attendance":
		data, err = es.exportAttendanceData(eventID)
	case "content":
		data, err = es.exportContentData(eventID)
	case "analytics":
		data, err = es.exportAnalyticsData(eventID)
	case "feedback":
		data, err = es.exportFeedbackData(eventID)
	default:
		err = fmt.Errorf("unsupported data type: %s", dataType)
	}

	if err != nil {
		es.updateExportStatus(requestID, "failed", "")
		return
	}
	var fileContent string
	switch format {
	case "csv":
		fileContent = es.convertToCSV(data)
	case "json":
		fileContent = es.convertToJSON(data)
	default:
		es.updateExportStatus(requestID, "failed", "")
		return
	}
	fmt.Print(fileContent)
	// In production, would upload to cloud storage
	fileURL := fmt.Sprintf("https://exports.lynkr.com/%s.%s", requestID, format)
	es.updateExportStatus(requestID, "completed", fileURL)
}

func (es *ExportService) exportAttendanceData(eventID string) ([]map[string]interface{}, error) {
	query := `
		SELECT u.id, u.email, a.checkin_time, a.checkout_time
		FROM attendances a
		JOIN users u ON a.user_id = u.id
		WHERE a.event_id = ?
	`

	rows, err := es.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var userID, email string
		var checkinTime, checkoutTime sql.NullString

		err := rows.Scan(&userID, &email, &checkinTime, &checkoutTime)
		if err != nil {
			continue
		}

		record := map[string]interface{}{
			"user_id":       userID,
			"email":         email,
			"checkin_time":  checkinTime.String,
			"checkout_time": checkoutTime.String,
		}
		data = append(data, record)
	}

	return data, nil
}

func (es *ExportService) exportContentData(eventID string) ([]map[string]interface{}, error) {
	query := `
		SELECT c.id, c.user_id, c.media_type, c.caption, c.view_count, c.share_count, c.created_at
		FROM content c
		WHERE c.event_id = ?
	`

	rows, err := es.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var contentID, userID, mediaType, caption string
		var viewCount, shareCount int
		var createdAt string

		err := rows.Scan(&contentID, &userID, &mediaType, &caption, &viewCount, &shareCount, &createdAt)
		if err != nil {
			continue
		}

		record := map[string]interface{}{
			"content_id":  contentID,
			"user_id":     userID,
			"media_type":  mediaType,
			"caption":     caption,
			"view_count":  viewCount,
			"share_count": shareCount,
			"created_at":  createdAt,
		}
		data = append(data, record)
	}

	return data, nil
}

func (es *ExportService) exportAnalyticsData(eventID string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(DISTINCT a.user_id) as total_attendees,
			COUNT(DISTINCT c.id) as total_content,
			SUM(c.view_count) as total_views,
			SUM(c.share_count) as total_shares
		FROM events e
		LEFT JOIN attendances a ON e.id = a.event_id
		LEFT JOIN content c ON e.id = c.event_id
		WHERE e.id = ?
	`

	var totalAttendees, totalContent, totalViews, totalShares int
	err := es.db.QueryRow(query, eventID).Scan(&totalAttendees, &totalContent, &totalViews, &totalShares)
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{
		{
			"metric":       "total_attendees",
			"value":        totalAttendees,
			"event_id":     eventID,
			"generated_at": time.Now().Format(time.RFC3339),
		},
		{
			"metric":       "total_content",
			"value":        totalContent,
			"event_id":     eventID,
			"generated_at": time.Now().Format(time.RFC3339),
		},
		{
			"metric":       "total_views",
			"value":        totalViews,
			"event_id":     eventID,
			"generated_at": time.Now().Format(time.RFC3339),
		},
		{
			"metric":       "total_shares",
			"value":        totalShares,
			"event_id":     eventID,
			"generated_at": time.Now().Format(time.RFC3339),
		},
	}

	return data, nil
}

func (es *ExportService) exportFeedbackData(eventID string) ([]map[string]interface{}, error) {
	query := `
		SELECT pv.user_id, pv.poll_id, pv.option_id, pv.created_at
		FROM poll_votes pv
		JOIN polls p ON pv.poll_id = p.id
		WHERE p.event_id = ?
	`

	rows, err := es.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var userID, pollID, optionID, createdAt string

		err := rows.Scan(&userID, &pollID, &optionID, &createdAt)
		if err != nil {
			continue
		}

		record := map[string]interface{}{
			"user_id":    userID,
			"poll_id":    pollID,
			"option_id":  optionID,
			"created_at": createdAt,
		}
		data = append(data, record)
	}

	return data, nil
}

func (es *ExportService) convertToCSV(data []map[string]interface{}) string {
	if len(data) == 0 {
		return ""
	}

	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	// Write header
	var headers []string
	for key := range data[0] {
		headers = append(headers, key)
	}
	writer.Write(headers)

	// Write data
	for _, record := range data {
		var row []string
		for _, header := range headers {
			row = append(row, fmt.Sprintf("%v", record[header]))
		}
		writer.Write(row)
	}

	writer.Flush()
	return builder.String()
}

func (es *ExportService) convertToJSON(data []map[string]interface{}) string {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return string(jsonData)
}

func (es *ExportService) updateExportStatus(requestID, status, fileURL string) {
	query := `UPDATE export_requests SET status = ?, file_url = ? WHERE id = ?`
	es.db.Exec(query, status, fileURL, requestID)
}

func (es *ExportService) GetExportStatus(requestID string) (*ExportRequest, error) {
	query := `
		SELECT id, brand_id, event_id, data_type, format, status, file_url, created_at, expires_at
		FROM export_requests WHERE id = ?
	`

	var req ExportRequest
	err := es.db.QueryRow(query, requestID).Scan(
		&req.ID, &req.BrandID, &req.EventID, &req.DataType,
		&req.Format, &req.Status, &req.FileURL, &req.CreatedAt, &req.ExpiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf("export request not found: %w", err)
	}

	return &req, nil
}
