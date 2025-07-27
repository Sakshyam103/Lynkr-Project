// package content

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"time"
// )

// // ContentType represents the type of content
// type ContentType string

// const (
// 	ContentTypePhoto ContentType = "photo"
// 	ContentTypeVideo ContentType = "video"
// 	ContentTypeText  ContentType = "text"
// )

// // Content represents user-generated content
// type Content struct {
// 	ID          uint       `json:"id"`
// 	UserID      uint       `json:"user_id"`
// 	EventID     uint       `json:"event_id"`
// 	Type        ContentType `json:"type"`
// 	URL         string     `json:"url"`
// 	Permissions string     `json:"permissions"`
// 	CreatedAt   time.Time  `json:"created_at"`
// }

// // Interaction represents a user's interaction with content
// type Interaction struct {
// 	ID        uint      `json:"id"`
// 	UserID    uint      `json:"user_id"`
// 	ContentID uint      `json:"content_id"`
// 	Type      string    `json:"type"`
// 	Data      string    `json:"data"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// // ContentService handles content-related operations
// type ContentService struct {
// 	DB *sql.DB
// }

// // NewContentService creates a new content service
// func NewContentService(db *sql.DB) *ContentService {
// 	return &ContentService{DB: db}
// }

// // Create creates new content
// func (s *ContentService) Create(userID, eventID uint, contentType ContentType, url string, permissions map[string]interface{}) (*Content, error) {
// 	// Convert permissions to JSON
// 	permissionsJSON, err := json.Marshal(permissions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := `
// 		INSERT INTO content (user_id, event_id, type, url, permissions)
// 		VALUES (?, ?, ?, ?, ?)
// 		RETURNING id, user_id, event_id, type, url, permissions, created_at
// 	`

// 	var content Content
// 	err = s.DB.QueryRow(
// 		query,
// 		userID,
// 		eventID,
// 		contentType,
// 		url,
// 		string(permissionsJSON),
// 	).Scan(
// 		&content.ID,
// 		&content.UserID,
// 		&content.EventID,
// 		&content.Type,
// 		&content.URL,
// 		&content.Permissions,
// 		&content.CreatedAt,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &content, nil
// }

// // GetByID retrieves content by ID
// func (s *ContentService) GetByID(id uint) (*Content, error) {
// 	query := `
// 		SELECT id, user_id, event_id, type, url, permissions, created_at
// 		FROM content
// 		WHERE id = ?
// 	`

// 	var content Content
// 	err := s.DB.QueryRow(query, id).Scan(
// 		&content.ID,
// 		&content.UserID,
// 		&content.EventID,
// 		&content.Type,
// 		&content.URL,
// 		&content.Permissions,
// 		&content.CreatedAt,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("content not found")
// 		}
// 		return nil, err
// 	}

// 	return &content, nil
// }

// // GetByEvent retrieves content for an event
// func (s *ContentService) GetByEvent(eventID uint) ([]Content, error) {
// 	query := `
// 		SELECT id, user_id, event_id, type, url, permissions, created_at
// 		FROM content
// 		WHERE event_id = ?
// 		ORDER BY created_at DESC
// 	`

// 	rows, err := s.DB.Query(query, eventID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var contents []Content
// 	for rows.Next() {
// 		var content Content
// 		err := rows.Scan(
// 			&content.ID,
// 			&content.UserID,
// 			&content.EventID,
// 			&content.Type,
// 			&content.URL,
// 			&content.Permissions,
// 			&content.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		contents = append(contents, content)
// 	}

// 	return contents, nil
// }

// // UpdatePermissions updates content permissions
// func (s *ContentService) UpdatePermissions(contentID uint, permissions map[string]interface{}) error {
// 	// Convert permissions to JSON
// 	permissionsJSON, err := json.Marshal(permissions)
// 	if err != nil {
// 		return err
// 	}

// 	query := `
// 		UPDATE content
// 		SET permissions = ?
// 		WHERE id = ?
// 	`

// 	result, err := s.DB.Exec(query, string(permissionsJSON), contentID)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("content not found")
// 	}

// 	return nil
// }

// // AddInteraction records a user's interaction with content
// func (s *ContentService) AddInteraction(userID, contentID uint, interactionType string, data map[string]interface{}) (*Interaction, error) {
// 	// Convert data to JSON
// 	dataJSON, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	query := `
// 		INSERT INTO interactions (user_id, content_id, type, data)
// 		VALUES (?, ?, ?, ?)
// 		RETURNING id, user_id, content_id, type, data, created_at
// 	`

// 	var interaction Interaction
// 	err = s.DB.QueryRow(
// 		query,
// 		userID,
// 		contentID,
// 		interactionType,
// 		string(dataJSON),
// 	).Scan(
// 		&interaction.ID,
// 		&interaction.UserID,
// 		&interaction.ContentID,
// 		&interaction.Type,
// 		&interaction.Data,
// 		&interaction.CreatedAt,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &interaction, nil
// }

// // GetInteractions retrieves interactions for content
// func (s *ContentService) GetInteractions(contentID uint) ([]Interaction, error) {
// 	query := `
// 		SELECT id, user_id, content_id, type, data, created_at
// 		FROM interactions
// 		WHERE content_id = ?
// 		ORDER BY created_at DESC
// 	`

// 	rows, err := s.DB.Query(query, contentID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var interactions []Interaction
// 	for rows.Next() {
// 		var interaction Interaction
// 		err := rows.Scan(
// 			&interaction.ID,
// 			&interaction.UserID,
// 			&interaction.ContentID,
// 			&interaction.Type,
// 			&interaction.Data,
// 			&interaction.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		interactions = append(interactions, interaction)
// 	}

// 	return interactions, nil
// }

//  // GetEventTags retrieves popular tags for an event
//  func (cs *ContentService) GetEventTags(eventID string) ([]ContentTag, error) {
// 	 query := `
// 		 SELECT name, type, COUNT(*) as usage_count
// 		 FROM content_tags
// 		 WHERE event_id = ?
// 		 GROUP BY name, type
// 		 ORDER BY usage_count DESC
// 		 LIMIT 10
// 	 `

// 	 rows, err := cs.db.Query(query, eventID)
// 	 if err != nil {
// 		 return nil, fmt.Errorf("failed to get event tags: %w", err)
// 	 }
// 	 defer rows.Close()

// 	 var tags []ContentTag
// 	 for rows.Next() {
// 		 var tag ContentTag
// 		 var usageCount int

// 		 err := rows.Scan(&tag.Name, &tag.Type, &usageCount)
// 		 if err != nil {
// 			 continue
// 		 }

// 		 tag.ID = fmt.Sprintf("tag_%s_%s", tag.Type, tag.Name)
// 		 tag.EventID = eventID

// 		 tags = append(tags, tag)
// 	 }

// 	 return tags, nil
//  }

// //  // TrackContentAnalytics tracks content performance metrics
// //  func (cs *ContentService) TrackContentAnalytics(contentID, action string, metadata map[string]interface{}) error {
// // 	 metadataJSON, _ := json.Marshal(metadata)

// // 	 query := `
// // 		 INSERT INTO content_analytics (content_id, action, metadata, created_at)
// // 		 VALUES (?, ?, ?, ?)
// // 	 `

// // 	 _, err := cs.db.Exec(query, contentID, action, string(metadataJSON), time.Now())
// // 	 if err != nil {
// // 		 return fmt.Errorf("failed to track content analytics: %w", err)
// // 	 }

// // 	 return nil
// //  }

/**
 * Content Service
 * Handles content storage, management, and analytics tracking
 */

package content

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type ContentPermissions struct {
	AllowBrandAccess   bool `json:"allowBrandAccess"`
	AllowCommercialUse bool `json:"allowCommercialUse"`
	AllowModification  bool `json:"allowModification"`
	AllowSocialSharing bool `json:"allowSocialSharing"`
	ExpirationDays     int  `json:"expirationDays"`
}

type ContentTag struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	BrandID string `json:"brandId,omitempty"`
	EventID string `json:"eventId,omitempty"`
}

type Content struct {
	ID          string             `json:"id"`
	UserID      uint               `json:"userId"`
	EventID     int                `json:"eventId,omitempty"`
	MediaURL    string             `json:"url"`
	MediaType   string             `json:"type"`
	Caption     string             `json:"caption"`
	Tags        []ContentTag       `json:"tags"`
	Permissions ContentPermissions `json:"permissions"`
	CreatedAt   time.Time          `json:"createdAt"`
	//UpdatedAt   time.Time          `json:"updatedAt"`
}

type Interaction struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ContentID uint      `json:"content_id"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type ContentService struct {
	db *sql.DB
}

func NewContentService(db *sql.DB) *ContentService {
	return &ContentService{db: db}
}

// CreateContent creates new content with tags and permissions
func (cs *ContentService) CreateContent(userID uint, eventID int, mediaURL, mediaType, caption string, tags []ContentTag, permissions ContentPermissions) (*Content, error) {
	contentID := fmt.Sprintf("content_%d", time.Now().UnixNano())

	tagsJSON, _ := json.Marshal(tags)
	permissionsJSON, _ := json.Marshal(permissions)

	query := `
		 INSERT INTO content (user_id, event_id, type, url, permissions, created_at, caption, tags)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	 `

	now := time.Now()
	//tagsJSONBytes, err := json.Marshal(tags)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to marshal tags: %w", err)
	//}
	//
	//permissionsJSONBytes, err := json.Marshal(permissions)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to marshal permissions: %w", err)
	//}
	_, err := cs.db.Exec(query, userID, eventID, mediaType, mediaURL, permissionsJSON, now, caption, tagsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}

	var permissions1 ContentPermissions
	err = json.Unmarshal([]byte(permissionsJSON), &permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions JSON: %w", err)
	}

	var tags1 []ContentTag
	err = json.Unmarshal([]byte(tagsJSON), &tags1)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions JSON: %w", err)
	}

	return &Content{
		ID:          contentID,
		UserID:      userID,
		EventID:     eventID,
		MediaURL:    mediaURL,
		MediaType:   mediaType,
		Caption:     caption,
		Tags:        tags1,
		Permissions: permissions1,
		CreatedAt:   now,
	}, nil
}

// GetContent retrieves content by ID
func (cs *ContentService) GetContent(contentID string) (*Content, error) {
	query := `
		 SELECT id, user_id, event_id, media_url, media_type, caption, tags, permissions, created_at
		 FROM content WHERE id = ?
	 `

	var content Content
	var tagsJSON, permissionsJSON string

	err := cs.db.QueryRow(query, contentID).Scan(
		&content.ID, &content.UserID, &content.EventID, &content.MediaURL,
		&content.MediaType, &content.Caption, &tagsJSON, &permissionsJSON,
		&content.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	json.Unmarshal([]byte(tagsJSON), &content.Tags)
	json.Unmarshal([]byte(permissionsJSON), &content.Permissions)

	return &content, nil
}

// GetEventContent retrieves content for a specific event
func (cs *ContentService) GetEventContent(eventID int) ([]Content, error) {
	query := `
		 SELECT c.id, c.user_id, c.event_id, c.url, c.type, c.caption, c.tags, c.permissions, c.created_at
		 FROM content c
		 WHERE c.event_id = ? AND JSON_EXTRACT(c.permissions, '$.allowBrandAccess') = 1
		 ORDER BY c.created_at DESC
	 `

	rows, err := cs.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event content: %w", err)
	}
	defer rows.Close()

	var contents []Content
	for rows.Next() {
		var content Content
		var tagsJSON, permissionsJSON string

		err := rows.Scan(
			&content.ID, &content.UserID, &content.EventID, &content.MediaURL,
			&content.MediaType, &content.Caption, &tagsJSON, &permissionsJSON,
			&content.CreatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(tagsJSON), &content.Tags)
		json.Unmarshal([]byte(permissionsJSON), &content.Permissions)

		contents = append(contents, content)
	}

	return contents, nil
}

// UpdateContentPermissions updates permissions for existing content
func (cs *ContentService) UpdateContentPermissions(contentID string, permissions ContentPermissions) error {
	permissionsJSON, _ := json.Marshal(permissions)

	query := `UPDATE content SET permissions = ?, updated_at = ? WHERE id = ?`
	_, err := cs.db.Exec(query, string(permissionsJSON), time.Now(), contentID)

	if err != nil {
		return fmt.Errorf("failed to update content permissions: %w", err)
	}

	return nil
}

// SearchTags searches for available tags
func (cs *ContentService) SearchTags(query, eventID string) ([]ContentTag, error) {
	sqlQuery := `
		 SELECT DISTINCT name, type, brand_id, event_id
		 FROM content_tags
		 WHERE name LIKE ? AND (event_id = ? OR event_id IS NULL)
		 ORDER BY name
		 LIMIT 20
	 `

	rows, err := cs.db.Query(sqlQuery, "%"+query+"%", eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to search tags: %w", err)
	}
	defer rows.Close()

	var tags []ContentTag
	for rows.Next() {
		var tag ContentTag
		var brandID, eventIDVal sql.NullString

		err := rows.Scan(&tag.Name, &tag.Type, &brandID, &eventIDVal)
		if err != nil {
			continue
		}

		tag.ID = fmt.Sprintf("tag_%s_%s", tag.Type, tag.Name)
		if brandID.Valid {
			tag.BrandID = brandID.String
		}
		if eventIDVal.Valid {
			tag.EventID = eventIDVal.String
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

// GetSuggestedTags returns AI-suggested tags for content
func (cs *ContentService) GetSuggestedTags(mediaURL, eventID string) ([]ContentTag, error) {
	// Simulate AI tagging - in real implementation, this would call ML service
	suggestedTags := []ContentTag{
		{ID: "tag_brand_sample", Name: "Sample Brand", Type: "brand"},
		{ID: "tag_product_demo", Name: "Product Demo", Type: "product"},
		{ID: "tag_event_activation", Name: "Brand Activation", Type: "event"},
	}

	// Add event-specific tags if eventID provided
	if eventID != "" {
		eventTags, _ := cs.GetEventTags(eventID)
		suggestedTags = append(suggestedTags, eventTags...)
	}

	return suggestedTags, nil
}

// GetEventTags retrieves popular tags for an event
func (cs *ContentService) GetEventTags(eventID string) ([]ContentTag, error) {
	query := `
		 SELECT name, type, COUNT(*) as usage_count
		 FROM content_tags
		 WHERE event_id = ?
		 GROUP BY name, type
		 ORDER BY usage_count DESC
		 LIMIT 10
	 `

	rows, err := cs.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event tags: %w", err)
	}
	defer rows.Close()

	var tags []ContentTag
	for rows.Next() {
		var tag ContentTag
		var usageCount int

		err := rows.Scan(&tag.Name, &tag.Type, &usageCount)
		if err != nil {
			continue
		}

		tag.ID = fmt.Sprintf("tag_%s_%s", tag.Type, tag.Name)
		tag.EventID = eventID

		tags = append(tags, tag)
	}

	return tags, nil
}

// TrackContentAnalytics tracks content performance metrics
func (cs *ContentService) TrackContentAnalytics(contentID, action string, metadata map[string]interface{}) error {
	metadataJSON, _ := json.Marshal(metadata)

	query := `
		 INSERT INTO content_analytics (content_id, action, metadata, created_at)
		 VALUES (?, ?, ?, ?)
	 `

	_, err := cs.db.Exec(query, contentID, action, string(metadataJSON), time.Now())
	if err != nil {
		return fmt.Errorf("failed to track content analytics: %w", err)
	}

	return nil
}

// GetContentAnalytics retrieves analytics for content
func (cs *ContentService) GetContentAnalytics(contentID string) (map[string]interface{}, error) {
	query := `
		 SELECT action, COUNT(*) as count
		 FROM content_analytics
		 WHERE content_id = ?
		 GROUP BY action
	 `

	rows, err := cs.db.Query(query, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get content analytics: %w", err)
	}
	defer rows.Close()

	analytics := make(map[string]interface{})
	for rows.Next() {
		var action string
		var count int

		err := rows.Scan(&action, &count)
		if err != nil {
			continue
		}

		analytics[action] = count
	}

	return analytics, nil
}

// AddInteraction records a user's interaction with content
func (s *ContentService) AddInteraction(userID, contentID uint, interactionType string, data map[string]interface{}) (*Interaction, error) {
	// Convert data to JSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO interactions (user_id, content_id, type, data)
		VALUES (?, ?, ?, ?)
		RETURNING id, user_id, content_id, type, data, created_at
	`

	var interaction Interaction
	err = s.db.QueryRow(
		query,
		userID,
		contentID,
		interactionType,
		string(dataJSON),
	).Scan(
		&interaction.ID,
		&interaction.UserID,
		&interaction.ContentID,
		&interaction.Type,
		&interaction.Data,
		&interaction.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &interaction, nil
}

// GetInteractions retrieves interactions for content
func (s *ContentService) GetInteractions(contentID uint) ([]Interaction, error) {
	query := `
		SELECT id, user_id, content_id, type, data, created_at
		FROM interactions
		WHERE content_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, contentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interactions []Interaction
	for rows.Next() {
		var interaction Interaction
		err := rows.Scan(
			&interaction.ID,
			&interaction.UserID,
			&interaction.ContentID,
			&interaction.Type,
			&interaction.Data,
			&interaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		interactions = append(interactions, interaction)
	}

	return interactions, nil
}
