package services

import (
	"database/sql"
	"fmt"
	"time"
)

type AttendanceService struct {
	db *sql.DB
}

type Attendance struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	EventID     string    `json:"eventId"`
	CheckinTime *time.Time `json:"checkinTime"`
	CheckoutTime *time.Time `json:"checkoutTime"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewAttendanceService(db *sql.DB) *AttendanceService {
	return &AttendanceService{db: db}
}

func (as *AttendanceService) CheckIn(userID, eventID string) error {
	query := `
		INSERT INTO attendances (user_id, event_id, checkin_time, created_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, event_id) DO UPDATE SET
		checkin_time = excluded.checkin_time
	`
	
	_, err := as.db.Exec(query, userID, eventID, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to check in: %w", err)
	}
	
	return nil
}

func (as *AttendanceService) CheckOut(userID, eventID string) error {
	query := `
		UPDATE attendances 
		SET checkout_time = ?
		WHERE user_id = ? AND event_id = ? AND checkout_time IS NULL
	`
	
	_, err := as.db.Exec(query, time.Now(), userID, eventID)
	if err != nil {
		return fmt.Errorf("failed to check out: %w", err)
	}
	
	return nil
}

func (as *AttendanceService) GetEventAttendances(eventID string) ([]Attendance, error) {
	query := `
		SELECT id, user_id, event_id, checkin_time, checkout_time, created_at
		FROM attendances
		WHERE event_id = ?
		ORDER BY checkin_time DESC
	`
	
	rows, err := as.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendances: %w", err)
	}
	defer rows.Close()
	
	var attendances []Attendance
	for rows.Next() {
		var a Attendance
		err := rows.Scan(&a.ID, &a.UserID, &a.EventID, &a.CheckinTime, &a.CheckoutTime, &a.CreatedAt)
		if err != nil {
			continue
		}
		attendances = append(attendances, a)
	}
	
	return attendances, nil
}

func (as *AttendanceService) IsUserCheckedIn(userID, eventID string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM attendances
		WHERE user_id = ? AND event_id = ? AND checkin_time IS NOT NULL AND checkout_time IS NULL
	`
	
	var count int
	err := as.db.QueryRow(query, userID, eventID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check attendance status: %w", err)
	}
	
	return count > 0, nil
}