/**
 * Pixel Tracking Handlers
 * HTTP handlers for pixel tracking and analytics
 */

package handlers

import (
	"net/http"

	"lynkr/internal/services"

	"github.com/gin-gonic/gin"
)

type PixelHandler struct {
	pixelService *services.PixelService
}

func NewPixelHandler(pixelService *services.PixelService) *PixelHandler {
	return &PixelHandler{
		pixelService: pixelService,
	}
}

func (ph *PixelHandler) TrackPixel(c *gin.Context) {
	eventID := c.Query("event")
	brandID := c.Query("brand")
	userID := c.Query("user")
	eventType := c.Query("type")

	if eventType == "" {
		eventType = "page_view"
	}

	url := c.Query("url")
	referrer := c.GetHeader("Referer")
	userAgent := c.GetHeader("User-Agent")

	// Track the pixel event
	_, err := ph.pixelService.TrackPixelEvent(userID, eventID, brandID, eventType, url, referrer, userAgent)
	if err != nil {
		// Don't return error to avoid breaking user experience
		c.Status(http.StatusOK)
		return
	}

	// Return 1x1 transparent pixel
	c.Header("Content-Type", "image/gif")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// 1x1 transparent GIF
	pixel := []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00,
		0x01, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0x21, 0xF9, 0x04, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
		0x01, 0x00, 0x3B,
	}

	c.Writer.Write(pixel)
}

func (ph *PixelHandler) TrackSearch(c *gin.Context) {
	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		EventID    string `json:"eventId"`
		BrandID    string `json:"brandId"`
		SearchTerm string `json:"searchTerm"`
		URL        string `json:"url"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := ph.pixelService.TrackPixelEvent(
		userID, request.EventID, request.BrandID, "search",
		request.URL, "", c.GetHeader("User-Agent"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track search"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tracked"})
}

func (ph *PixelHandler) GetPixelAnalytics(c *gin.Context) {
	eventID := c.Param("id")

	analytics, err := ph.pixelService.GetPixelAnalytics(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (ph *PixelHandler) GeneratePixelURL(c *gin.Context) {
	eventID := c.Query("eventId")
	brandID := c.Query("brandId")

	if eventID == "" || brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID and Brand ID required"})
		return
	}

	pixelURL := ph.pixelService.GeneratePixelURL(eventID, brandID)

	c.JSON(http.StatusOK, gin.H{
		"pixelUrl": pixelURL,
	})
}
