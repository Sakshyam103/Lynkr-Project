/**
 * Feedback Service
 * Handles feedback collection and processing
 */

package services

import (
	"database/sql"
	"fmt"
	"time"
)

type FeedbackService struct {
	db *sql.DB
}

type Poll struct {
	ID       string       `json:"id"`
	Question string       `json:"question"`
	Options  []PollOption `json:"options"`
	EventID  string       `json:"eventId"`
}

type PollOption struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

type SliderFeedback struct {
	ID       string  `json:"id"`
	UserID   string  `json:"userId"`
	SliderID string  `json:"sliderId"`
	Value    float64 `json:"value"`
	EventID  string  `json:"eventId"`
}

func NewFeedbackService(db *sql.DB) *FeedbackService {
	return &FeedbackService{db: db}
}

// CreatePoll creates a new poll
func (fs *FeedbackService) CreatePoll(question, eventID string, options []string) (*Poll, error) {
	pollID := fmt.Sprintf("poll_%d", time.Now().UnixNano())
	
	poll := &Poll{
		ID:       pollID,
		Question: question,
		EventID:  eventID,
		Options:  make([]PollOption, len(options)),
	}
	
	for i, option := range options {
		poll.Options[i] = PollOption{
			ID:    fmt.Sprintf("option_%d", i),
			Text:  option,
			Votes: 0,
		}
	}
	
	return poll, nil
}

// SubmitPollVote records a poll vote
func (fs *FeedbackService) SubmitPollVote(userID, pollID, optionID string) error {
	query := `
		INSERT INTO poll_votes (user_id, poll_id, option_id, created_at)
		VALUES (?, ?, ?, ?)
	`
	
	_, err := fs.db.Exec(query, userID, pollID, optionID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to submit poll vote: %w", err)
	}
	
	return nil
}

// SubmitSliderFeedback records slider feedback
func (fs *FeedbackService) SubmitSliderFeedback(userID, sliderID string, value float64, eventID string) error {
	query := `
		INSERT INTO slider_feedback (user_id, slider_id, value, event_id, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	
	_, err := fs.db.Exec(query, userID, sliderID, value, eventID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to submit slider feedback: %w", err)
	}
	
	return nil
}

// GetEventFeedbackSummary returns feedback summary for an event
func (fs *FeedbackService) GetEventFeedbackSummary(eventID string) (map[string]interface{}, error) {
	// Get poll results
	pollQuery := `
		SELECT pv.poll_id, pv.option_id, COUNT(*) as votes
		FROM poll_votes pv
		JOIN polls p ON pv.poll_id = p.id
		WHERE p.event_id = ?
		GROUP BY pv.poll_id, pv.option_id
	`
	
	pollRows, err := fs.db.Query(pollQuery, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get poll results: %w", err)
	}
	defer pollRows.Close()
	
	pollResults := make(map[string]map[string]int)
	for pollRows.Next() {
		var pollID, optionID string
		var votes int
		
		if err := pollRows.Scan(&pollID, &optionID, &votes); err != nil {
			continue
		}
		
		if pollResults[pollID] == nil {
			pollResults[pollID] = make(map[string]int)
		}
		pollResults[pollID][optionID] = votes
	}
	
	// Get slider averages
	sliderQuery := `
		SELECT slider_id, AVG(value) as avg_value, COUNT(*) as responses
		FROM slider_feedback
		WHERE event_id = ?
		GROUP BY slider_id
	`
	
	sliderRows, err := fs.db.Query(sliderQuery, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get slider results: %w", err)
	}
	defer sliderRows.Close()
	
	sliderResults := make(map[string]map[string]interface{})
	for sliderRows.Next() {
		var sliderID string
		var avgValue float64
		var responses int
		
		if err := sliderRows.Scan(&sliderID, &avgValue, &responses); err != nil {
			continue
		}
		
		sliderResults[sliderID] = map[string]interface{}{
			"average":   avgValue,
			"responses": responses,
		}
	}
	
	return map[string]interface{}{
		"polls":   pollResults,
		"sliders": sliderResults,
	}, nil
}