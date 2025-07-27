/**
 * CRM Integration Service
 * Handles CRM system integrations and data synchronization
 */

package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CRMIntegration struct {
	ID           string `json:"id"`
	BrandID      string `json:"brandId"`
	CRMType      string `json:"crmType"`
	APIKey       string `json:"apiKey"`
	APISecret    string `json:"apiSecret"`
	WebhookURL   string `json:"webhookUrl"`
	SyncInterval int    `json:"syncInterval"`
	Status       string `json:"status"`
}

type CRMContact struct {
	ID       string                 `json:"id"`
	Email    string                 `json:"email"`
	Name     string                 `json:"name"`
	EventID  string                 `json:"eventId"`
	Metadata map[string]interface{} `json:"metadata"`
}

type CRMIntegrationService struct {
	db *sql.DB
}

func NewCRMIntegrationService(db *sql.DB) *CRMIntegrationService {
	return &CRMIntegrationService{db: db}
}

func (cis *CRMIntegrationService) CreateIntegration(brandID, crmType, apiKey, apiSecret, webhookURL string, syncInterval int) (*CRMIntegration, error) {
	integrationID := fmt.Sprintf("crm_%d", time.Now().UnixNano())
	
	query := `
		INSERT INTO crm_integrations (id, brand_id, crm_type, api_key, api_secret, webhook_url, sync_interval, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := cis.db.Exec(query, integrationID, brandID, crmType, apiKey, apiSecret, webhookURL, syncInterval, "active", time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create CRM integration: %w", err)
	}
	
	return &CRMIntegration{
		ID:           integrationID,
		BrandID:      brandID,
		CRMType:      crmType,
		APIKey:       apiKey,
		APISecret:    apiSecret,
		WebhookURL:   webhookURL,
		SyncInterval: syncInterval,
		Status:       "active",
	}, nil
}

func (cis *CRMIntegrationService) SyncEventData(integrationID, eventID string) error {
	integration, err := cis.getIntegration(integrationID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}
	
	contacts, err := cis.getEventContacts(eventID)
	if err != nil {
		return fmt.Errorf("failed to get event contacts: %w", err)
	}
	
	switch integration.CRMType {
	case "salesforce":
		return cis.syncToSalesforce(integration, contacts)
	case "hubspot":
		return cis.syncToHubspot(integration, contacts)
	case "mailchimp":
		return cis.syncToMailchimp(integration, contacts)
	default:
		return fmt.Errorf("unsupported CRM type: %s", integration.CRMType)
	}
}

func (cis *CRMIntegrationService) syncToSalesforce(integration *CRMIntegration, contacts []CRMContact) error {
	for _, contact := range contacts {
		payload := map[string]interface{}{
			"Email":     contact.Email,
			"FirstName": contact.Name,
			"EventId__c": contact.EventID,
			"Source":    "Lynkr Event",
		}
		
		err := cis.sendToCRM("https://api.salesforce.com/services/data/v52.0/sobjects/Contact/", integration, payload)
		if err != nil {
			return fmt.Errorf("failed to sync contact to Salesforce: %w", err)
		}
	}
	return nil
}

func (cis *CRMIntegrationService) syncToHubspot(integration *CRMIntegration, contacts []CRMContact) error {
	for _, contact := range contacts {
		payload := map[string]interface{}{
			"properties": map[string]interface{}{
				"email":     contact.Email,
				"firstname": contact.Name,
				"event_id":  contact.EventID,
				"source":    "Lynkr Event",
			},
		}
		
		err := cis.sendToCRM("https://api.hubapi.com/crm/v3/objects/contacts", integration, payload)
		if err != nil {
			return fmt.Errorf("failed to sync contact to HubSpot: %w", err)
		}
	}
	return nil
}

func (cis *CRMIntegrationService) syncToMailchimp(integration *CRMIntegration, contacts []CRMContact) error {
	for _, contact := range contacts {
		payload := map[string]interface{}{
			"email_address": contact.Email,
			"status":        "subscribed",
			"merge_fields": map[string]interface{}{
				"FNAME":    contact.Name,
				"EVENT_ID": contact.EventID,
			},
			"tags": []string{"lynkr-event"},
		}
		
		err := cis.sendToCRM("https://us1.api.mailchimp.com/3.0/lists/LIST_ID/members", integration, payload)
		if err != nil {
			return fmt.Errorf("failed to sync contact to Mailchimp: %w", err)
		}
	}
	return nil
}

func (cis *CRMIntegrationService) sendToCRM(url string, integration *CRMIntegration, payload map[string]interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+integration.APIKey)
	
	// Add signature for security
	signature := cis.generateSignature(string(jsonData), integration.APISecret)
	req.Header.Set("X-Lynkr-Signature", signature)
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("CRM API error: %d", resp.StatusCode)
	}
	
	return nil
}

func (cis *CRMIntegrationService) getIntegration(integrationID string) (*CRMIntegration, error) {
	query := `
		SELECT id, brand_id, crm_type, api_key, api_secret, webhook_url, sync_interval, status
		FROM crm_integrations WHERE id = ?
	`
	
	var integration CRMIntegration
	err := cis.db.QueryRow(query, integrationID).Scan(
		&integration.ID, &integration.BrandID, &integration.CRMType,
		&integration.APIKey, &integration.APISecret, &integration.WebhookURL,
		&integration.SyncInterval, &integration.Status,
	)
	
	if err != nil {
		return nil, fmt.Errorf("integration not found: %w", err)
	}
	
	return &integration, nil
}

func (cis *CRMIntegrationService) getEventContacts(eventID string) ([]CRMContact, error) {
	query := `
		SELECT u.id, u.email, u.name, a.event_id
		FROM users u
		JOIN attendances a ON u.id = a.user_id
		WHERE a.event_id = ? AND u.consent_marketing = 1
	`
	
	rows, err := cis.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var contacts []CRMContact
	for rows.Next() {
		var contact CRMContact
		err := rows.Scan(&contact.ID, &contact.Email, &contact.Name, &contact.EventID)
		if err != nil {
			continue
		}
		
		contact.Metadata = map[string]interface{}{
			"source":      "lynkr_event",
			"sync_date":   time.Now().Format(time.RFC3339),
		}
		
		contacts = append(contacts, contact)
	}
	
	return contacts, nil
}

func (cis *CRMIntegrationService) generateSignature(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func (cis *CRMIntegrationService) ScheduleSync(integrationID string) error {
	integration, err := cis.getIntegration(integrationID)
	if err != nil {
		return err
	}
	
	// Schedule sync based on interval
	go func() {
		ticker := time.NewTicker(time.Duration(integration.SyncInterval) * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			// Get recent events for this brand
			events := cis.getRecentEvents(integration.BrandID)
			for _, eventID := range events {
				cis.SyncEventData(integrationID, eventID)
			}
		}
	}()
	
	return nil
}

func (cis *CRMIntegrationService) getRecentEvents(brandID string) []string {
	query := `
		SELECT id FROM events 
		WHERE brand_id = ? AND created_at > datetime('now', '-1 day')
	`
	
	rows, err := cis.db.Query(query, brandID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()
	
	var events []string
	for rows.Next() {
		var eventID string
		if err := rows.Scan(&eventID); err == nil {
			events = append(events, eventID)
		}
	}
	
	return events
}