/**
 * Analytics Handlers
 * HTTP handlers for analytics reporting APIs
 */

package handlers

import (
	"net/http"

	"lynkr/internal/services"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (ah *AnalyticsHandler) GetEngagementMetrics(c *gin.Context) {
	eventID := c.Param("id")

	metrics, err := ah.analyticsService.GetEngagementMetrics(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get engagement metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (ah *AnalyticsHandler) GetAttendanceAnalytics(c *gin.Context) {
	eventID := c.Param("id")

	analytics, err := ah.analyticsService.GetAttendanceAnalytics(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attendance analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendance": analytics})
}

func (ah *AnalyticsHandler) GetContentPerformance(c *gin.Context) {
	eventID := c.Param("id")

	performance, err := ah.analyticsService.GetContentPerformance(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get content performance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"performance": performance,
	})
}

func (ah *AnalyticsHandler) TrackEvent(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		Type    string                 `json:"type"`
		EventID string                 `json:"eventId"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := ah.analyticsService.TrackEvent(request.Type, userID, request.EventID, request.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ah *AnalyticsHandler) GetRealtimeStats(c *gin.Context) {
	eventID := c.Param("id")

	stats, err := ah.analyticsService.GetRealtimeStats(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get realtime stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
