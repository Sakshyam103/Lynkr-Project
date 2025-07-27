/**
 * E-commerce Service
 * Core service for e-commerce integrations and purchase tracking
 */

package services

import (
	"database/sql"
	"fmt"
	"time"
)

type EcommerceService struct {
	db *sql.DB
}

type Integration struct {
	ID           string `json:"id"`
	BrandID      string `json:"brandId"`
	PlatformType string `json:"platformType"`
	APIKey       string `json:"apiKey"`
	StoreURL     string `json:"storeUrl"`
	WebhookURL   string `json:"webhookUrl"`
	Status       string `json:"status"`
}

type Purchase struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	ProductID string    `json:"productId"`
	EventID   string    `json:"eventId"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewEcommerceService(db *sql.DB) *EcommerceService {
	return &EcommerceService{db: db}
}

func (es *EcommerceService) CreateIntegration(brandID, platformType, apiKey, storeURL string) (*Integration, error) {
	integrationID := fmt.Sprintf("integration_%d", time.Now().UnixNano())
	webhookURL := fmt.Sprintf("https://api.lynkr.com/webhooks/%s", integrationID)

	query := `
		INSERT INTO ecommerce_integrations (id, brand_id, platform_type, api_key, store_url, webhook_url, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := es.db.Exec(query, integrationID, brandID, platformType, apiKey, storeURL, webhookURL, "active", time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}

	return &Integration{
		ID:           integrationID,
		BrandID:      brandID,
		PlatformType: platformType,
		APIKey:       apiKey,
		StoreURL:     storeURL,
		WebhookURL:   webhookURL,
		Status:       "active",
	}, nil
}

func (es *EcommerceService) GetIntegration(brandID string) (*Integration, error) {
	query := `
		SELECT id, brand_id, platform_type, api_key, store_url, webhook_url, status
		FROM ecommerce_integrations WHERE brand_id = ? AND status = 'active'
	`

	var integration Integration
	err := es.db.QueryRow(query, brandID).Scan(
		&integration.ID, &integration.BrandID, &integration.PlatformType,
		&integration.APIKey, &integration.StoreURL, &integration.WebhookURL,
		&integration.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	return &integration, nil
}

func (es *EcommerceService) TrackPurchase(userID, productID, eventID string, amount float64) (*Purchase, error) {
	purchaseID := fmt.Sprintf("purchase_%d", time.Now().UnixNano())

	query := `
		INSERT INTO purchases (id, user_id, product_id, event_id, amount, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := es.db.Exec(query, purchaseID, userID, productID, eventID, amount, "completed", time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to track purchase: %w", err)
	}

	return &Purchase{
		ID:        purchaseID,
		UserID:    userID,
		ProductID: productID,
		EventID:   eventID,
		Amount:    amount,
		Status:    "completed",
		CreatedAt: time.Now(),
	}, nil
}

func (es *EcommerceService) GetPurchaseAnalytics(eventID string) (map[string]interface{}, error) {
	query := `
-- 		SELECT 
-- 			COUNT(*) as total_purchases,
-- 			SUM(amount) as total_revenue,
-- 			AVG(amount) as avg_order_value,
-- 			COUNT(DISTINCT user_id) as unique_buyers
-- 		FROM purchases 
-- 		WHERE event_id = ? AND status = 'completed'
SELECT 
			COUNT(*) as total_purchases,
			COALESCE(SUM(amount), 0) as total_revenue,
			COALESCE(AVG(amount), 0) as avg_order_value,
			COUNT(DISTINCT user_id) as unique_buyers
		FROM purchases 
		WHERE event_id = ? AND status = 'completed'
	`

	var totalPurchases, uniqueBuyers int
	var totalRevenue, avgOrderValue float64

	err := es.db.QueryRow(query, eventID).Scan(&totalPurchases, &totalRevenue, &avgOrderValue, &uniqueBuyers)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase analytics: %w", err)
	}

	conversionRate := 0.0
	if totalPurchases > 0 {
		// Get total attendees for conversion rate
		attendeeQuery := `SELECT COUNT(DISTINCT user_id) FROM attendances WHERE event_id = ?`
		var totalAttendees int
		es.db.QueryRow(attendeeQuery, eventID).Scan(&totalAttendees)

		if totalAttendees > 0 {
			conversionRate = float64(uniqueBuyers) / float64(totalAttendees) * 100
		}
	}

	return map[string]interface{}{
		"totalPurchases": totalPurchases,
		"totalRevenue":   totalRevenue,
		"avgOrderValue":  avgOrderValue,
		"uniqueBuyers":   uniqueBuyers,
		"conversionRate": conversionRate,
	}, nil
}

func (es *EcommerceService) GetTopProducts(eventID string, limit int) ([]map[string]interface{}, error) {
	query := `
-- 		SELECT 
-- 			product_id,
-- 			COUNT(*) as purchase_count,
-- 			SUM(amount) as total_revenue
-- 		FROM purchase_products 
-- 		WHERE event_id = ? AND status = 'completed'
-- 		GROUP BY product_id
-- 		ORDER BY purchase_count DESC
-- 		LIMIT ?
        SELECT pp.product_name, SUM(pp.quantity) as total_sold
        FROM purchase_products pp
        JOIN purchases p ON pp.purchase_id = p.id  
        WHERE p.event_id = ? AND p.status = 'completed'
        GROUP BY pp.product_id
        ORDER BY total_sold DESC
	`

	rows, err := es.db.Query(query, eventID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top products: %w", err)
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var productID string
		var purchaseCount int
		var totalRevenue float64

		err := rows.Scan(&productID, &purchaseCount, &totalRevenue)
		if err != nil {
			continue
		}

		products = append(products, map[string]interface{}{
			"productId":     productID,
			"purchaseCount": purchaseCount,
			"totalRevenue":  totalRevenue,
		})
	}

	return products, nil
}
