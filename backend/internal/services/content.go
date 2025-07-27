// /**
//  * Content Service
//  * Handles content storage, management, and analytics tracking
//  */

package services

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"time"
// )

// type ContentPermissions struct {
// 	AllowBrandAccess   bool `json:"allowBrandAccess"`
// 	AllowCommercialUse bool `json:"allowCommercialUse"`
// 	AllowModification  bool `json:"allowModification"`
// 	AllowSocialSharing bool `json:"allowSocialSharing"`
// 	ExpirationDays     int  `json:"expirationDays"`
// }

// type ContentTag struct {
// 	ID      string `json:"id"`
// 	Name    string `json:"name"`
// 	Type    string `json:"type"`
// 	BrandID string `json:"brandId,omitempty"`
// 	EventID string `json:"eventId,omitempty"`
// }

// type Content struct {
// 	ID          string             `json:"id"`
// 	UserID      string             `json:"userId"`
// 	EventID     string             `json:"eventId,omitempty"`
// 	MediaURL    string             `json:"mediaUrl"`
// 	MediaType   string             `json:"mediaType"`
// 	Caption     string             `json:"caption"`
// 	Tags        []ContentTag       `json:"tags"`
// 	Permissions ContentPermissions `json:"permissions"`
// 	CreatedAt   time.Time          `json:"createdAt"`
// 	UpdatedAt   time.Time          `json:"updatedAt"`
// }

// type ContentService struct {
// 	db *sql.DB
// }

// func NewContentService(db *sql.DB) *ContentService {
// 	return &ContentService{db: db}
// }

// // CreateContent creates new content with tags and permissions
// func (cs *ContentService) CreateContent(userID, eventID, mediaURL, mediaType, caption string, tags []ContentTag, permissions ContentPermissions) (*Content, error) {
// 	contentID := fmt.Sprintf("content_%d", time.Now().UnixNano())

// 	tagsJSON, _ := json.Marshal(tags)
// 	permissionsJSON, _ := json.Marshal(permissions)

// 	query := `
// 		INSERT INTO content (id, user_id, event_id, media_url, media_type, caption, tags, permissions, created_at, updated_at)
// 		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
// 	`

// 	now := time.Now()
// 	_, err := cs.db.Exec(query, contentID, userID, eventID, mediaURL, mediaType, caption, string(tagsJSON), string(permissionsJSON), now, now)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create content: %w", err)
// 	}

// 	return &Content{
// 		ID:          contentID,
// 		UserID:      userID,
// 		EventID:     eventID,
// 		MediaURL:    mediaURL,
// 		MediaType:   mediaType,
// 		Caption:     caption,
// 		Tags:        tags,
// 		Permissions: permissions,
// 		CreatedAt:   now,
// 		UpdatedAt:   now,
// 	}, nil
// }

// // GetContent retrieves content by ID
// func (cs *ContentService) GetContent(contentID string) (*Content, error) {
// 	query := `
// 		SELECT id, user_id, event_id, media_url, media_type, caption, tags, permissions, created_at, updated_at
// 		FROM content WHERE id = ?
// 	`

// 	var content Content
// 	var tagsJSON, permissionsJSON string

// 	err := cs.db.QueryRow(query, contentID).Scan(
// 		&content.ID, &content.UserID, &content.EventID, &content.MediaURL,
// 		&content.MediaType, &content.Caption, &tagsJSON, &permissionsJSON,
// 		&content.CreatedAt, &content.UpdatedAt,
// 	)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get content: %w", err)
// 	}

// 	json.Unmarshal([]byte(tagsJSON), &content.Tags)
// 	json.Unmarshal([]byte(permissionsJSON), &content.Permissions)

// 	return &content, nil
// }

// // GetEventContent retrieves content for a specific event
// func (cs *ContentService) GetEventContent(eventID string, brandID string) ([]Content, error) {
// 	query := `
// 		SELECT c.id, c.user_id, c.event_id, c.media_url, c.media_type, c.caption, c.tags, c.permissions, c.created_at, c.updated_at
// 		FROM content c
// 		WHERE c.event_id = ? AND JSON_EXTRACT(c.permissions, '$.allowBrandAccess') = true
// 		ORDER BY c.created_at DESC
// 	`

// 	rows, err := cs.db.Query(query, eventID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get event content: %w", err)
// 	}
// 	defer rows.Close()

// 	var contents []Content
// 	for rows.Next() {
// 		var content Content
// 		var tagsJSON, permissionsJSON string

// 		err := rows.Scan(
// 			&content.ID, &content.UserID, &content.EventID, &content.MediaURL,
// 			&content.MediaType, &content.Caption, &tagsJSON, &permissionsJSON,
// 			&content.CreatedAt, &content.UpdatedAt,
// 		)
// 		if err != nil {
// 			continue
// 		}

// 		json.Unmarshal([]byte(tagsJSON), &content.Tags)
// 		json.Unmarshal([]byte(permissionsJSON), &content.Permissions)

// 		contents = append(contents, content)
// 	}

// 	return contents, nil
// }

// // UpdateContentPermissions updates permissions for existing content
// func (cs *ContentService) UpdateContentPermissions(contentID string, permissions ContentPermissions) error {
// 	permissionsJSON, _ := json.Marshal(permissions)

// 	query := `UPDATE content SET permissions = ?, updated_at = ? WHERE id = ?`
// 	_, err := cs.db.Exec(query, string(permissionsJSON), time.Now(), contentID)

// 	if err != nil {
// 		return fmt.Errorf("failed to update content permissions: %w", err)
// 	}

// 	return nil
// }

// // SearchTags searches for available tags
// func (cs *ContentService) SearchTags(query, eventID string) ([]ContentTag, error) {
// 	sqlQuery := `
// 		SELECT DISTINCT name, type, brand_id, event_id
// 		FROM content_tags
// 		WHERE name LIKE ? AND (event_id = ? OR event_id IS NULL)
// 		ORDER BY name
// 		LIMIT 20
// 	`

// 	rows, err := cs.db.Query(sqlQuery, "%"+query+"%", eventID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to search tags: %w", err)
// 	}
// 	defer rows.Close()

// 	var tags []ContentTag
// 	for rows.Next() {
// 		var tag ContentTag
// 		var brandID, eventIDVal sql.NullString

// 		err := rows.Scan(&tag.Name, &tag.Type, &brandID, &eventIDVal)
// 		if err != nil {
// 			continue
// 		}

// 		tag.ID = fmt.Sprintf("tag_%s_%s", tag.Type, tag.Name)
// 		if brandID.Valid {
// 			tag.BrandID = brandID.String
// 		}
// 		if eventIDVal.Valid {
// 			tag.EventID = eventIDVal.String
// 		}

// 		tags = append(tags, tag)
// 	}

// 	return tags, nil
// }

// // GetSuggestedTags returns AI-suggested tags for content
// func (cs *ContentService) GetSuggestedTags(mediaURL, eventID string) ([]ContentTag, error) {
// 	// Simulate AI tagging - in real implementation, this would call ML service
// 	suggestedTags := []ContentTag{
// 		{ID: "tag_brand_sample", Name: "Sample Brand", Type: "brand"},
// 		{ID: "tag_product_demo", Name: "Product Demo", Type: "product"},
// 		{ID: "tag_event_activation", Name: "Brand Activation", Type: "event"},
// 	}

// 	// Add event-specific tags if eventID provided
// 	if eventID != "" {
// 		eventTags, _ := cs.GetEventTags(eventID)
// 		suggestedTags = append(suggestedTags, eventTags...)
// 	}

// 	return suggestedTags, nil
// }

// // GetEventTags retrieves popular tags for an event
// func (cs *ContentService) GetEventTags(eventID string) ([]ContentTag, error) {
// 	query := `
// 		SELECT name, type, COUNT(*) as usage_count
// 		FROM content_tags
// 		WHERE event_id = ?
// 		GROUP BY name, type
// 		ORDER BY usage_count DESC
// 		LIMIT 10
// 	`

// 	rows, err := cs.db.Query(query, eventID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get event tags: %w", err)
// 	}
// 	defer rows.Close()

// 	var tags []ContentTag
// 	for rows.Next() {
// 		var tag ContentTag
// 		var usageCount int

// 		err := rows.Scan(&tag.Name, &tag.Type, &usageCount)
// 		if err != nil {
// 			continue
// 		}

// 		tag.ID = fmt.Sprintf("tag_%s_%s", tag.Type, tag.Name)
// 		tag.EventID = eventID

// 		tags = append(tags, tag)
// 	}

// 	return tags, nil
// }

// // TrackContentAnalytics tracks content performance metrics
// func (cs *ContentService) TrackContentAnalytics(contentID, action string, metadata map[string]interface{}) error {
// 	metadataJSON, _ := json.Marshal(metadata)

// 	query := `
// 		INSERT INTO content_analytics (content_id, action, metadata, created_at)
// 		VALUES (?, ?, ?, ?)
// 	`

// 	_, err := cs.db.Exec(query, contentID, action, string(metadataJSON), time.Now())
// 	if err != nil {
// 		return fmt.Errorf("failed to track content analytics: %w", err)
// 	}

// 	return nil
// }

// // GetContentAnalytics retrieves analytics for content
// func (cs *ContentService) GetContentAnalytics(contentID string) (map[string]interface{}, error) {
// 	query := `
// 		SELECT action, COUNT(*) as count
// 		FROM content_analytics
// 		WHERE content_id = ?
// 		GROUP BY action
// 	`

// 	rows, err := cs.db.Query(query, contentID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get content analytics: %w", err)
// 	}
// 	defer rows.Close()

// 	analytics := make(map[string]interface{})
// 	for rows.Next() {
// 		var action string
// 		var count int

// 		err := rows.Scan(&action, &count)
// 		if err != nil {
// 			continue
// 		}

// 		analytics[action] = count
// 	}

// 	return analytics, nil
// }
