/**
 * Feedback Handlers
 * HTTP handlers for feedback collection and sentiment analysis
 */

package handlers

import (
	"net/http"

	"lynkr/internal/services"

	"github.com/gin-gonic/gin"
)

type FeedbackHandler struct {
	feedbackService  *services.FeedbackService
	sentimentService *services.SentimentService
}

func NewFeedbackHandler(feedbackService *services.FeedbackService, sentimentService *services.SentimentService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService:  feedbackService,
		sentimentService: sentimentService,
	}
}

// SubmitPollVote handles poll vote submission
func (fh *FeedbackHandler) SubmitPollVote(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		PollID   string `json:"pollId"`
		OptionID string `json:"optionId"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Vote recorded successfully",
	})
}

// SubmitSliderFeedback handles slider feedback submission
func (fh *FeedbackHandler) SubmitSliderFeedback(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		SliderID string  `json:"sliderId"`
		Value    float64 `json:"value"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Feedback recorded successfully",
	})
}

// SubmitQuickFeedback handles quick feedback submission
func (fh *FeedbackHandler) SubmitQuickFeedback(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		WidgetID string `json:"widgetId"`
		OptionID string `json:"optionId"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Quick feedback recorded successfully",
	})
}

// AnalyzeSentiment handles sentiment analysis requests
func (fh *FeedbackHandler) AnalyzeSentiment(c *gin.Context) {
	var request struct {
		ContentID string `json:"contentId"`
		Text      string `json:"text"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	analysis, err := fh.sentimentService.AnalyzeContent(request.ContentID, request.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze sentiment"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetEventSentiment returns sentiment summary for an event
func (fh *FeedbackHandler) GetEventSentiment(c *gin.Context) {
	eventID := c.Param("id")

	summary, err := fh.sentimentService.GetEventSentimentSummary(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get event sentiment"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetUserBadges returns user's gamification badges
func (fh *FeedbackHandler) GetUserBadges(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Simulate badge data
	badges := []map[string]interface{}{
		{
			"id":          "badge_1",
			"title":       "Feedback Champion",
			"description": "Provided feedback on 10 products",
			"icon":        "🏆",
			"rarity":      "rare",
			"points":      100,
			"isNew":       false,
		},
		{
			"id":          "badge_2",
			"title":       "Event Explorer",
			"description": "Attended 5 different events",
			"icon":        "🗺️",
			"rarity":      "common",
			"points":      50,
			"isNew":       true,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"badges": badges,
	})
}
