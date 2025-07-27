/**
 * E-commerce Handlers
 * HTTP handlers for e-commerce integration and purchase tracking
 */

package handlers

import (
	"net/http"

	"lynkr/internal/services"

	"github.com/gin-gonic/gin"
)

type EcommerceHandler struct {
	ecommerceService *services.EcommerceService
}

func NewEcommerceHandler(ecommerceService *services.EcommerceService) *EcommerceHandler {
	return &EcommerceHandler{
		ecommerceService: ecommerceService,
	}
}

func (eh *EcommerceHandler) CreateIntegration(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	var request struct {
		PlatformType string `json:"platformType"`
		APIKey       string `json:"apiKey"`
		StoreURL     string `json:"storeUrl"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	integration, err := eh.ecommerceService.CreateIntegration(
		brandID, request.PlatformType, request.APIKey, request.StoreURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create integration"})
		return
	}

	c.JSON(http.StatusOK, integration)
}

func (eh *EcommerceHandler) GetIntegration(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	integration, err := eh.ecommerceService.GetIntegration(brandID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Integration not found"})
		return
	}

	c.JSON(http.StatusOK, integration)
}

func (eh *EcommerceHandler) TrackPurchase(c *gin.Context) {
	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		ProductID string  `json:"productId"`
		EventID   string  `json:"eventId"`
		Amount    float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	purchase, err := eh.ecommerceService.TrackPurchase(
		userID, request.ProductID, request.EventID, request.Amount,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track purchase"})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (eh *EcommerceHandler) GetPurchaseAnalytics(c *gin.Context) {
	eventID := c.Param("id")

	analytics, err := eh.ecommerceService.GetPurchaseAnalytics(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get purchase analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (eh *EcommerceHandler) GetTopProducts(c *gin.Context) {
	eventID := c.GetString("brandID")

	products, err := eh.ecommerceService.GetTopProducts(eventID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (eh *EcommerceHandler) HandleWebhook(c *gin.Context) {
	integrationID := c.Param("integrationId")
	if integrationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Integration ID required"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook payload"})
		return
	}

	eventType, ok := payload["event"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type"})
		return
	}

	switch eventType {
	case "order.created", "order.completed":
		c.JSON(http.StatusOK, gin.H{"status": "processed"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown event type"})
	}
}
