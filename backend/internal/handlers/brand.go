/**
 * Brand Handlers
 * HTTP handlers for brand authentication and management
 */

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	// "golang.org/x/crypto/bcrypt"
	"lynkr/internal/services"
)

type BrandHandler struct {
	brandService *services.BrandService
	jwtSecret    []byte
}

type BrandLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BrandLoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type BrandClaims struct {
	BrandID string `json:"brandId"`
	UserID  string `json:"userId"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

func NewBrandHandler(brandService *services.BrandService, jwtSecret string) *BrandHandler {
	return &BrandHandler{
		brandService: brandService,
		jwtSecret:    []byte(jwtSecret),
	}
}

// Login handles brand user authentication
func (bh *BrandHandler) Login(c *gin.Context) {
	var req BrandLoginRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate credentials (simplified - in real implementation would check database)
	if req.Email == "brand@example.com" && req.Password == "password" {
		// Generate JWT token
		claims := BrandClaims{
			BrandID: "brand_1",
			UserID:  "user_1",
			Role:    "brand",
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        "brand_1",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(bh.jwtSecret)
		if err != nil {
			// http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		user := map[string]interface{}{
			"id":      "user_1",
			"name":    "Brand Manager",
			"email":   req.Email,
			"brandId": 123,
			"role":    "brand",
		}

		response := BrandLoginResponse{
			Token: tokenString,
			User:  user,
		}

		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(response)
		c.JSON(http.StatusOK, response)
	} else {
		// http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// GetDashboardStats returns dashboard statistics
func (bh *BrandHandler) GetDashboardStats(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		// http.Error(w, "Brand ID required", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}
	// Simulate dashboard data
	stats := map[string]interface{}{
		"totalAttendees": 1247,
		"contentPieces":  89,
		"engagementRate": 23.5,
		"conversionRate": 4.2,
		"attendanceData": []map[string]interface{}{
			{"date": "2024-01", "attendees": 120},
			{"date": "2024-02", "attendees": 180},
			{"date": "2024-03", "attendees": 250},
			{"date": "2024-04", "attendees": 320},
			{"date": "2024-05", "attendees": 280},
			{"date": "2024-06", "attendees": 380},
		},
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(stats)
	//fmt.Printf(brandID.[string])
	c.JSON(http.StatusOK, stats)
}

// GetCampaigns returns brand campaigns
func (bh *BrandHandler) GetCampaigns(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	campaigns := []map[string]interface{}{
		{
			"id":        "1",
			"name":      "Tech Conference 2024",
			"status":    "active",
			"startDate": "2024-03-15",
			"endDate":   "2024-03-17",
			"budget":    50000,
			"attendees": 1247,
		},
		{
			"id":        "2",
			"name":      "Product Launch Event",
			"status":    "completed",
			"startDate": "2024-02-10",
			"endDate":   "2024-02-12",
			"budget":    75000,
			"attendees": 890,
		},
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"campaigns": campaigns,
	// })
	c.JSON(http.StatusOK, map[string]interface{}{
		"campaigns": campaigns,
	})
}

// CreateCampaign creates a new campaign
func (bh *BrandHandler) CreateCampaign(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		// http.Error(w, "Brand ID required", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	var campaign map[string]interface{}
	// if err := json.NewDecoder(r.Body).Decode(&campaign); err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	if err := c.ShouldBindJSON(&campaign); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Add campaign ID and default values
	campaign["id"] = time.Now().UnixNano()
	campaign["status"] = "draft"
	campaign["attendees"] = 0
	campaign["brandId"] = brandID

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"campaign": campaign,
	// 	"message":  "Campaign created successfully",
	// })
	c.JSON(http.StatusOK, map[string]interface{}{
		"campaign": campaign,
		"message":  "Campaign created successfully",
	})
}

// GetBrandContent returns content accessible to the brand
func (bh *BrandHandler) GetBrandContent(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		// http.Error(w, "Brand ID required", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	content := []map[string]interface{}{
		{
			"id":        "1",
			"mediaUrl":  "https://via.placeholder.com/300x300",
			"mediaType": "photo",
			"caption":   "Amazing product demo at the tech conference!",
			"tags":      []string{"tech", "product-demo", "innovation"},
			"createdAt": "2024-03-15T10:30:00Z",
			"eventName": "Tech Conference 2024",
			"engagement": map[string]int{
				"views":  1250,
				"shares": 45,
				"likes":  189,
			},
		},
		{
			"id":        "2",
			"mediaUrl":  "https://via.placeholder.com/300x300",
			"mediaType": "photo",
			"caption":   "Great networking session with industry leaders",
			"tags":      []string{"networking", "conference", "business"},
			"createdAt": "2024-03-15T14:20:00Z",
			"eventName": "Tech Conference 2024",
			"engagement": map[string]int{
				"views":  890,
				"shares": 23,
				"likes":  156,
			},
		},
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"content": content,
	// })
	c.JSON(http.StatusOK, map[string]interface{}{
		"content": content,
	})

}
