package event

import (
	"database/sql"
	"errors"
	"lynkr/pkg/geofencing"
	"math"
	"time"
)

// Event represents an event in the system
type Event struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	GeofenceData string    `json:"geofence_data"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	BrandID      string    `json:"brand_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Attendance represents a user's attendance at an event
type Attendance struct {
	ID           uint       `json:"id"`
	UserID       uint       `json:"user_id"`
	EventID      uint       `json:"event_id"`
	CheckInTime  time.Time  `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	CreatedAt    time.Time  `json:"created_at"`
	Latitude     float64    `json:"latitude,omitempty"`
	Longitude    float64    `json:"longitude,omitempty"`
}

// CheckInRequest represents a request to check in to an event
type CheckInRequest struct {
	UserID    uint    `json:"user_id"`
	EventID   uint    `json:"event_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// EventService handles event-related operations
type EventService struct {
	DB *sql.DB
}

// NewEventService creates a new event service
func NewEventService(db *sql.DB) *EventService {
	return &EventService{DB: db}
}

// Create creates a new event
func (s *EventService) Create(name, description, location, geofenceData string, startTime, endTime time.Time, brandID uint) (*Event, error) {
	// Validate geofence data if provided
	if geofenceData != "" {
		_, err := geofencing.ParseGeofenceData(geofenceData)
		if err != nil {
			return nil, err
		}
	}

	query := `
		INSERT INTO events (name, description, location, geofence_data, start_time, end_time, brand_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id, name, description, location, geofence_data, start_time, end_time, brand_id, created_at, updated_at
	`

	var event Event
	err := s.DB.QueryRow(
		query,
		name,
		description,
		location,
		geofenceData,
		startTime,
		endTime,
		brandID,
	).Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.GeofenceData,
		&event.StartTime,
		&event.EndTime,
		&event.BrandID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

// GetByID retrieves an event by ID
func (s *EventService) GetByID(id uint) (*Event, error) {
	query := `
		SELECT id, name, description, location, geofence_data, start_time, end_time, brand_id, created_at, updated_at
		FROM events
		WHERE id = ?
	`

	var event Event
	err := s.DB.QueryRow(query, id).Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.GeofenceData,
		&event.StartTime,
		&event.EndTime,
		&event.BrandID,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("event not found")
		}
		return nil, err
	}

	return &event, nil
}

// List retrieves a list of upcoming events
func (s *EventService) List(limit, offset int) ([]Event, error) {
	query := `
		SELECT id, name, description, location, geofence_data, start_time, end_time, brand_id, created_at, updated_at
		FROM events
		ORDER BY start_time ASC
		LIMIT ? OFFSET ?
	`

	rows, err := s.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Location,
			&event.GeofenceData,
			&event.StartTime,
			&event.EndTime,
			&event.BrandID,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// CheckIn records a user's attendance at an event
func (s *EventService) CheckIn(req CheckInRequest) (*Attendance, error) {
	// Check if the event exists
	event, err := s.GetByID(req.EventID)
	if err != nil {
		return nil, err
	}

	// Check if event has started and not ended
	now := time.Now()
	if now.Before(event.StartTime) {
		return nil, errors.New("event has not started yet")
	}
	if now.After(event.EndTime) {
		return nil, errors.New("event has already ended")
	}

	// Verify user is within geofence if geofence data is available
	if event.GeofenceData != "" {
		geofenceData, err := geofencing.ParseGeofenceData(event.GeofenceData)
		if err != nil {
			return nil, err
		}

		point := geofencing.Point{
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		}

		if !geofencing.IsPointInGeofence(point, geofenceData) {
			return nil, errors.New("user is not within event geofence")
		}
	}

	// Record the check-in
	query := `
		INSERT INTO attendances (user_id, event_id, check_in_time, latitude, longitude)
		VALUES (?, ?, CURRENT_TIMESTAMP, ?, ?)
		RETURNING id, user_id, event_id, check_in_time, check_out_time, created_at
	`

	var attendance Attendance
	err = s.DB.QueryRow(
		query,
		req.UserID,
		req.EventID,
		req.Latitude,
		req.Longitude,
	).Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.EventID,
		&attendance.CheckInTime,
		&attendance.CheckOutTime,
		&attendance.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Add location data to response
	attendance.Latitude = req.Latitude
	attendance.Longitude = req.Longitude

	return &attendance, nil
}

// CheckOut records a user's departure from an event
func (s *EventService) CheckOut(userID, eventID uint, latitude, longitude float64) error {
	query := `
		UPDATE attendances
		SET check_out_time = CURRENT_TIMESTAMP, 
		    latitude = COALESCE(latitude, ?), 
		    longitude = COALESCE(longitude, ?)
		WHERE user_id = ? AND event_id = ? AND check_out_time IS NULL
	`

	result, err := s.DB.Exec(query, latitude, longitude, userID, eventID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no active check-in found")
	}

	return nil
}

// GetAttendees retrieves a list of users who attended an event
func (s *EventService) GetAttendees(eventID uint) ([]uint, error) {
	query := `
		SELECT user_id
		FROM attendances
		WHERE event_id = ?
	`

	rows, err := s.DB.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uint
	for rows.Next() {
		var userID uint
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

// GetNearbyEvents retrieves events near a location
func (s *EventService) GetNearbyEvents(latitude, longitude float64, radiusKm float64, limit, offset int) ([]Event, error) {
	// This is a simplified implementation that doesn't use spatial indexing
	// In a production environment, you would use a spatial database or service

	// Get all upcoming events
	events, err := s.List(100, 0) // Get more events than needed for filtering
	if err != nil {
		return nil, err
	}

	// Filter events by distance
	var nearbyEvents []Event
	point := geofencing.Point{Latitude: latitude, Longitude: longitude}

	for _, event := range events {
		// If event has geofence data, check if point is within or near the geofence
		if event.GeofenceData != "" {
			geofenceData, err := geofencing.ParseGeofenceData(event.GeofenceData)
			if err != nil {
				continue
			}

			// For circle geofences, check if point is within radius + radiusKm
			if geofenceData.Type == geofencing.CircleGeofence {
				center := geofenceData.Circle.Center
				distance := calculateDistance(point.Latitude, point.Longitude, center.Latitude, center.Longitude)

				// Convert radiusKm to meters
				if distance <= (geofenceData.Circle.Radius + (radiusKm * 1000)) {
					nearbyEvents = append(nearbyEvents, event)
				}
			}
		}
	}

	// Apply limit and offset
	if offset >= len(nearbyEvents) {
		return []Event{}, nil
	}

	end := offset + limit
	if end > len(nearbyEvents) {
		end = len(nearbyEvents)
	}

	return nearbyEvents[offset:end], nil
}

// calculateDistance calculates the distance between two points using the Haversine formula
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // Earth radius in meters

	// Convert degrees to radians
	// lat1Rad := lat1 * (3.14159265359 / 180)
	// lat2Rad := lat2 * (3.14159265359 / 180)
	// lonDiff := (lon2 - lon1) * (3.14159265359 / 180)
	// latDiff := (lat2 - lat1) * (3.14159265359 / 180)
	lat1Rad := lat1 * (math.Pi / 180)
	lat2Rad := lat2 * (math.Pi / 180)
	deltaLat := (lat2 - lat1) * (math.Pi / 180)
	deltaLon := (lon2 - lon1) * (math.Pi / 180)

	// a := (1-2*latDiff)*(1-2*latDiff) +
	// 	2*lat1Rad*2*lat2Rad*(1-2*lonDiff)*(1-2*lonDiff)
	// c := 2 * 2 * latDiff * 2 * latDiff
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
