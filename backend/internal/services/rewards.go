/**
 * Rewards Service
 * Handles incentive system for user engagement and quality content
 */

package services

import (
	"database/sql"
	"fmt"
	"time"
)

type Reward struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Type        string    `json:"type"`
	Points      int       `json:"points"`
	Description string    `json:"description"`
	EventID     string    `json:"eventId"`
	ContentID   string    `json:"contentId"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type UserRewards struct {
	UserID      string `json:"userId"`
	TotalPoints int    `json:"totalPoints"`
	Level       int    `json:"level"`
	Badges      []string `json:"badges"`
}

type RewardsService struct {
	db *sql.DB
}

func NewRewardsService(db *sql.DB) *RewardsService {
	return &RewardsService{db: db}
}

func (rs *RewardsService) AwardReward(userID, rewardType, description, eventID, contentID string, points int) (*Reward, error) {
	rewardID := fmt.Sprintf("reward_%d", time.Now().UnixNano())
	
	query := `
		INSERT INTO rewards (id, user_id, type, points, description, event_id, content_id, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	_, err := rs.db.Exec(query, rewardID, userID, rewardType, points, description, eventID, contentID, "pending", now)
	if err != nil {
		return nil, fmt.Errorf("failed to award reward: %w", err)
	}
	
	// Update user points
	rs.updateUserPoints(userID, points)
	
	return &Reward{
		ID:          rewardID,
		UserID:      userID,
		Type:        rewardType,
		Points:      points,
		Description: description,
		EventID:     eventID,
		ContentID:   contentID,
		Status:      "pending",
		CreatedAt:   now,
	}, nil
}

func (rs *RewardsService) GetUserRewards(userID string) (*UserRewards, error) {
	query := `
		SELECT COALESCE(SUM(points), 0) as total_points
		FROM rewards 
		WHERE user_id = ? AND status = 'delivered'
	`
	
	var totalPoints int
	err := rs.db.QueryRow(query, userID).Scan(&totalPoints)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rewards: %w", err)
	}
	
	level := rs.calculateLevel(totalPoints)
	badges := rs.getUserBadges(userID)
	
	return &UserRewards{
		UserID:      userID,
		TotalPoints: totalPoints,
		Level:       level,
		Badges:      badges,
	}, nil
}

func (rs *RewardsService) EvaluateContentQuality(contentID string) (int, error) {
	// Simulate quality evaluation - in production would use ML models
	query := `
		SELECT view_count, share_count, 
		       (SELECT COUNT(*) FROM content_analytics WHERE content_id = ? AND action = 'like') as likes
		FROM content WHERE id = ?
	`
	
	var views, shares, likes int
	err := rs.db.QueryRow(query, contentID, contentID).Scan(&views, &shares, &likes)
	if err != nil {
		return 0, fmt.Errorf("failed to evaluate content quality: %w", err)
	}
	
	// Quality score based on engagement
	qualityScore := 0
	if views > 100 {
		qualityScore += 10
	}
	if shares > 5 {
		qualityScore += 20
	}
	if likes > 10 {
		qualityScore += 15
	}
	
	return qualityScore, nil
}

func (rs *RewardsService) ProcessQualityRewards() error {
	// Get unprocessed content from last 24 hours
	query := `
		SELECT id, user_id, event_id 
		FROM content 
		WHERE created_at > datetime('now', '-1 day') 
		AND id NOT IN (SELECT content_id FROM rewards WHERE content_id IS NOT NULL)
	`
	
	rows, err := rs.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to get content for quality evaluation: %w", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var contentID, userID, eventID string
		if err := rows.Scan(&contentID, &userID, &eventID); err != nil {
			continue
		}
		
		qualityScore, err := rs.EvaluateContentQuality(contentID)
		if err != nil {
			continue
		}
		
		if qualityScore > 20 {
			rs.AwardReward(userID, "quality_content", "High-quality content reward", eventID, contentID, qualityScore)
		}
	}
	
	return nil
}

func (rs *RewardsService) updateUserPoints(userID string, points int) error {
	query := `
		INSERT INTO user_points (user_id, total_points, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
		total_points = total_points + ?, updated_at = ?
	`
	
	now := time.Now()
	_, err := rs.db.Exec(query, userID, points, now, points, now)
	return err
}

func (rs *RewardsService) calculateLevel(points int) int {
	if points < 100 {
		return 1
	} else if points < 500 {
		return 2
	} else if points < 1000 {
		return 3
	} else if points < 2500 {
		return 4
	}
	return 5
}

func (rs *RewardsService) getUserBadges(userID string) []string {
	query := `SELECT badge_id FROM user_badges WHERE user_id = ?`
	
	rows, err := rs.db.Query(query, userID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()
	
	var badges []string
	for rows.Next() {
		var badgeID string
		if err := rows.Scan(&badgeID); err == nil {
			badges = append(badges, badgeID)
		}
	}
	
	return badges
}