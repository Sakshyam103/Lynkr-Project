/**
 * AI Tagging Service
 * Handles AI-powered product detection and automated content tagging
 */

package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ProductDetection struct {
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	BrandID     string  `json:"brandId"`
	Confidence  float64 `json:"confidence"`
	BoundingBox struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"boundingBox"`
}

type AITaggingResult struct {
	ContentID   string             `json:"contentId"`
	Products    []ProductDetection `json:"products"`
	Tags        []string           `json:"tags"`
	ProcessedAt time.Time          `json:"processedAt"`
}

type AITaggingService struct {
	db *sql.DB
}

func NewAITaggingService(db *sql.DB) *AITaggingService {
	return &AITaggingService{db: db}
}

func (ats *AITaggingService) ProcessContent(contentID, mediaURL string) (*AITaggingResult, error) {
	// Simulate AI processing - in production would call ML service
	products := ats.detectProducts(mediaURL)
	tags := ats.generateTags(products)
	
	result := &AITaggingResult{
		ContentID:   contentID,
		Products:    products,
		Tags:        tags,
		ProcessedAt: time.Now(),
	}
	
	// Store results
	err := ats.storeResults(result)
	if err != nil {
		return nil, fmt.Errorf("failed to store AI tagging results: %w", err)
	}
	
	return result, nil
}

func (ats *AITaggingService) detectProducts(mediaURL string) []ProductDetection {
	// Simulate product detection
	return []ProductDetection{
		{
			ProductID:   "prod_1",
			ProductName: "Sample Product",
			BrandID:     "brand_1",
			Confidence:  0.85,
			BoundingBox: struct {
				X      int `json:"x"`
				Y      int `json:"y"`
				Width  int `json:"width"`
				Height int `json:"height"`
			}{X: 100, Y: 150, Width: 200, Height: 250},
		},
	}
}

func (ats *AITaggingService) generateTags(products []ProductDetection) []string {
	tags := []string{}
	for _, product := range products {
		if product.Confidence > 0.7 {
			tags = append(tags, strings.ToLower(product.ProductName))
			tags = append(tags, "brand_"+product.BrandID)
		}
	}
	return tags
}

func (ats *AITaggingService) storeResults(result *AITaggingResult) error {
	productsJSON, _ := json.Marshal(result.Products)
	tagsJSON, _ := json.Marshal(result.Tags)
	
	query := `
		INSERT INTO ai_tagging_results (content_id, products, tags, processed_at)
		VALUES (?, ?, ?, ?)
	`
	
	_, err := ats.db.Exec(query, result.ContentID, string(productsJSON), string(tagsJSON), result.ProcessedAt)
	return err
}

func (ats *AITaggingService) GetProductAnalytics(brandID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			atr.products,
			COUNT(*) as detection_count
		FROM ai_tagging_results atr
		JOIN content c ON atr.content_id = c.id
		WHERE JSON_EXTRACT(atr.products, '$[0].brandId') = ?
		GROUP BY atr.products
		ORDER BY detection_count DESC
		LIMIT 10
	`
	
	rows, err := ats.db.Query(query, brandID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product analytics: %w", err)
	}
	defer rows.Close()
	
	var analytics []map[string]interface{}
	totalDetections := 0
	
	for rows.Next() {
		var productsJSON string
		var count int
		
		err := rows.Scan(&productsJSON, &count)
		if err != nil {
			continue
		}
		
		var products []ProductDetection
		json.Unmarshal([]byte(productsJSON), &products)
		
		if len(products) > 0 {
			analytics = append(analytics, map[string]interface{}{
				"productName":     products[0].ProductName,
				"detectionCount":  count,
				"avgConfidence":   products[0].Confidence,
			})
			totalDetections += count
		}
	}
	
	return map[string]interface{}{
		"products":         analytics,
		"totalDetections":  totalDetections,
	}, nil
}