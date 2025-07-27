package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"lynkr/internal/middleware"
	"lynkr/internal/services/content"
	"lynkr/internal/services/event"
	"lynkr/internal/services/user"

	"github.com/gin-gonic/gin"
)

// Handler contains all the services needed for API handlers
type Handler struct {
	UserService    *user.UserService
	EventService   *event.EventService
	ContentService *content.ContentService
}

// NewHandler creates a new handler with the given services
func NewHandler(userService *user.UserService, eventService *event.EventService, contentService *content.ContentService) *Handler {
	return &Handler{
		UserService:    userService,
		EventService:   eventService,
		ContentService: contentService,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// API versioning
	v1 := router.Group("/api/v1")

	// Public routes
	v1.GET("/health", h.HealthCheck)
	v1.POST("/users", h.CreateUser)
	v1.POST("/users/login", h.Login)
	v1.GET("/events", h.ListEvents)
	v1.GET("/events/:id", h.GetEvent)
	v1.GET("/events/:id/content", h.GetEventContent)

	// Protected routes
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PUT("/users/consent", h.UpdateConsent)
		protected.POST("/events", h.CreateEvent)
		protected.POST("/events/:id/checkin", h.CheckInEvent)
		protected.POST("/events/:id/checkout", h.CheckOutEvent)
		protected.GET("/events/nearby", h.GetNearbyEvents)
		protected.GET("/events/:id/attendance/status", h.GetAttendanceStatus)
		protected.POST("/content", h.CreateContent)
		protected.PUT("/content/:id/permissions", h.UpdateContentPermissions)
		protected.GET("/content/:id/interactions", h.GetContentInteractions)
		protected.POST("/content/:id/interactions", h.AddContentInteraction)
	}

	// Admin routes
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		admin.GET("/users", h.ListUsers)
	}
}

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "connected",
	})
}

// CreateUser handles user registration
func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		//PrivacySettings map[string]interface{} `json:"privacy_settings"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	privacysetting := map[string]interface{}{}

	// Create the user
	user, err := h.UserService.Create(req.Username, req.Email, req.Password, privacysetting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles user authentication
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate the user
	user, err := h.UserService.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a JWT token
	token, err := middleware.GenerateToken(user.ID, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// UpdateConsent handles updating user consent settings
func (h *Handler) UpdateConsent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		PrivacySettings map[string]interface{} `json:"privacy_settings" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the user's privacy settings
	err := h.UserService.UpdatePrivacySettings(userID.(uint), req.PrivacySettings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update privacy settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Privacy settings updated"})
}

// ListEvents handles listing upcoming events
func (h *Handler) ListEvents(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Get the events
	events, err := h.EventService.List(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEvent handles retrieving a specific event
func (h *Handler) GetEvent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Get the event
	event, err := h.EventService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// CreateEvent handles creating a new event
func (h *Handler) CreateEvent(c *gin.Context) {
	var req struct {
		Name         string    `json:"name" binding:"required"`
		Description  string    `json:"description"`
		Location     string    `json:"location" binding:"required"`
		GeofenceData string    `json:"geofence_data"`
		StartTime    time.Time `json:"start_time" binding:"required"`
		EndTime      time.Time `json:"end_time" binding:"required"`
		BrandID      uint      `json:"brand_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the event
	event, err := h.EventService.Create(
		req.Name,
		req.Description,
		req.Location,
		req.GeofenceData,
		req.StartTime,
		req.EndTime,
		req.BrandID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// CheckInEvent handles checking in to an event
func (h *Handler) CheckInEvent(c *gin.Context) {
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

	// Create check-in request
	checkInReq := event.CheckInRequest{
		UserID:    userID.(uint),
		EventID:   uint(eventID),
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	// Check in to the event
	attendance, err := h.EventService.CheckIn(checkInReq)
	if err != nil {
		if err.Error() == "user is not within event geofence" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You must be at the event location to check in"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check in to event"})
		return
	}

	c.JSON(http.StatusCreated, attendance)
}

// GetEventContent handles retrieving content for an event
func (h *Handler) GetEventContent(c *gin.Context) {
	idStr := c.Param("id")
	// eventID, err := strconv.ParseUint(idStr, 10, 32)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
	// 	return
	// }

	// Get the content
	id, _ := strconv.Atoi(idStr)
	contents, err := h.ContentService.GetEventContent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve content"})
		return
	}

	c.JSON(http.StatusOK, contents)
}

// CreateContent handles creating new content
// func (h *Handler) CreateContent(c *gin.Context) {
// 	userID, _ := c.Get("userID")

// 	var req struct {
// 		EventID     uint                   `json:"event_id" binding:"required"`
// 		Type        content.ContentType    `json:"type" binding:"required"`
// 		URL         string                 `json:"url" binding:"required"`
// 		Permissions map[string]interface{} `json:"permissions"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Create the content
// 	content, err := h.ContentService.Create(
// 		userID.(uint),
// 		req.EventID,
// 		req.Type,
// 		req.URL,
// 		req.Permissions,
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create content"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, content)
// }

func (ch *Handler) CreateContent(c *gin.Context) {
	// userID := r.Header.Get("X-User-ID")
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	//if userID == "" {
	//	// http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//	c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	//	return
	//}
	userID := userIDVal.(uint)

	// Parse multipart form
	// err := r.ParseMultipartForm(10 << 20) // 10MB max
	err := c.Request.ParseMultipartForm(10 << 20)
	fmt.Printf("%#v\n", c.Request.MultipartForm.File)

	if err != nil {
		// http.Error(w, "Failed to parse form", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse form"})
		return
	}

	// eventID := r.FormValue("eventId")
	// caption := r.FormValue("caption")
	// tagsJSON := r.FormValue("tags")
	// permissionsJSON := r.FormValue("permissions")
	eventID := c.PostForm("eventID")
	caption := c.PostForm("caption")
	tagsJSON := c.PostForm("tags")
	permissionsJSON := c.PostForm("permissions")
	eventID1, _ := strconv.Atoi(eventID)
	// Parse tags
	var tags []content.ContentTag
	if tagsJSON != "" {
		err := json.Unmarshal([]byte(tagsJSON), &tags)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tags JSON"})
			return
		}
	}

	// Parse permissions
	var permissions content.ContentPermissions
	if permissionsJSON != "" {
		err := json.Unmarshal([]byte(permissionsJSON), &permissions)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permissions JSON"})
			return
		}
	}

	// Handle file upload (simplified - in real implementation would upload to cloud storage)
	// file, header, err := r.FormFile("media")
	file, header, err := c.Request.FormFile("media")
	if err != nil {
		// http.Error(w, "No media file provided", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"message": "No media file provided"})
		return
	}
	defer file.Close()

	// Simulate file upload and get URL
	mediaURL := "https://storage.example.com/" + header.Filename
	mediaType := "photo"
	if header.Header.Get("Content-Type") == "video/mp4" {
		mediaType = "video"
	}

	// Create content
	// content, err := ch.contentService.CreateContent(userID, eventID, mediaURL, mediaType, caption, tags, permissions)
	content, err := ch.ContentService.CreateContent(userID, eventID1, mediaURL, mediaType, caption, tags, permissions)
	if err != nil {
		// http.Error(w, "Failed to create content", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create content"})
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"content": content,
	// 	"message": "Content created successfully",
	// })
	c.JSON(http.StatusOK, gin.H{"content": content, "message": "Content created successfully"})
}

// UpdateContentPermissions handles updating content permissions
// func (h *Handler) UpdateContentPermissions(c *gin.Context) {
// 	idStr := c.Param("id")
// 	contentID, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
// 		return
// 	}

// 	var req struct {
// 		Permissions map[string]interface{} `json:"permissions" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Update the permissions
// 	err = h.ContentService.UpdateContentPermissions(contentID, req.Permissions)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permissions"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Permissions updated"})
// }

// UpdateContentPermissions updates content permissions
func (ch *Handler) UpdateContentPermissions(c *gin.Context) {
	contentID := c.Param("id")
	userID := c.Request.Header.Get("userID")

	if userID == "" {
		// http.Error(w, "Unauthorized", http.StatusUnauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var permissions content.ContentPermissions
	if err := json.NewDecoder(c.Request.Body).Decode(&permissions); err != nil {
		// http.Error(w, "Invalid request body", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	err := ch.ContentService.UpdateContentPermissions(contentID, permissions)
	if err != nil {
		// http.Error(w, "Failed to update permissions", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update permissions"})
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"message": "Permissions updated successfully",
	// })
	c.JSON(http.StatusOK, gin.H{"message": "Permissions updated successfully"})
}

// GetContentInteractions handles retrieving interactions for content
func (h *Handler) GetContentInteractions(c *gin.Context) {
	idStr := c.Param("id")
	contentID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	// Get the interactions
	interactions, err := h.ContentService.GetInteractions(uint(contentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve interactions"})
		return
	}

	c.JSON(http.StatusOK, interactions)
}

// GetEventTags retrieves popular tags for an event
func (ch *Handler) GetEventTags(c *gin.Context) {
	eventID := c.Param("id")
	tags, err := ch.ContentService.GetEventTags(eventID)
	if err != nil {
		// http.Error(w, "Failed to get event tags", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get event tags"})
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"tags": tags,
	// })
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// AddContentInteraction handles adding an interaction to content
func (h *Handler) AddContentInteraction(c *gin.Context) {
	userID, _ := c.Get("userID")

	idStr := c.Param("id")
	contentID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	var req struct {
		Type string                 `json:"type" binding:"required"`
		Data map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the interaction
	interaction, err := h.ContentService.AddInteraction(
		userID.(uint),
		uint(contentID),
		req.Type,
		req.Data,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add interaction"})
		return
	}

	c.JSON(http.StatusCreated, interaction)
}

// ListUsers handles listing all users (admin only)
func (h *Handler) ListUsers(c *gin.Context) {
	// This is a placeholder for an admin endpoint
	c.JSON(http.StatusOK, gin.H{"message": "Admin endpoint"})
}

// TrackContentAnalytics tracks content interactions
func (ch *Handler) TrackContentAnalytics(c *gin.Context) {
	contentID := c.Param("id")
	var request struct {
		Action   string                 `json:"action"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := ch.ContentService.TrackContentAnalytics(contentID, request.Action, request.Metadata)
	if err != nil {
		// http.Error(w, "Failed to track analytics", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track analytics"})
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"message": "Analytics tracked successfully",
	// })
	c.JSON(http.StatusOK, gin.H{"message": "Analytics tracked successfully"})
}

// GetContentAnalytics retrieves content analytics
func (ch *Handler) GetContentAnalytics(c *gin.Context) {
	contentID := c.Param("id")

	analytics, err := ch.ContentService.GetContentAnalytics(contentID)
	if err != nil {
		// http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"analytics": analytics,
	// })
	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
}

// ) {
// 	vars := mux.Vars(r)
// 	contentID := vars["id"]

// 	analytics, err := ch.contentService.GetContentAnalytics(contentID)
// 	if err != nil {
// 		http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"analytics": analytics,
// 	})
// }

// SearchTags searches for available tags
func (ch *Handler) SearchTags(c *gin.Context) {
	query := c.Query("query")
	eventID := c.Query("eventId")
	limitStr := c.Query("limit")

	if query == "" {
		// http.Error(w, "Query parameter required", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter required"})
		return
	}

	tags, err := ch.ContentService.SearchTags(query, eventID)
	if err != nil {
		// http.Error(w, "Failed to search tags", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tags"})
		return
	}

	// Apply limit if specified
	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit < len(tags) {
			tags = tags[:limit]
		}
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"tags": tags,
	// })
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// GetSuggestedTags returns AI-suggested tags
func (ch *Handler) GetSuggestedTags(c *gin.Context) {
	// Parse multipart form for media file
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		// http.Error(w, "Failed to parse form", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	eventID := c.Request.FormValue("id")

	// Handle file upload (simplified)
	file, header, err := c.Request.FormFile("content")
	if err != nil {
		// http.Error(w, "No content file provided", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No content file provided"})
		return
	}
	defer file.Close()

	// Simulate file processing and get suggested tags
	mediaURL := "temp://" + header.Filename
	suggestedTags, err := ch.ContentService.GetSuggestedTags(mediaURL, eventID)
	if err != nil {
		// http.Error(w, "Failed to get suggested tags", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get suggested tags"})
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"suggestedTags": suggestedTags,
	// })
	c.JSON(http.StatusOK, gin.H{"suggestedTags": suggestedTags})
}
