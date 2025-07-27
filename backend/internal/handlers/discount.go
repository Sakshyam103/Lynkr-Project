/**
 * Discount Code Handlers
 * HTTP handlers for discount code management and redemption
 */

package handlers

import (
	"net/http"
	"time"

	"lynkr/internal/services"

	"github.com/gin-gonic/gin"
)

type DiscountHandler struct {
	discountService *services.DiscountService
}

func NewDiscountHandler(discountService *services.DiscountService) *DiscountHandler {
	return &DiscountHandler{
		discountService: discountService,
	}
}

func (dh *DiscountHandler) GenerateCode(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	var request struct {
		EventID     string  `json:"eventId"`
		DiscountPct float64 `json:"discountPct"`
		MaxUses     int     `json:"maxUses"`
		ExpiresIn   int     `json:"expiresIn"` // days
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	expiresAt := time.Now().AddDate(0, 0, request.ExpiresIn)

	code, err := dh.discountService.GenerateCode(
		request.EventID, brandID, request.DiscountPct, request.MaxUses, expiresAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate code"})
		return
	}

	c.JSON(http.StatusOK, code)
}

func (dh *DiscountHandler) ValidateCode(c *gin.Context) {
	code := c.Param("code")

	discountCode, err := dh.discountService.ValidateCode(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount code"})
		return
	}

	c.JSON(http.StatusOK, discountCode)
}

func (dh *DiscountHandler) RedeemCode(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		Code    string  `json:"code"`
		OrderID string  `json:"orderId"`
		Amount  float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate code first
	discountCode, err := dh.discountService.ValidateCode(request.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount code"})
		return
	}

	// Redeem code
	redemption, err := dh.discountService.RedeemCode(
		discountCode.ID, userID, request.OrderID, request.Amount,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to redeem code"})
		return
	}

	c.JSON(http.StatusOK, redemption)
}

func (dh *DiscountHandler) GetCodeAnalytics(c *gin.Context) {
	eventID := c.Param("id")

	analytics, err := dh.discountService.GetCodeAnalytics(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (dh *DiscountHandler) GetBrandCodes(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	// Simulate getting brand codes
	codes := []map[string]interface{}{
		{
			"id":          "discount_1",
			"code":        "EVENT20",
			"discountPct": 20.0,
			"usedCount":   45,
			"maxUses":     100,
			"expiresAt":   time.Now().AddDate(0, 0, 30),
		},
		{
			"id":          "discount_2",
			"code":        "SPECIAL15",
			"discountPct": 15.0,
			"usedCount":   23,
			"maxUses":     50,
			"expiresAt":   time.Now().AddDate(0, 0, 15),
		},
	}

	c.JSON(http.StatusOK, gin.H{"codes": codes})
}
