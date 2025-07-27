package privacy

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

// RetentionManager handles data retention policies
type RetentionManager struct {
	db        *sql.DB
	anonymizer *Anonymizer
}

// NewRetentionManager creates a new retention manager
func NewRetentionManager(db *sql.DB, anonymizer *Anonymizer) *RetentionManager {
	return &RetentionManager{
		db:        db,
		anonymizer: anonymizer,
	}
}

// ApplyRetentionPolicies applies retention policies to all user data
func (r *RetentionManager) ApplyRetentionPolicies() error {
	// Get all users with their privacy settings
	rows, err := r.db.Query(`
		SELECT id, privacy_settings, created_at
		FROM users
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			userID         uint
			privacySettingsJSON string
			createdAt      time.Time
		)

		err := rows.Scan(&userID, &privacySettingsJSON, &createdAt)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}

		// Parse privacy settings
		var privacySettings map[string]interface{}
		err = json.Unmarshal([]byte(privacySettingsJSON), &privacySettings)
		if err != nil {
			log.Printf("Error parsing privacy settings for user %d: %v", userID, err)
			continue
		}

		// Apply retention policy to user data
		err = r.applyUserRetentionPolicy(userID, privacySettings)
		if err != nil {
			log.Printf("Error applying retention policy for user %d: %v", userID, err)
		}
	}

	return nil
}

// applyUserRetentionPolicy applies retention policy to a specific user's data
func (r *RetentionManager) applyUserRetentionPolicy(userID uint, privacySettings map[string]interface{}) error {
	// Apply retention policy to content
	err := r.applyContentRetentionPolicy(userID, privacySettings)
	if err != nil {
		return err
	}

	// Apply retention policy to interactions
	err = r.applyInteractionsRetentionPolicy(userID, privacySettings)
	if err != nil {
		return err
	}

	// Apply retention policy to attendances
	err = r.applyAttendancesRetentionPolicy(userID, privacySettings)
	if err != nil {
		return err
	}

	return nil
}

// applyContentRetentionPolicy applies retention policy to user content
func (r *RetentionManager) applyContentRetentionPolicy(userID uint, privacySettings map[string]interface{}) error {
	// Get all content for the user
	rows, err := r.db.Query(`
		SELECT id, created_at
		FROM content
		WHERE user_id = ?
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			contentID  uint
			createdAt  time.Time
		)

		err := rows.Scan(&contentID, &createdAt)
		if err != nil {
			log.Printf("Error scanning content: %v", err)
			continue
		}

		// Check if content should be retained
		if !r.anonymizer.ShouldRetainData(createdAt, privacySettings) {
			// Delete content that exceeds retention period
			_, err := r.db.Exec(`DELETE FROM content WHERE id = ?`, contentID)
			if err != nil {
				log.Printf("Error deleting content %d: %v", contentID, err)
			} else {
				log.Printf("Deleted content %d due to retention policy", contentID)
			}
		}
	}

	return nil
}

// applyInteractionsRetentionPolicy applies retention policy to user interactions
func (r *RetentionManager) applyInteractionsRetentionPolicy(userID uint, privacySettings map[string]interface{}) error {
	// Get all interactions for the user
	rows, err := r.db.Query(`
		SELECT id, created_at
		FROM interactions
		WHERE user_id = ?
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			interactionID uint
			createdAt     time.Time
		)

		err := rows.Scan(&interactionID, &createdAt)
		if err != nil {
			log.Printf("Error scanning interaction: %v", err)
			continue
		}

		// Check if interaction should be retained
		if !r.anonymizer.ShouldRetainData(createdAt, privacySettings) {
			// Delete interaction that exceeds retention period
			_, err := r.db.Exec(`DELETE FROM interactions WHERE id = ?`, interactionID)
			if err != nil {
				log.Printf("Error deleting interaction %d: %v", interactionID, err)
			} else {
				log.Printf("Deleted interaction %d due to retention policy", interactionID)
			}
		}
	}

	return nil
}

// applyAttendancesRetentionPolicy applies retention policy to user attendances
func (r *RetentionManager) applyAttendancesRetentionPolicy(userID uint, privacySettings map[string]interface{}) error {
	// Get all attendances for the user
	rows, err := r.db.Query(`
		SELECT id, created_at
		FROM attendances
		WHERE user_id = ?
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			attendanceID uint
			createdAt    time.Time
		)

		err := rows.Scan(&attendanceID, &createdAt)
		if err != nil {
			log.Printf("Error scanning attendance: %v", err)
			continue
		}

		// Check if attendance should be retained
		if !r.anonymizer.ShouldRetainData(createdAt, privacySettings) {
			// Delete attendance that exceeds retention period
			_, err := r.db.Exec(`DELETE FROM attendances WHERE id = ?`, attendanceID)
			if err != nil {
				log.Printf("Error deleting attendance %d: %v", attendanceID, err)
			} else {
				log.Printf("Deleted attendance %d due to retention policy", attendanceID)
			}
		}
	}

	return nil
}

// ScheduleRetentionJob schedules the retention job to run periodically
func (r *RetentionManager) ScheduleRetentionJob(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			log.Println("Running scheduled data retention job")
			err := r.ApplyRetentionPolicies()
			if err != nil {
				log.Printf("Error running retention job: %v", err)
			}
		}
	}()
	log.Printf("Scheduled data retention job to run every %v", interval)
}