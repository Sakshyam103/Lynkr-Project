/**
 * Rewards Handlers
 * HTTP handlers for rewards and pulse survey management
 */

package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"lynkr/internal/services"

	"github.com/gin-gonic/gin"
)

type RewardsHandler struct {
	rewardsService     *services.RewardsService
	pulseSurveyService *services.PulseSurveyService
}

func NewRewardsHandler(rewardsService *services.RewardsService, pulseSurveyService *services.PulseSurveyService) *RewardsHandler {
	return &RewardsHandler{
		rewardsService:     rewardsService,
		pulseSurveyService: pulseSurveyService,
	}
}

func (rh *RewardsHandler) GetUserRewards(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rewards, err := rh.rewardsService.GetUserRewards(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user rewards"})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (rh *RewardsHandler) AwardReward(c *gin.Context) {
	var request struct {
		UserID      string `json:"userId"`
		Type        string `json:"type"`
		Description string `json:"description"`
		EventID     string `json:"eventId"`
		ContentID   string `json:"contentId"`
		Points      int    `json:"points"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	reward, err := rh.rewardsService.AwardReward(
		request.UserID, request.Type, request.Description,
		request.EventID, request.ContentID, request.Points,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to award reward"})
		return
	}

	c.JSON(http.StatusOK, reward)
}

func (rh *RewardsHandler) ProcessQualityRewards(c *gin.Context) {
	err := rh.rewardsService.ProcessQualityRewards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process quality rewards"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

func (rh *RewardsHandler) GetAvailableSurveys(c *gin.Context) {
	userID := c.GetInt("userID")
	if userID == -1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDString := strconv.Itoa(userID)
	surveys, err := rh.pulseSurveyService.GetAvailableSurveys(userIDString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get available surveys"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"surveys": surveys})
}

func (rh *RewardsHandler) SubmitSurveyResponse(c *gin.Context) {
	userID := strconv.Itoa(c.GetInt("userID"))
	if userID == "-1" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		SurveyID  string                 `json:"surveyId"`
		Responses map[string]interface{} `json:"responses"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := rh.pulseSurveyService.SubmitResponse(userID, request.SurveyID, request.Responses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit survey response"})
		return
	}

	// Get survey details to know how many points to award
	survey, err := rh.pulseSurveyService.GetSurveyByID(request.SurveyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get survey details"})
		return
	}

	//// Award points to the user
	//err = rh.rewardsService.AwardPoints(userID, survey.RewardPoints, "survey_completion")
	_, err = rh.rewardsService.AwardReward(
		userID, "survey_completion", fmt.Sprintf("Completed survey: %s", survey.Title),
		survey.EventID, "", survey.RewardPoints,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to award points"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":        "submitted",
		"pointsAwarded": survey.RewardPoints,
	})
}

func (rh *RewardsHandler) ScheduleSurveys(c *gin.Context) {
	eventID := c.Param("id")
	brandID := c.GetHeader("X-Brand-ID")

	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	err := rh.pulseSurveyService.ScheduleSurveys(eventID, brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule surveys"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "scheduled"})
}

func (rh *RewardsHandler) GetSurveyAnalytics(c *gin.Context) {
	eventID := c.Param("id")

	analytics, err := rh.pulseSurveyService.GetSurveyAnalytics(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get survey analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
