/**
 * Content Handlers
 * HTTP handlers for content creation, management, and sharing
 */

package handlers

// import (
// 	"encoding/json"
// 	"net/http"
// 	"strconv"

// 	"lynkr/internal/services"

// 	"github.com/gorilla/mux"
// )

// type ContentHandler struct {
// 	contentService *services.ContentService
// }

// func NewContentHandler(contentService *services.ContentService) *ContentHandler {
// 	return &ContentHandler{
// 		contentService: contentService,
// 	}
// }

// // CreateContent handles content creation
// func (ch *ContentHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Header.Get("X-User-ID")
// 	if userID == "" {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Parse multipart form
// 	err := r.ParseMultipartForm(10 << 20) // 10MB max
// 	if err != nil {
// 		http.Error(w, "Failed to parse form", http.StatusBadRequest)
// 		return
// 	}

// 	eventID := r.FormValue("eventId")
// 	caption := r.FormValue("caption")
// 	tagsJSON := r.FormValue("tags")
// 	permissionsJSON := r.FormValue("permissions")

// 	// Parse tags
// 	var tags []services.ContentTag
// 	if tagsJSON != "" {
// 		json.Unmarshal([]byte(tagsJSON), &tags)
// 	}

// 	// Parse permissions
// 	var permissions services.ContentPermissions
// 	if permissionsJSON != "" {
// 		json.Unmarshal([]byte(permissionsJSON), &permissions)
// 	}

// 	// Handle file upload (simplified - in real implementation would upload to cloud storage)
// 	file, header, err := r.FormFile("media")
// 	if err != nil {
// 		http.Error(w, "No media file provided", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Simulate file upload and get URL
// 	mediaURL := "https://storage.example.com/" + header.Filename
// 	mediaType := "photo"
// 	if header.Header.Get("Content-Type") == "video/mp4" {
// 		mediaType = "video"
// 	}

// 	// Create content
// 	content, err := ch.contentService.CreateContent(userID, eventID, mediaURL, mediaType, caption, tags, permissions)
// 	if err != nil {
// 		http.Error(w, "Failed to create content", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"content": content,
// 		"message": "Content created successfully",
// 	})
// }

// // GetContent retrieves content by ID
// func (ch *ContentHandler) GetContent(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	contentID := vars["id"]

// 	content, err := ch.contentService.GetContent(contentID)
// 	if err != nil {
// 		http.Error(w, "Content not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(content)
// }

// // GetEventContent retrieves content for an event
// func (ch *ContentHandler) GetEventContent(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	eventID := vars["eventId"]
// 	brandID := r.Header.Get("X-Brand-ID")

// 	contents, err := ch.contentService.GetEventContent(eventID, brandID)
// 	if err != nil {
// 		http.Error(w, "Failed to get event content", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"contents": contents,
// 		"total":    len(contents),
// 	})
// }

// // UpdateContentPermissions updates content permissions
// func (ch *ContentHandler) UpdateContentPermissions(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	contentID := vars["id"]
// 	userID := r.Header.Get("X-User-ID")

// 	if userID == "" {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var permissions services.ContentPermissions
// 	if err := json.NewDecoder(r.Body).Decode(&permissions); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	err := ch.contentService.UpdateContentPermissions(contentID, permissions)
// 	if err != nil {
// 		http.Error(w, "Failed to update permissions", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "Permissions updated successfully",
// 	})
// }

// // SearchTags searches for available tags
// func (ch *ContentHandler) SearchTags(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query().Get("q")
// 	eventID := r.URL.Query().Get("eventId")
// 	limitStr := r.URL.Query().Get("limit")

// 	if query == "" {
// 		http.Error(w, "Query parameter required", http.StatusBadRequest)
// 		return
// 	}

// 	tags, err := ch.contentService.SearchTags(query, eventID)
// 	if err != nil {
// 		http.Error(w, "Failed to search tags", http.StatusInternalServerError)
// 		return
// 	}

// 	// Apply limit if specified
// 	if limitStr != "" {
// 		if limit, err := strconv.Atoi(limitStr); err == nil && limit < len(tags) {
// 			tags = tags[:limit]
// 		}
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"tags": tags,
// 	})
// }

// // GetSuggestedTags returns AI-suggested tags
// func (ch *ContentHandler) GetSuggestedTags(w http.ResponseWriter, r *http.Request) {
// 	// Parse multipart form for media file
// 	err := r.ParseMultipartForm(10 << 20)
// 	if err != nil {
// 		http.Error(w, "Failed to parse form", http.StatusBadRequest)
// 		return
// 	}

// 	eventID := r.FormValue("eventId")

// 	// Handle file upload (simplified)
// 	file, header, err := r.FormFile("content")
// 	if err != nil {
// 		http.Error(w, "No content file provided", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Simulate file processing and get suggested tags
// 	mediaURL := "temp://" + header.Filename
// 	suggestedTags, err := ch.contentService.GetSuggestedTags(mediaURL, eventID)
// 	if err != nil {
// 		http.Error(w, "Failed to get suggested tags", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"suggestedTags": suggestedTags,
// 	})
// }

// // // GetEventTags retrieves popular tags for an event
// func (ch *ContentHandler) GetEventTags(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	eventID := vars["eventId"]

// 	tags, err := ch.contentService.GetEventTags(eventID)
// 	if err != nil {
// 		http.Error(w, "Failed to get event tags", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"tags": tags,
// 	})
// }

// // TrackContentAnalytics tracks content interactions
// func (ch *ContentHandler) TrackContentAnalytics(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	contentID := vars["id"]

// 	var request struct {
// 		Action   string                 `json:"action"`
// 		Metadata map[string]interface{} `json:"metadata"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	err := ch.contentService.TrackContentAnalytics(contentID, request.Action, request.Metadata)
// 	if err != nil {
// 		http.Error(w, "Failed to track analytics", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "Analytics tracked successfully",
// 	})
// }

// // GetContentAnalytics retrieves content analytics
// func (ch *ContentHandler) GetContentAnalytics(w http.ResponseWriter, r *http.Request) {
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
