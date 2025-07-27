package handlers

//
//import (
//	"net/http"
//	"strconv"
//
//	"github.com/gin-gonic/gin"
//	"lynkr/internal/services"
//)
//
//type AnalyticsHandler struct {
//	analyticsService   *services.AnalyticsService
//	attendanceService  *services.AttendanceService
//	sentimentService   *services.SentimentService
//}
//
//func NewAnalyticsHandler(analyticsService *services.AnalyticsService, attendanceService *services.AttendanceService, sentimentService *services.SentimentService) *AnalyticsHandler {
//	return &AnalyticsHandler{
//		analyticsService:  analyticsService,
//		attendanceService: attendanceService,
//		sentimentService:  sentimentService,
//	}
//}
//
//// @Summary Get engagement metrics for an event
//// @Description Get engagement metrics including views, likes, shares, comments
//// @Tags analytics
//// @Accept json
//// @Produce json
//// @Param eventId path string true "Event ID"
//// @Success 200 {object} services.EngagementMetrics
//// @Router /events/{eventId}/analytics/engagement [get]
//func (ah *AnalyticsHandler) GetEngagementMetrics(c *gin.Context) {
//	eventID := c.Param("id")
//	if eventID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
//		return
//	}
//
//	metrics, err := ah.analyticsService.GetEngagementMetrics(eventID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get engagement metrics"})
//		return
//	}
//
//	c.JSON(http.StatusOK, metrics)
//}
//
//// @Summary Get attendance analytics for an event
//// @Description Get attendance analytics including unique attendees, total check-ins
//// @Tags analytics
//// @Accept json
//// @Produce json
//// @Param eventId path string true "Event ID"
//// @Success 200 {object} map[string]interface{}
//// @Router /events/{eventId}/analytics/attendance [get]
//func (ah *AnalyticsHandler) GetAttendanceAnalytics(c *gin.Context) {
//	eventID := c.Param("id")
//	if eventID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
//		return
//	}
//
//	analytics, err := ah.analyticsService.GetAttendanceAnalytics(eventID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attendance analytics"})
//		return
//	}
//
//	c.JSON(http.StatusOK, analytics)
//}
//
//// @Summary Get content performance for an event
//// @Description Get content performance metrics for all content in an event
//// @Tags analytics
//// @Accept json
//// @Produce json
//// @Param eventId path string true "Event ID"
//// @Success 200 {array} services.ContentPerformance
//// @Router /events/{eventId}/analytics/content [get]
//func (ah *AnalyticsHandler) GetContentPerformance(c *gin.Context) {
//	eventID := c.Param("id")
//	if eventID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
//		return
//	}
//
//	performance, err := ah.analyticsService.GetContentPerformance(eventID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get content performance"})
//		return
//	}
//
//	c.JSON(http.StatusOK, performance)
//}
//
//// @Summary Get real-time stats for an event
//// @Description Get real-time statistics including active users and recent content
//// @Tags analytics
//// @Accept json
//// @Produce json
//// @Param eventId path string true "Event ID"
//// @Success 200 {object} map[string]interface{}
//// @Router /events/{eventId}/analytics/realtime [get]
//func (ah *AnalyticsHandler) GetRealtimeStats(c *gin.Context) {
//	eventID := c.Param("id")
//	if eventID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
//		return
//	}
//
//	stats, err := ah.analyticsService.GetRealtimeStats(eventID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get real-time stats"})
//		return
//	}
//
//	c.JSON(http.StatusOK, stats)
//}
//
//// @Summary Track an analytics event
//// @Description Track a custom analytics event
//// @Tags analytics
//// @Accept json
//// @Produce json
//// @Param request body map[string]interface{} true "Event data"
//// @Success 200 {object} map[string]string
//// @Router /analytics/track [post]
//func (ah *AnalyticsHandler) TrackEvent(c *gin.Context) {
//	var request struct {
//		Type    string                 `json:"type" binding:"required"`
//		UserID  string                 `json:"userId"`
//		EventID string                 `json:"eventId"`
//		Data    map[string]interface{} `json:"data"`
//	}
//
//	if err := c.ShouldBindJSON(&request); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
//		return
//	}
//
//	err := ah.analyticsService.TrackEvent(request.Type, request.UserID, request.EventID, request.Data)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track event"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Event tracked successfully"})
//}
//
//// @Summary Get event sentiment analysis
//// @Description Get aggregated sentiment analysis for an event
//// @Tags analytics
//// @Accept json
//// @Produce json
//// @Param eventId path string true "Event ID"
//// @Success 200 {object} map[string]interface{}
//// @Router /events/{eventId}/sentiment [get]
//func (ah *AnalyticsHandler) GetEventSentiment(c *gin.Context) {
//	eventID := c.Param("id")
//	if eventID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
//		return
//	}
//
//	sentiment, err := ah.sentimentService.GetEventSentiment(eventID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get event sentiment"})
//		return
//	}
//
//	c.JSON(http.StatusOK, sentiment)
//}
//
//// @Summary Analyze text sentiment
//// @Description Analyze sentiment of provided text
//// @Tags analytics
//// @Accept json
//// @Produce json
//// @Param request body map[string]string true "Text to analyze"
//// @Success 200 {object} services.SentimentResult
//// @Router /sentiment/analyze [post]
//func (ah *AnalyticsHandler) AnalyzeSentiment(c *gin.Context) {
//	var request struct {
//		Text string `json:"text" binding:"required"`
//	}
//
//	if err := c.ShouldBindJSON(&request); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Text is required"})
//		return
//	}
//
//	result, err := ah.sentimentService.AnalyzeSentiment(request.Text)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze sentiment"})
//		return
//	}
//
//	c.JSON(http.StatusOK, result)
//}
//
//// @Summary Check in to event
//// @Description Check user into an event
//// @Tags attendance
//// @Accept json
//// @Produce json
//// @Param eventId path string true "Event ID"
//// @Param request body map[string]interface{} true "Check-in data"
//// @Success 200 {object} map[string]string
//// @Router /events/{eventId}/checkin [post]
//func (ah *AnalyticsHandler) CheckInEvent(c *gin.Context) {
//	eventID := c.Param("id")
//	userID := c.GetString("user_id")
//
//	if eventID == "" || userID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID and user ID are required"})
//		return
//	}
//
//	err := ah.attendanceService.CheckIn(userID, eventID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check in"})
//		return
//	}
//
//	// Track analytics event
//	ah.analyticsService.TrackEvent("checkin", userID, eventID, map[string]interface{}{
//		"action": "checkin",
//		"value": 1,
//	})
//
//	c.JSON(http.StatusOK, gin.H{"message": "Checked in successfully"})
//}
//
//// @Summary Get event attendances
//// @Description Get all attendances for an event
//// @Tags attendance
//// @Accept json
//// @Produce json
//// @Param eventId path string true "Event ID"
//// @Success 200 {array} services.Attendance
//// @Router /events/{eventId}/attendances [get]
//func (ah *AnalyticsHandler) GetEventAttendances(c *gin.Context) {
//	eventID := c.Param("id")
//	if eventID == "" {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
//		return
//	}
//
//	attendances, err := ah.attendanceService.GetEventAttendances(eventID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attendances"})
//		return
//	}
//
//	c.JSON(http.StatusOK, attendances)
//}
