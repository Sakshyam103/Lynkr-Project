/**
 * Advanced Analytics Handlers
 * HTTP handlers for AI tagging and conversion funnel analytics
 */

package handlers

import (
	"net/http"

	"lynkr/internal/services"

	"github.com/gin-gonic/gin"
)

type AdvancedAnalyticsHandler struct {
	aiTaggingService        *services.AITaggingService
	conversionFunnelService *services.ConversionFunnelService
}

func NewAdvancedAnalyticsHandler(aiTaggingService *services.AITaggingService, conversionFunnelService *services.ConversionFunnelService) *AdvancedAnalyticsHandler {
	return &AdvancedAnalyticsHandler{
		aiTaggingService:        aiTaggingService,
		conversionFunnelService: conversionFunnelService,
	}
}

func (aah *AdvancedAnalyticsHandler) ProcessContentAI(c *gin.Context) {
	contentID := c.Param("id")

	var request struct {
		MediaURL string `json:"mediaUrl"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := aah.aiTaggingService.ProcessContent(contentID, request.MediaURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process content"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (aah *AdvancedAnalyticsHandler) GetProductAnalytics(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	analytics, err := aah.aiTaggingService.GetProductAnalytics(brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (aah *AdvancedAnalyticsHandler) GetConversionFunnel(c *gin.Context) {
	eventID := c.Param("id")
	brandID := c.GetString("brandID")

	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	funnel, err := aah.conversionFunnelService.GetConversionFunnel(eventID, brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversion funnel"})
		return
	}

	c.JSON(http.StatusOK, funnel)
}

func (aah *AdvancedAnalyticsHandler) GetAttributionReport(c *gin.Context) {
	eventID := c.Param("id")

	report, err := aah.conversionFunnelService.GetAttributionReport(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get attribution report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (aah *AdvancedAnalyticsHandler) TrackConversion(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		EventID  string                 `json:"eventId"`
		Stage    string                 `json:"stage"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := aah.conversionFunnelService.TrackConversion(userID, request.EventID, request.Stage, request.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track conversion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tracked"})
}
