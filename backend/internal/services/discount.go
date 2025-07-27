/**
 * Discount Code Service
 * Handles unique discount code generation and redemption tracking
 */

package services

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type DiscountCode struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	EventID     string    `json:"eventId"`
	BrandID     string    `json:"brandId"`
	DiscountPct float64   `json:"discountPct"`
	MaxUses     int       `json:"maxUses"`
	UsedCount   int       `json:"usedCount"`
	ExpiresAt   time.Time `json:"expiresAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CodeRedemption struct {
	ID        string    `json:"id"`
	CodeID    string    `json:"codeId"`
	UserID    string    `json:"userId"`
	OrderID   string    `json:"orderId"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}

type DiscountService struct {
	db *sql.DB
}

func NewDiscountService(db *sql.DB) *DiscountService {
	return &DiscountService{db: db}
}

func (ds *DiscountService) GenerateCode(eventID, brandID string, discountPct float64, maxUses int, expiresAt time.Time) (*DiscountCode, error) {
	code := ds.generateUniqueCode()
	codeID := fmt.Sprintf("discount_%d", time.Now().UnixNano())

	query := `
		INSERT INTO discount_codes (id, code, event_id, brand_id, discount_pct, max_uses, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err := ds.db.Exec(query, codeID, code, eventID, brandID, discountPct, maxUses, expiresAt, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create discount code: %w", err)
	}

	return &DiscountCode{
		ID:          codeID,
		Code:        code,
		EventID:     eventID,
		BrandID:     brandID,
		DiscountPct: discountPct,
		MaxUses:     maxUses,
		UsedCount:   0,
		ExpiresAt:   expiresAt,
		CreatedAt:   now,
	}, nil
}

func (ds *DiscountService) ValidateCode(code string) (*DiscountCode, error) {
	query := `
		SELECT id, code, event_id, brand_id, discount_pct, max_uses, used_count, expires_at, created_at
		FROM discount_codes 
		WHERE code = ? AND expires_at > ? AND (max_uses = 0 OR used_count < max_uses)
	`

	var dc DiscountCode
	err := ds.db.QueryRow(query, code, time.Now()).Scan(
		&dc.ID, &dc.Code, &dc.EventID, &dc.BrandID,
		&dc.DiscountPct, &dc.MaxUses, &dc.UsedCount,
		&dc.ExpiresAt, &dc.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("invalid or expired discount code")
	}

	return &dc, nil
}

func (ds *DiscountService) RedeemCode(codeID, userID, orderID string, amount float64) (*CodeRedemption, error) {
	redemptionID := fmt.Sprintf("redemption_%d", time.Now().UnixNano())

	tx, err := ds.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert redemption
	query1 := `
		INSERT INTO code_redemptions (id, code_id, user_id, order_id, amount, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err = tx.Exec(query1, redemptionID, codeID, userID, orderID, amount, now)
	if err != nil {
		return nil, fmt.Errorf("failed to record redemption: %w", err)
	}

	// Update used count
	query2 := `UPDATE discount_codes SET used_count = used_count + 1 WHERE id = ?`
	_, err = tx.Exec(query2, codeID)
	if err != nil {
		return nil, fmt.Errorf("failed to update code usage: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &CodeRedemption{
		ID:        redemptionID,
		CodeID:    codeID,
		UserID:    userID,
		OrderID:   orderID,
		Amount:    amount,
		CreatedAt: now,
	}, nil
}

func (ds *DiscountService) GetCodeAnalytics(eventID string) (map[string]interface{}, error) {
	query := `
-- 		SELECT 
-- 			dc.code,
-- 			dc.discount_pct,
-- 			dc.used_count,
-- 			dc.max_uses,
-- 			COALESCE(SUM(cr.amount), 0) as total_revenue
-- 		FROM discount_codes dc
-- 		LEFT JOIN code_redemptions cr ON dc.id = cr.code_id
-- 		WHERE dc.event_id = ?
-- 		GROUP BY dc.id
-- 		ORDER BY dc.used_count DESC
SELECT 
    dc.code,
    dc.discount_pct,
    dc.used_count,
    dc.max_uses,
    COALESCE(COUNT(dcu.id), 0) as actual_usage
FROM discount_codes dc
LEFT JOIN discount_code_usage dcu ON dc.id = dcu.discount_code_id
WHERE dc.event_id = ?
GROUP BY dc.id
ORDER BY dc.used_count DESC
	`

	rows, err := ds.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get code analytics: %w", err)
	}
	defer rows.Close()

	var codes []map[string]interface{}
	totalRedemptions := 0
	totalRevenue := 0.0

	for rows.Next() {
		var code string
		var discountPct float64
		var usedCount, maxUses int
		var revenue float64

		err := rows.Scan(&code, &discountPct, &usedCount, &maxUses, &revenue)
		if err != nil {
			continue
		}

		codes = append(codes, map[string]interface{}{
			"code":        code,
			"discountPct": discountPct,
			"usedCount":   usedCount,
			"maxUses":     maxUses,
			"revenue":     revenue,
		})

		totalRedemptions += usedCount
		totalRevenue += revenue
	}

	return map[string]interface{}{
		"codes":            codes,
		"totalRedemptions": totalRedemptions,
		"totalRevenue":     totalRevenue,
	}, nil
}

func (ds *DiscountService) generateUniqueCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	b := make([]byte, length)
	rand.Read(b)

	var result strings.Builder
	for _, v := range b {
		result.WriteByte(charset[v%byte(len(charset))])
	}

	return result.String()
}
