/**
 * E-commerce Platform Connector
 * Generic interface for e-commerce platform integrations
 */

package ecommerce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
	BrandID     string  `json:"brandId"`
}

type Purchase struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	ProductID string    `json:"productId"`
	Amount    float64   `json:"amount"`
	EventID   string    `json:"eventId"`
	Timestamp time.Time `json:"timestamp"`
}

type Connector interface {
	GetProducts(brandID string) ([]Product, error)
	CreatePurchase(userID, productID string, amount float64) (*Purchase, error)
	TrackConversion(eventID, userID string, purchase *Purchase) error
}

type ShopifyConnector struct {
	apiKey    string
	apiSecret string
	shopURL   string
}

func NewShopifyConnector(apiKey, apiSecret, shopURL string) *ShopifyConnector {
	return &ShopifyConnector{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		shopURL:   shopURL,
	}
}

func (sc *ShopifyConnector) GetProducts(brandID string) ([]Product, error) {
	// Simulate Shopify API call
	products := []Product{
		{
			ID:          "prod_1",
			Name:        "Sample Product",
			Price:       29.99,
			Description: "A great product from the event",
			ImageURL:    "https://example.com/product.jpg",
			BrandID:     brandID,
		},
	}
	return products, nil
}

func (sc *ShopifyConnector) CreatePurchase(userID, productID string, amount float64) (*Purchase, error) {
	purchase := &Purchase{
		ID:        fmt.Sprintf("purchase_%d", time.Now().UnixNano()),
		UserID:    userID,
		ProductID: productID,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	return purchase, nil
}

func (sc *ShopifyConnector) TrackConversion(eventID, userID string, purchase *Purchase) error {
	// Track conversion in analytics
	return nil
}

type WooCommerceConnector struct {
	apiKey    string
	apiSecret string
	storeURL  string
}

func NewWooCommerceConnector(apiKey, apiSecret, storeURL string) *WooCommerceConnector {
	return &WooCommerceConnector{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		storeURL:  storeURL,
	}
}

func (wc *WooCommerceConnector) GetProducts(brandID string) ([]Product, error) {
	url := fmt.Sprintf("%s/wp-json/wc/v3/products", wc.storeURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.SetBasicAuth(wc.apiKey, wc.apiSecret)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var wcProducts []struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Price       string  `json:"price"`
		Description string  `json:"description"`
		Images      []struct {
			Src string `json:"src"`
		} `json:"images"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&wcProducts); err != nil {
		return nil, err
	}
	
	products := make([]Product, len(wcProducts))
	for i, wp := range wcProducts {
		imageURL := ""
		if len(wp.Images) > 0 {
			imageURL = wp.Images[0].Src
		}
		
		products[i] = Product{
			ID:          fmt.Sprintf("%d", wp.ID),
			Name:        wp.Name,
			Price:       parsePrice(wp.Price),
			Description: wp.Description,
			ImageURL:    imageURL,
			BrandID:     brandID,
		}
	}
	
	return products, nil
}

func (wc *WooCommerceConnector) CreatePurchase(userID, productID string, amount float64) (*Purchase, error) {
	// Create order via WooCommerce API
	purchase := &Purchase{
		ID:        fmt.Sprintf("wc_purchase_%d", time.Now().UnixNano()),
		UserID:    userID,
		ProductID: productID,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	return purchase, nil
}

func (wc *WooCommerceConnector) TrackConversion(eventID, userID string, purchase *Purchase) error {
	return nil
}

func parsePrice(priceStr string) float64 {
	// Simple price parsing - in production would handle currency formatting
	var price float64
	fmt.Sscanf(priceStr, "%f", &price)
	return price
}