/**
 * Usability Tester
 * Tools for conducting usability testing and UX analysis
 */

package ux

import (
	"database/sql"
	// "encoding/json"
	"fmt"
	"time"
)

type UsabilityTester struct {
	db *sql.DB
}

type UserSession struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	SessionStart time.Time `json:"sessionStart"`
	SessionEnd   time.Time `json:"sessionEnd"`
	Duration     int       `json:"duration"`
	Actions      []Action  `json:"actions"`
	Errors       []Error   `json:"errors"`
}

type Action struct {
	Type      string    `json:"type"`
	Screen    string    `json:"screen"`
	Element   string    `json:"element"`
	Timestamp time.Time `json:"timestamp"`
	Duration  int       `json:"duration"`
}

type Error struct {
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Screen      string    `json:"screen"`
	Timestamp   time.Time `json:"timestamp"`
	Recoverable bool      `json:"recoverable"`
}

type UsabilityMetrics struct {
	TaskCompletionRate   float64 `json:"taskCompletionRate"`
	AverageTaskTime      float64 `json:"averageTaskTime"`
	ErrorRate            float64 `json:"errorRate"`
	UserSatisfaction     float64 `json:"userSatisfaction"`
	NavigationEfficiency float64 `json:"navigationEfficiency"`
}

func NewUsabilityTester(db *sql.DB) *UsabilityTester {
	return &UsabilityTester{db: db}
}

func (ut *UsabilityTester) StartSession(userID string) (*UserSession, error) {
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())

	session := &UserSession{
		ID:           sessionID,
		UserID:       userID,
		SessionStart: time.Now(),
		Actions:      []Action{},
		Errors:       []Error{},
	}

	query := `
		INSERT INTO user_sessions (id, user_id, session_start, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := ut.db.Exec(query, sessionID, userID, session.SessionStart, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}

	return session, nil
}

func (ut *UsabilityTester) TrackAction(sessionID, actionType, screen, element string, duration int) error {
	query := `
		INSERT INTO user_actions (session_id, type, screen, element, duration, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := ut.db.Exec(query, sessionID, actionType, screen, element, duration, time.Now())
	return err
}

func (ut *UsabilityTester) TrackError(sessionID, errorType, message, screen string, recoverable bool) error {
	query := `
		INSERT INTO user_errors (session_id, type, message, screen, recoverable, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := ut.db.Exec(query, sessionID, errorType, message, screen, recoverable, time.Now())
	return err
}

func (ut *UsabilityTester) EndSession(sessionID string) error {
	query := `UPDATE user_sessions SET session_end = ?, duration = ? WHERE id = ?`

	now := time.Now()

	// Get session start time
	var sessionStart time.Time
	ut.db.QueryRow("SELECT session_start FROM user_sessions WHERE id = ?", sessionID).Scan(&sessionStart)

	duration := int(now.Sub(sessionStart).Seconds())

	_, err := ut.db.Exec(query, now, duration, sessionID)
	return err
}

func (ut *UsabilityTester) GetUsabilityMetrics(timeframe string) (*UsabilityMetrics, error) {
	var whereClause string
	switch timeframe {
	case "24h":
		whereClause = "WHERE created_at > datetime('now', '-1 day')"
	case "7d":
		whereClause = "WHERE created_at > datetime('now', '-7 days')"
	case "30d":
		whereClause = "WHERE created_at > datetime('now', '-30 days')"
	default:
		whereClause = ""
	}

	metrics := &UsabilityMetrics{}

	// Task completion rate
	completionQuery := fmt.Sprintf(`
		SELECT 
			COUNT(CASE WHEN session_end IS NOT NULL THEN 1 END) * 100.0 / COUNT(*) as completion_rate
		FROM user_sessions %s
	`, whereClause)
	ut.db.QueryRow(completionQuery).Scan(&metrics.TaskCompletionRate)

	// Average task time
	timeQuery := fmt.Sprintf(`
		SELECT AVG(duration) FROM user_sessions 
		WHERE session_end IS NOT NULL %s
	`, whereClause)
	ut.db.QueryRow(timeQuery).Scan(&metrics.AverageTaskTime)

	// Error rate
	errorQuery := fmt.Sprintf(`
		SELECT 
			COUNT(ue.id) * 100.0 / COUNT(DISTINCT us.id) as error_rate
		FROM user_sessions us
		LEFT JOIN user_errors ue ON us.id = ue.session_id
		%s
	`, whereClause)
	ut.db.QueryRow(errorQuery).Scan(&metrics.ErrorRate)

	// User satisfaction (from feedback)
	satisfactionQuery := fmt.Sprintf(`
		SELECT AVG(rating) FROM app_feedback 
		WHERE type = 'satisfaction' %s
	`, whereClause)
	ut.db.QueryRow(satisfactionQuery).Scan(&metrics.UserSatisfaction)

	// Navigation efficiency
	navQuery := fmt.Sprintf(`
		SELECT 
			AVG(CASE WHEN type = 'navigation' THEN duration END) as nav_efficiency
		FROM user_actions ua
		JOIN user_sessions us ON ua.session_id = us.id
		%s
	`, whereClause)
	ut.db.QueryRow(navQuery).Scan(&metrics.NavigationEfficiency)

	return metrics, nil
}

func (ut *UsabilityTester) GetHeatmapData(screen string) ([]map[string]interface{}, error) {
	query := `
		SELECT element, COUNT(*) as clicks, AVG(duration) as avg_duration
		FROM user_actions 
		WHERE screen = ? AND type = 'tap'
		GROUP BY element
		ORDER BY clicks DESC
	`

	rows, err := ut.db.Query(query, screen)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var heatmapData []map[string]interface{}
	for rows.Next() {
		var element string
		var clicks int
		var avgDuration float64

		if rows.Scan(&element, &clicks, &avgDuration) == nil {
			heatmapData = append(heatmapData, map[string]interface{}{
				"element":     element,
				"clicks":      clicks,
				"avgDuration": avgDuration,
			})
		}
	}

	return heatmapData, nil
}

func (ut *UsabilityTester) GetUserJourney(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT ua.screen, ua.type, ua.element, ua.created_at, ua.duration
		FROM user_actions ua
		JOIN user_sessions us ON ua.session_id = us.id
		WHERE us.user_id = ?
		ORDER BY ua.created_at DESC
		LIMIT 50
	`

	rows, err := ut.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journey []map[string]interface{}
	for rows.Next() {
		var screen, actionType, element, timestamp string
		var duration int

		if rows.Scan(&screen, &actionType, &element, &timestamp, &duration) == nil {
			journey = append(journey, map[string]interface{}{
				"screen":    screen,
				"action":    actionType,
				"element":   element,
				"timestamp": timestamp,
				"duration":  duration,
			})
		}
	}

	return journey, nil
}

func (ut *UsabilityTester) IdentifyPainPoints() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			screen,
			element,
			COUNT(*) as error_count,
			AVG(CASE WHEN recoverable = 0 THEN 1 ELSE 0 END) as critical_rate
		FROM user_errors
		WHERE created_at > datetime('now', '-7 days')
		GROUP BY screen, element
		HAVING error_count > 5
		ORDER BY error_count DESC, critical_rate DESC
	`

	rows, err := ut.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var painPoints []map[string]interface{}
	for rows.Next() {
		var screen, element string
		var errorCount int
		var criticalRate float64

		if rows.Scan(&screen, &element, &errorCount, &criticalRate) == nil {
			painPoints = append(painPoints, map[string]interface{}{
				"screen":       screen,
				"element":      element,
				"errorCount":   errorCount,
				"criticalRate": criticalRate,
				"severity":     ut.calculateSeverity(errorCount, criticalRate),
			})
		}
	}

	return painPoints, nil
}

func (ut *UsabilityTester) calculateSeverity(errorCount int, criticalRate float64) string {
	score := float64(errorCount) * (1 + criticalRate)

	if score > 50 {
		return "critical"
	} else if score > 20 {
		return "high"
	} else if score > 10 {
		return "medium"
	}
	return "low"
}
