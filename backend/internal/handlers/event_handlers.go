package handlers

import (
	"net/http"
	"strconv"

	"lynkr/internal/services/event"

	"github.com/gin-gonic/gin"
)

// CheckOutEvent handles checking out from an event
func (h *Handler) CheckOutEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	idStr := c.Param("id")
	eventID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Parse request body
	var req struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check out from the event
	err = h.EventService.CheckOut(userID.(uint), uint(eventID), req.Latitude, req.Longitude)
	if err != nil {
		if err.Error() == "no active check-in found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No active check-in found for this event"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check out from event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully checked out from event"})
}

// GetNearbyEvents handles retrieving events near a location
func (h *Handler) GetNearbyEvents(c *gin.Context) {
	// Parse query parameters
	latStr := c.Query("latitude")
	lngStr := c.Query("longitude")
	radiusStr := c.DefaultQuery("radiusKm", "5")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Convert parameters to appropriate types
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		radius = 5.0 // Default to 5km
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 // Default to 10 events
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0 // Default to offset 0
	}

	// Get nearby events
	events, err := h.EventService.GetNearbyEvents(lat, lng, radius, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve nearby events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetAttendanceStatus checks if a user is checked in to an event
func (h *Handler) GetAttendanceStatus(c *gin.Context) {
	userID, _ := c.Get("userID")

	idStr := c.Param("id")
	eventID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Query for active attendance
	query := `
		SELECT id, user_id, event_id, check_in_time, check_out_time, created_at, latitude, longitude
		FROM attendances
		WHERE user_id = ? AND event_id = ? AND check_out_time IS NULL
		LIMIT 1
	`

	var attendance event.Attendance
	err = h.EventService.DB.QueryRow(
		query,
		userID,
		eventID,
	).Scan(
		&attendance.ID,
		&attendance.UserID,
		&attendance.EventID,
		&attendance.CheckInTime,
		&attendance.CheckOutTime,
		&attendance.CreatedAt,
		&attendance.Latitude,
		&attendance.Longitude,
	)

	if err != nil {
		// If no rows, user is not checked in
		c.JSON(http.StatusNotFound, gin.H{"error": "User is not checked in to this event"})
		return
	}

	c.JSON(http.StatusOK, attendance)
}
