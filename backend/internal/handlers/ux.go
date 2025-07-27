/**
 * UX Handlers
 * HTTP handlers for user experience monitoring and usability testing
 */

package handlers

import (
	"lynkr/internal/ux"

	"github.com/gin-gonic/gin"
)

type UXHandler struct {
	usabilityTester *ux.UsabilityTester
}

func NewUXHandler(usabilityTester *ux.UsabilityTester) *UXHandler {
	return &UXHandler{
		usabilityTester: usabilityTester,
	}
}

func (uh *UXHandler) StartUsabilitySession(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	session, err := uh.usabilityTester.StartSession(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to start session"})
		return
	}

	c.JSON(200, session)
}

func (uh *UXHandler) TrackUserAction(c *gin.Context) {
	var request struct {
		SessionID string `json:"sessionId"`
		Type      string `json:"type"`
		Screen    string `json:"screen"`
		Element   string `json:"element"`
		Duration  int    `json:"duration"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := uh.usabilityTester.TrackAction(request.SessionID, request.Type, request.Screen, request.Element, request.Duration); err != nil {
		c.JSON(500, gin.H{"error": "Failed to track action"})
		return
	}

	c.JSON(200, gin.H{"status": "tracked"})
}

func (uh *UXHandler) TrackUserError(c *gin.Context) {
	var request struct {
		SessionID   string `json:"sessionId"`
		Type        string `json:"type"`
		Message     string `json:"message"`
		Screen      string `json:"screen"`
		Recoverable bool   `json:"recoverable"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := uh.usabilityTester.TrackError(request.SessionID, request.Type, request.Message, request.Screen, request.Recoverable); err != nil {
		c.JSON(500, gin.H{"error": "Failed to track error"})
		return
	}

	c.JSON(200, gin.H{"status": "tracked"})
}

func (uh *UXHandler) EndUsabilitySession(c *gin.Context) {
	var request struct {
		SessionID string `json:"sessionId"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if err := uh.usabilityTester.EndSession(request.SessionID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to end session"})
		return
	}

	c.JSON(200, gin.H{"status": "session_ended"})
}

func (uh *UXHandler) GetUsabilityMetrics(c *gin.Context) {
	timeframe := c.Query("timeframe")
	if timeframe == "" {
		timeframe = "7d"
	}

	metrics, err := uh.usabilityTester.GetUsabilityMetrics(timeframe)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get usability metrics"})
		return
	}

	c.JSON(200, metrics)
}

func (uh *UXHandler) GetHeatmapData(c *gin.Context) {
	screen := c.Query("screen")
	if screen == "" {
		c.JSON(400, gin.H{"error": "Screen parameter required"})
		return
	}

	heatmapData, err := uh.usabilityTester.GetHeatmapData(screen)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get heatmap data"})
		return
	}

	c.JSON(200, gin.H{
		"screen": screen,
		"data":   heatmapData,
	})
}

func (uh *UXHandler) GetUserJourney(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
	}

	if userID == "" {
		c.JSON(400, gin.H{"error": "User ID required"})
		return
	}

	journey, err := uh.usabilityTester.GetUserJourney(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user journey"})
		return
	}

	c.JSON(200, gin.H{
		"userId":  userID,
		"journey": journey,
	})
}

func (uh *UXHandler) GetPainPoints(c *gin.Context) {
	painPoints, err := uh.usabilityTester.IdentifyPainPoints()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get pain points"})
		return
	}

	c.JSON(200, gin.H{
		"painPoints": painPoints,
	})
}
