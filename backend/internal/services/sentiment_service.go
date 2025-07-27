package services

//
//import (
//	"database/sql"
//	"encoding/json"
//	"fmt"
//	"strings"
//	"time"
//)
//
//type SentimentService struct {
//	db *sql.DB
//}
//
//type SentimentResult struct {
//	Score     float64 `json:"score"`     // -1 to 1 (-1 = negative, 0 = neutral, 1 = positive)
//	Sentiment string  `json:"sentiment"` // "positive", "negative", "neutral"
//	Confidence float64 `json:"confidence"` // 0 to 1
//}
//
//func NewSentimentService(db *sql.DB) *SentimentService {
//	return &SentimentService{db: db}
//}
//
//func (ss *SentimentService) AnalyzeSentiment(text string) (*SentimentResult, error) {
//	// Simple sentiment analysis (in production, use ML service like AWS Comprehend)
//	result := ss.simpleSentimentAnalysis(text)
//
//	// Store result in database
//	resultJSON, _ := json.Marshal(result)
//	query := `
//		INSERT INTO sentiment_analysis (content_id, result, created_at)
//		VALUES (?, ?, ?)
//	`
//
//	// For now, use text hash as content_id (in real app, pass actual content_id)
//	contentID := fmt.Sprintf("text_%d", len(text))
//	_, err := ss.db.Exec(query, contentID, string(resultJSON), time.Now())
//	if err != nil {
//		return nil, fmt.Errorf("failed to store sentiment analysis: %w", err)
//	}
//
//	return result, nil
//}
//
//func (ss *SentimentService) AnalyzeContentSentiment(contentID, text string) (*SentimentResult, error) {
//	result := ss.simpleSentimentAnalysis(text)
//
//	// Store result in database
//	resultJSON, _ := json.Marshal(result)
//	query := `
//		INSERT INTO sentiment_analysis (content_id, result, created_at)
//		VALUES (?, ?, ?)
//	`
//
//	_, err := ss.db.Exec(query, contentID, string(resultJSON), time.Now())
//	if err != nil {
//		return nil, fmt.Errorf("failed to store sentiment analysis: %w", err)
//	}
//
//	return result, nil
//}
//
//func (ss *SentimentService) GetEventSentiment(eventID string) (map[string]interface{}, error) {
//	query := `
//		SELECT sa.result
//		FROM sentiment_analysis sa
//		JOIN content c ON sa.content_id = c.id
//		WHERE c.event_id = ?
//		ORDER BY sa.created_at DESC
//	`
//
//	rows, err := ss.db.Query(query, eventID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get event sentiment: %w", err)
//	}
//	defer rows.Close()
//
//	var results []SentimentResult
//	for rows.Next() {
//		var resultJSON string
//		if err := rows.Scan(&resultJSON); err != nil {
//			continue
//		}
//
//		var result SentimentResult
//		if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
//			continue
//		}
//
//		results = append(results, result)
//	}
//
//	// Calculate aggregated sentiment
//	if len(results) == 0 {
//		return map[string]interface{}{
//			"overallSentiment": "neutral",
//			"averageScore": 0.0,
//			"totalAnalyzed": 0,
//			"breakdown": map[string]int{
//				"positive": 0,
//				"negative": 0,
//				"neutral": 0,
//			},
//		}, nil
//	}
//
//	var totalScore float64
//	breakdown := map[string]int{
//		"positive": 0,
//		"negative": 0,
//		"neutral": 0,
//	}
//
//	for _, result := range results {
//		totalScore += result.Score
//		breakdown[result.Sentiment]++
//	}
//
//	avgScore := totalScore / float64(len(results))
//	overallSentiment := "neutral"
//	if avgScore > 0.1 {
//		overallSentiment = "positive"
//	} else if avgScore < -0.1 {
//		overallSentiment = "negative"
//	}
//
//	return map[string]interface{}{
//		"overallSentiment": overallSentiment,
//		"averageScore": avgScore,
//		"totalAnalyzed": len(results),
//		"breakdown": breakdown,
//	}, nil
//}
//
//// Simple sentiment analysis (replace with ML service in production)
//func (ss *SentimentService) simpleSentimentAnalysis(text string) *SentimentResult {
//	text = strings.ToLower(text)
//
//	positiveWords := []string{"amazing", "great", "awesome", "love", "excellent", "fantastic", "wonderful", "good", "best", "perfect"}
//	negativeWords := []string{"bad", "terrible", "awful", "hate", "worst", "horrible", "disappointing", "poor", "boring", "sucks"}
//
//	positiveCount := 0
//	negativeCount := 0
//
//	for _, word := range positiveWords {
//		positiveCount += strings.Count(text, word)
//	}
//
//	for _, word := range negativeWords {
//		negativeCount += strings.Count(text, word)
//	}
//
//	totalWords := positiveCount + negativeCount
//	if totalWords == 0 {
//		return &SentimentResult{
//			Score: 0.0,
//			Sentiment: "neutral",
//			Confidence: 0.5,
//		}
//	}
//
//	score := float64(positiveCount-negativeCount) / float64(totalWords)
//	sentiment := "neutral"
//	confidence := 0.6
//
//	if score > 0.2 {
//		sentiment = "positive"
//		confidence = 0.8
//	} else if score < -0.2 {
//		sentiment = "negative"
//		confidence = 0.8
//	}
//
//	return &SentimentResult{
//		Score: score,
//		Sentiment: sentiment,
//		Confidence: confidence,
//	}
//}
