/**
 * Privacy Enhancer
 * Enhanced privacy features and data anonymization
 */

package security

import (
	// "crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"

	// "strings"
	"time"
)

type PrivacyEnhancer struct {
	db *sql.DB
}

type ConsentRecord struct {
	UserID      string    `json:"userId"`
	ConsentType string    `json:"consentType"`
	Granted     bool      `json:"granted"`
	Version     string    `json:"version"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type DataRetentionPolicy struct {
	DataType      string    `json:"dataType"`
	RetentionDays int       `json:"retentionDays"`
	AutoDelete    bool      `json:"autoDelete"`
	LastCleanup   time.Time `json:"lastCleanup"`
}

func NewPrivacyEnhancer(db *sql.DB) *PrivacyEnhancer {
	return &PrivacyEnhancer{db: db}
}

func (pe *PrivacyEnhancer) UpdateConsentFlow() error {
	// Update consent management with enhanced granular controls
	consentTypes := []string{
		"analytics_tracking",
		"marketing_communications",
		"data_sharing_partners",
		"location_tracking",
		"content_usage_rights",
		"personalized_recommendations",
	}

	for _, consentType := range consentTypes {
		query := `
			INSERT OR IGNORE INTO consent_types (type, description, required, version, created_at)
			VALUES (?, ?, ?, ?, ?)
		`

		description := pe.getConsentDescription(consentType)
		required := pe.isConsentRequired(consentType)

		pe.db.Exec(query, consentType, description, required, "2.0", time.Now())
	}

	return nil
}

func (pe *PrivacyEnhancer) AnonymizeUserData(userID string) error {
	// Generate anonymous identifier
	hash := sha256.Sum256([]byte(userID + time.Now().String()))
	anonID := hex.EncodeToString(hash[:])[:16]

	// Anonymize user data
	queries := []string{
		`UPDATE users SET email = ?, name = ?, phone = NULL WHERE id = ?`,
		`UPDATE attendances SET user_id = ? WHERE user_id = ?`,
		`UPDATE content SET user_id = ? WHERE user_id = ?`,
		`UPDATE analytics_events SET user_id = ? WHERE user_id = ?`,
	}

	anonEmail := fmt.Sprintf("anon_%s@privacy.local", anonID)
	anonName := fmt.Sprintf("Anonymous_%s", anonID[:8])

	for i, query := range queries {
		if i == 0 {
			pe.db.Exec(query, anonEmail, anonName, userID)
		} else {
			pe.db.Exec(query, anonID, userID)
		}
	}

	// Log anonymization
	pe.logPrivacyAction("data_anonymized", userID, "User data anonymized upon request")

	return nil
}

func (pe *PrivacyEnhancer) ImplementDataRetention() error {
	policies := []DataRetentionPolicy{
		{DataType: "analytics_events", RetentionDays: 365, AutoDelete: true},
		{DataType: "security_events", RetentionDays: 90, AutoDelete: true},
		{DataType: "export_requests", RetentionDays: 30, AutoDelete: true},
		{DataType: "query_performance_logs", RetentionDays: 7, AutoDelete: true},
	}

	for _, policy := range policies {
		if policy.AutoDelete {
			pe.cleanupOldData(policy.DataType, policy.RetentionDays)
		}

		// Store policy
		query := `
			INSERT OR REPLACE INTO data_retention_policies 
			(data_type, retention_days, auto_delete, last_cleanup)
			VALUES (?, ?, ?, ?)
		`
		pe.db.Exec(query, policy.DataType, policy.RetentionDays, policy.AutoDelete, time.Now())
	}

	return nil
}

func (pe *PrivacyEnhancer) cleanupOldData(dataType string, retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	queries := map[string]string{
		"analytics_events":       `DELETE FROM analytics_events WHERE created_at < ?`,
		"security_events":        `DELETE FROM security_events WHERE created_at < ?`,
		"export_requests":        `DELETE FROM export_requests WHERE created_at < ? AND status = 'completed'`,
		"query_performance_logs": `DELETE FROM query_performance_logs WHERE created_at < ?`,
	}

	if query, exists := queries[dataType]; exists {
		result, err := pe.db.Exec(query, cutoffDate)
		if err == nil {
			rowsAffected, _ := result.RowsAffected()
			pe.logPrivacyAction("data_cleanup", "system", fmt.Sprintf("Cleaned up %d records from %s", rowsAffected, dataType))
		}
		return err
	}

	return fmt.Errorf("unknown data type: %s", dataType)
}

func (pe *PrivacyEnhancer) ProcessDataDeletionRequest(userID string) error {
	// Mark user for deletion
	query := `UPDATE users SET deletion_requested = 1, deletion_requested_at = ? WHERE id = ?`
	pe.db.Exec(query, time.Now(), userID)

	// Schedule data deletion (30-day grace period)
	deletionDate := time.Now().AddDate(0, 0, 30)
	scheduleQuery := `
		INSERT INTO scheduled_deletions (user_id, scheduled_for, status, created_at)
		VALUES (?, ?, ?, ?)
	`
	pe.db.Exec(scheduleQuery, userID, deletionDate, "scheduled", time.Now())

	pe.logPrivacyAction("deletion_requested", userID, "User requested data deletion")

	return nil
}

func (pe *PrivacyEnhancer) GetUserDataExport(userID string) (map[string]interface{}, error) {
	export := make(map[string]interface{})

	// User profile data
	userQuery := `SELECT id, email, name, created_at FROM users WHERE id = ?`
	var user map[string]interface{}
	row := pe.db.QueryRow(userQuery, userID)
	var id, email, name, createdAt string
	if row.Scan(&id, &email, &name, &createdAt) == nil {
		user = map[string]interface{}{
			"id": id, "email": email, "name": name, "created_at": createdAt,
		}
	}
	export["profile"] = user

	// Attendance data
	attendanceQuery := `SELECT event_id, checkin_time, checkout_time FROM attendances WHERE user_id = ?`
	rows, _ := pe.db.Query(attendanceQuery, userID)
	defer rows.Close()

	var attendances []map[string]interface{}
	for rows.Next() {
		var eventID, checkinTime, checkoutTime string
		if rows.Scan(&eventID, &checkinTime, &checkoutTime) == nil {
			attendances = append(attendances, map[string]interface{}{
				"event_id": eventID, "checkin_time": checkinTime, "checkout_time": checkoutTime,
			})
		}
	}
	export["attendances"] = attendances

	// Content data
	contentQuery := `SELECT id, media_type, caption, created_at FROM content WHERE user_id = ?`
	contentRows, _ := pe.db.Query(contentQuery, userID)
	defer contentRows.Close()

	var content []map[string]interface{}
	for contentRows.Next() {
		var id, mediaType, caption, createdAt string
		if contentRows.Scan(&id, &mediaType, &caption, &createdAt) == nil {
			content = append(content, map[string]interface{}{
				"id": id, "media_type": mediaType, "caption": caption, "created_at": createdAt,
			})
		}
	}
	export["content"] = content

	pe.logPrivacyAction("data_exported", userID, "User data exported")

	return export, nil
}

func (pe *PrivacyEnhancer) logPrivacyAction(action, userID, details string) {
	query := `
		INSERT INTO privacy_audit_log (action, user_id, details, created_at)
		VALUES (?, ?, ?, ?)
	`
	pe.db.Exec(query, action, userID, details, time.Now())
}

func (pe *PrivacyEnhancer) getConsentDescription(consentType string) string {
	descriptions := map[string]string{
		"analytics_tracking":           "Allow collection of usage analytics to improve the app",
		"marketing_communications":     "Receive marketing emails and promotional content",
		"data_sharing_partners":        "Share anonymized data with trusted partners",
		"location_tracking":            "Track location for event check-ins and geofencing",
		"content_usage_rights":         "Allow brands to use your content for marketing",
		"personalized_recommendations": "Receive personalized event and product recommendations",
	}
	return descriptions[consentType]
}

func (pe *PrivacyEnhancer) isConsentRequired(consentType string) bool {
	required := map[string]bool{
		"analytics_tracking":           false,
		"marketing_communications":     false,
		"data_sharing_partners":        false,
		"location_tracking":            true,
		"content_usage_rights":         false,
		"personalized_recommendations": false,
	}
	return required[consentType]
}
