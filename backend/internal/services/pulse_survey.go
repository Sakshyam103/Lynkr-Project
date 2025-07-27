/**
 * Pulse Survey Service
 * Handles delayed pulse surveys for post-event engagement
 */

package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type PulseSurvey struct {
	ID        string     `json:"id"`
	EventID   string     `json:"eventId"`
	BrandID   string     `json:"brandId"`
	Type      string     `json:"type"`
	Questions []Question `json:"questions"`
	ExpiresAt time.Time  `json:"expiresAt"`
}

type Question struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Type    string   `json:"type"`
	Options []string `json:"options,omitempty"`
}

type SurveyResponse struct {
	SurveyID  string                 `json:"surveyId"`
	UserID    string                 `json:"userId"`
	Responses map[string]interface{} `json:"responses"`
}

type Survey struct {
	ID           string `json:"id"`
	EventID      string `json:"event_id"`
	BrandID      string `json:"brand_id"`
	Title        string `json:"title"`
	RewardPoints int    `json:"reward_points"`
}

type PulseSurveyService struct {
	db *sql.DB
}

func NewPulseSurveyService(db *sql.DB) *PulseSurveyService {
	return &PulseSurveyService{db: db}
}

func (pss *PulseSurveyService) ScheduleSurveys(eventID, brandID string) error {
	surveys := []struct {
		surveyType string
		delay      time.Duration
		questions  []Question
	}{
		{
			surveyType: "24h",
			delay:      24 * time.Hour,
			questions: []Question{
				{ID: "q1", Text: "How likely are you to recommend this brand?", Type: "scale"},
				{ID: "q2", Text: "What was your favorite part of the event?", Type: "text"},
			},
		},
		{
			surveyType: "72h",
			delay:      72 * time.Hour,
			questions: []Question{
				{ID: "q1", Text: "Have you visited the brand website since the event?", Type: "boolean"},
				{ID: "q2", Text: "Have you made any purchases from this brand?", Type: "boolean"},
			},
		},
		{
			surveyType: "7d",
			delay:      7 * 24 * time.Hour,
			questions: []Question{
				{ID: "q1", Text: "How has your perception of the brand changed?", Type: "multiple_choice", Options: []string{"Much better", "Better", "Same", "Worse", "Much worse"}},
				{ID: "q2", Text: "Would you attend another event by this brand?", Type: "boolean"},
			},
		},
	}

	for _, survey := range surveys {
		err := pss.createSurvey(eventID, brandID, survey.surveyType, survey.questions, survey.delay)
		if err != nil {
			return fmt.Errorf("failed to schedule %s survey: %w", survey.surveyType, err)
		}
	}

	return nil
}

func (pss *PulseSurveyService) createSurvey(eventID, brandID, surveyType string, questions []Question, delay time.Duration) error {
	surveyID := fmt.Sprintf("survey_%s_%s_%d", surveyType, eventID, time.Now().UnixNano())
	questionsJSON, _ := json.Marshal(questions)
	expiresAt := time.Now().Add(delay + 24*time.Hour) // Survey available for 24h after delay

	query := `
		INSERT INTO pulse_surveys (id, event_id, brand_id, survey_type, questions, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := pss.db.Exec(query, surveyID, eventID, brandID, surveyType, string(questionsJSON), expiresAt, time.Now())
	return err
}

func (pss *PulseSurveyService) GetAvailableSurveys(userID string) ([]PulseSurvey, error) {
	query := `
		SELECT ps.id, ps.event_id, ps.brand_id, ps.questions, ps.expires_at
		FROM pulse_surveys ps
		JOIN attendances a ON ps.event_id = a.event_id
		WHERE a.user_id = ? 
		AND ps.expires_at > datetime('now')
		AND ps.id NOT IN (SELECT survey_id FROM pulse_survey_responses WHERE user_id = ?)
		AND datetime('now') >= datetime(ps.created_at, '+' || 
			CASE ps.survey_type 
				WHEN '24h' THEN '1 day'
				WHEN '72h' THEN '3 days' 
				WHEN '7d' THEN '7 days'
			END)
	`

	rows, err := pss.db.Query(query, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available surveys: %w", err)
	}
	defer rows.Close()

	var surveys []PulseSurvey
	for rows.Next() {
		var survey PulseSurvey
		var questionsJSON string

		err := rows.Scan(&survey.ID, &survey.EventID, &survey.BrandID, &questionsJSON, &survey.ExpiresAt)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(questionsJSON), &survey.Questions)
		surveys = append(surveys, survey)
	}

	return surveys, nil
}

func (pss *PulseSurveyService) SubmitResponse(userID, surveyID string, responses map[string]interface{}) error {
	responsesJSON, _ := json.Marshal(responses)

	query := `
		INSERT INTO pulse_survey_responses (survey_id, user_id, responses, completed_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := pss.db.Exec(query, surveyID, userID, string(responsesJSON), time.Now())
	if err != nil {
		return fmt.Errorf("failed to submit survey response: %w", err)
	}

	return nil
}

func (pss *PulseSurveyService) GetSurveyByID(surveyID string) (*Survey, error) {
	var survey Survey
	query := `SELECT id, title, event_id, brand_id, reward_points FROM pulse_surveys WHERE id = ?`
	err := pss.db.QueryRow(query, surveyID).Scan(&survey.ID, &survey.Title, &survey.EventID, &survey.BrandID, &survey.RewardPoints)
	if err != nil {
		return nil, fmt.Errorf("failed to get survey: %w", err)
	}
	return &survey, nil
}

func (pss *PulseSurveyService) GetSurveyAnalytics(eventID string) (map[string]interface{}, error) {
	query := `
		SELECT ps.survey_type, COUNT(psr.id) as responses, ps.questions
		FROM pulse_surveys ps
		LEFT JOIN pulse_survey_responses psr ON ps.id = psr.survey_id
		WHERE ps.event_id = ?
		GROUP BY ps.id, ps.survey_type
	`

	rows, err := pss.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get survey analytics: %w", err)
	}
	defer rows.Close()

	analytics := make(map[string]interface{})
	for rows.Next() {
		var surveyType, questionsJSON string
		var responses int

		err := rows.Scan(&surveyType, &responses, &questionsJSON)
		if err != nil {
			continue
		}

		analytics[surveyType] = map[string]interface{}{
			"responses": responses,
		}
	}

	return analytics, nil
}
