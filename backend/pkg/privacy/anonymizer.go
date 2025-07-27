package privacy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Anonymizer provides methods for anonymizing user data
type Anonymizer struct {
	salt      string
	retention map[string]time.Duration
}

// NewAnonymizer creates a new anonymizer with the given salt
func NewAnonymizer(salt string) *Anonymizer {
	// Default retention periods
	retention := map[string]time.Duration{
		"30days":  30 * 24 * time.Hour,
		"90days":  90 * 24 * time.Hour,
		"1year":   365 * 24 * time.Hour,
		"forever": 100 * 365 * 24 * time.Hour, // Effectively forever
	}

	return &Anonymizer{
		salt:      salt,
		retention: retention,
	}
}

// AnonymizeUserID hashes a user ID to create a pseudonymous identifier
func (a *Anonymizer) AnonymizeUserID(userID uint) string {
	data := fmt.Sprintf("%d:%s", userID, a.salt)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// AnonymizeEmail partially masks an email address
func (a *Anonymizer) AnonymizeEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "invalid@example.com"
	}

	username := parts[0]
	domain := parts[1]

	// Keep first and last character, mask the rest
	if len(username) > 2 {
		masked := username[0:1] + strings.Repeat("*", len(username)-2) + username[len(username)-1:]
		return masked + "@" + domain
	}

	// If username is too short, just mask one character
	if len(username) == 2 {
		return username[0:1] + "*@" + domain
	}

	// If username is just one character, add a mask
	return username + "*@" + domain
}

// AnonymizeLocation reduces precision of location data
func (a *Anonymizer) AnonymizeLocation(lat, lng float64) (float64, float64) {
	// Reduce precision to ~1km (roughly 0.01 degrees)
	return float64(int(lat*100)) / 100, float64(int(lng*100)) / 100
}

// AnonymizeData anonymizes a map of user data based on privacy settings
func (a *Anonymizer) AnonymizeData(data map[string]interface{}, privacySettings map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Copy data that doesn't need anonymization
	for k, v := range data {
		result[k] = v
	}

	// Apply anonymization based on privacy settings
	settings := getPrivacySettings(privacySettings)

	// Anonymize location if not allowed to share
	if !settings.ShareLocation {
		if location, ok := result["location"].(map[string]interface{}); ok {
			if lat, ok := location["latitude"].(float64); ok {
				if lng, ok := location["longitude"].(float64); ok {
					anonLat, anonLng := a.AnonymizeLocation(lat, lng)
					location["latitude"] = anonLat
					location["longitude"] = anonLng
					result["location"] = location
				}
			}
		}
	}

	// Anonymize email
	if email, ok := result["email"].(string); ok {
		result["email"] = a.AnonymizeEmail(email)
	}

	// Remove analytics data if not allowed
	if !settings.ShareAnalytics {
		delete(result, "analytics")
		delete(result, "usage_data")
	}

	return result
}

// ShouldRetainData checks if data should be retained based on privacy settings
func (a *Anonymizer) ShouldRetainData(createdAt time.Time, privacySettings map[string]interface{}) bool {
	settings := getPrivacySettings(privacySettings)
	
	// Get retention period
	retention, ok := a.retention[settings.DataRetention]
	if !ok {
		// Default to 90 days if not specified
		retention = a.retention["90days"]
	}
	
	// Check if data is within retention period
	return time.Since(createdAt) <= retention
}

// PrivacySettings represents user privacy preferences
type PrivacySettings struct {
	ShareLocation      bool   `json:"shareLocation"`
	ShareContent       bool   `json:"shareContent"`
	AllowNotifications bool   `json:"allowNotifications"`
	DataRetention      string `json:"dataRetention"`
	ShareAnalytics     bool   `json:"shareAnalytics"`
	PersonalizedContent bool  `json:"personalizedContent"`
}

// getPrivacySettings converts a map to PrivacySettings struct
func getPrivacySettings(settings map[string]interface{}) PrivacySettings {
	// Default settings
	result := PrivacySettings{
		ShareLocation:      false,
		ShareContent:       true,
		AllowNotifications: true,
		DataRetention:      "90days",
		ShareAnalytics:     true,
		PersonalizedContent: true,
	}
	
	// Convert map to JSON and then to struct
	if settings != nil {
		jsonData, err := json.Marshal(settings)
		if err == nil {
			json.Unmarshal(jsonData, &result)
		}
	}
	
	return result
}