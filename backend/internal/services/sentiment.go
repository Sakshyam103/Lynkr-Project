/**
 * Sentiment Analysis Service
 * NLP processing for comment and post sentiment analysis
 */

package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type SentimentResult struct {
	Score     float64 `json:"score"`     // -1 to 1 (negative to positive)
	Magnitude float64 `json:"magnitude"` // 0 to 1 (intensity)
	Label     string  `json:"label"`     // positive, negative, neutral
}

type SentimentAnalysis struct {
	ID        string          `json:"id"`
	ContentID string          `json:"contentId"`
	Text      string          `json:"text"`
	Result    SentimentResult `json:"result"`
	CreatedAt time.Time       `json:"createdAt"`
}

type SentimentService struct {
	db *sql.DB
}

func NewSentimentService(db *sql.DB) *SentimentService {
	return &SentimentService{db: db}
}

// AnalyzeText performs sentiment analysis on text
func (ss *SentimentService) AnalyzeText(text string) (*SentimentResult, error) {
	// Simplified sentiment analysis - in production would use ML model
	cleanText := ss.cleanText(text)
	
	positiveWords := []string{"good", "great", "amazing", "awesome", "love", "excellent", "fantastic", "wonderful", "perfect", "best"}
	negativeWords := []string{"bad", "terrible", "awful", "hate", "worst", "horrible", "disappointing", "poor", "useless", "boring"}
	
	words := strings.Fields(strings.ToLower(cleanText))
	positiveCount := 0
	negativeCount := 0
	
	for _, word := range words {
		for _, pos := range positiveWords {
			if strings.Contains(word, pos) {
				positiveCount++
			}
		}
		for _, neg := range negativeWords {
			if strings.Contains(word, neg) {
				negativeCount++
			}
		}
	}
	
	totalSentimentWords := positiveCount + negativeCount
	magnitude := float64(totalSentimentWords) / float64(len(words))
	if magnitude > 1 {
		magnitude = 1
	}
	
	var score float64
	var label string
	
	if totalSentimentWords == 0 {
		score = 0
		label = "neutral"
	} else {
		score = (float64(positiveCount) - float64(negativeCount)) / float64(totalSentimentWords)
		if score > 0.1 {
			label = "positive"
		} else if score < -0.1 {
			label = "negative"
		} else {
			label = "neutral"
		}
	}
	
	return &SentimentResult{
		Score:     score,
		Magnitude: magnitude,
		Label:     label,
	}, nil
}

// AnalyzeContent analyzes sentiment for content and stores result
func (ss *SentimentService) AnalyzeContent(contentID, text string) (*SentimentAnalysis, error) {
	result, err := ss.AnalyzeText(text)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze sentiment: %w", err)
	}
	
	analysisID := fmt.Sprintf("sentiment_%d", time.Now().UnixNano())
	resultJSON, _ := json.Marshal(result)
	
	query := `
		INSERT INTO sentiment_analysis (id, content_id, text, result, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	now := time.Now()
	_, err = ss.db.Exec(query, analysisID, contentID, text, string(resultJSON), now)
	if err != nil {
		return nil, fmt.Errorf("failed to store sentiment analysis: %w", err)
	}
	
	return &SentimentAnalysis{
		ID:        analysisID,
		ContentID: contentID,
		Text:      text,
		Result:    *result,
		CreatedAt: now,
	}, nil
}

// GetContentSentiment retrieves sentiment analysis for content
func (ss *SentimentService) GetContentSentiment(contentID string) (*SentimentAnalysis, error) {
	query := `
		SELECT id, content_id, text, result, created_at
		FROM sentiment_analysis WHERE content_id = ?
		ORDER BY created_at DESC LIMIT 1
	`
	
	var analysis SentimentAnalysis
	var resultJSON string
	
	err := ss.db.QueryRow(query, contentID).Scan(
		&analysis.ID, &analysis.ContentID, &analysis.Text,
		&resultJSON, &analysis.CreatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get sentiment analysis: %w", err)
	}
	
	json.Unmarshal([]byte(resultJSON), &analysis.Result)
	
	return &analysis, nil
}

// GetEventSentimentSummary returns sentiment summary for an event
func (ss *SentimentService) GetEventSentimentSummary(eventID string) (map[string]interface{}, error) {
	query := `
		SELECT sa.result
		FROM sentiment_analysis sa
		JOIN content c ON sa.content_id = c.id
		WHERE c.event_id = ?
	`
	
	rows, err := ss.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event sentiment: %w", err)
	}
	defer rows.Close()
	
	var results []SentimentResult
	for rows.Next() {
		var resultJSON string
		if err := rows.Scan(&resultJSON); err != nil {
			continue
		}
		
		var result SentimentResult
		if json.Unmarshal([]byte(resultJSON), &result) == nil {
			results = append(results, result)
		}
	}
	
	if len(results) == 0 {
		return map[string]interface{}{
			"totalAnalyses": 0,
			"averageScore": 0,
			"distribution": map[string]int{
				"positive": 0,
				"negative": 0,
				"neutral":  0,
			},
		}, nil
	}
	
	// Calculate summary statistics
	var totalScore float64
	distribution := map[string]int{
		"positive": 0,
		"negative": 0,
		"neutral":  0,
	}
	
	for _, result := range results {
		totalScore += result.Score
		distribution[result.Label]++
	}
	
	return map[string]interface{}{
		"totalAnalyses": len(results),
		"averageScore":  totalScore / float64(len(results)),
		"distribution":  distribution,
	}, nil
}

// cleanText removes special characters and normalizes text
func (ss *SentimentService) cleanText(text string) string {
	// Remove URLs
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	text = urlRegex.ReplaceAllString(text, "")
	
	// Remove mentions and hashtags
	mentionRegex := regexp.MustCompile(`@[^\s]+`)
	text = mentionRegex.ReplaceAllString(text, "")
	
	hashtagRegex := regexp.MustCompile(`#[^\s]+`)
	text = hashtagRegex.ReplaceAllString(text, "")
	
	// Remove extra whitespace
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	return strings.TrimSpace(text)
}