/**
 * Security Handlers
 * HTTP handlers for security auditing and privacy management
 */

package handlers

import (
	"strconv"

	"lynkr/internal/security"

	"github.com/gin-gonic/gin"
)

type SecurityHandler struct {
	securityAudit   *security.SecurityAudit
	privacyEnhancer *security.PrivacyEnhancer
}

func NewSecurityHandler(securityAudit *security.SecurityAudit, privacyEnhancer *security.PrivacyEnhancer) *SecurityHandler {
	return &SecurityHandler{
		securityAudit:   securityAudit,
		privacyEnhancer: privacyEnhancer,
	}
}

func (sh *SecurityHandler) RunSecurityScan(c *gin.Context) {
	vulnerabilities, err := sh.securityAudit.RunSecurityScan()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to run security scan"})
		return
	}
	c.JSON(200, gin.H{
		"vulnerabilities": vulnerabilities,
		"scan_completed":  true,
	})
}

func (sh *SecurityHandler) GetSecurityEvents(c *gin.Context) {
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	events, err := sh.securityAudit.GetSecurityEvents(limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get security events"})
		return
	}

	c.JSON(200, gin.H{
		"events": events,
	})
}

func (sh *SecurityHandler) UpdatePrivacySettings(c *gin.Context) {
	if err := sh.privacyEnhancer.UpdateConsentFlow(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update privacy settings"})
		return
	}

	if err := sh.privacyEnhancer.ImplementDataRetention(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to implement data retention"})
		return
	}

	c.JSON(200, gin.H{"status": "privacy_settings_updated"})
}

func (sh *SecurityHandler) RequestDataDeletion(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := sh.privacyEnhancer.ProcessDataDeletionRequest(userID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to process deletion request"})
		return
	}

	c.JSON(200, gin.H{
		"status":  "deletion_scheduled",
		"message": "Your data will be deleted in 30 days",
	})
}

func (sh *SecurityHandler) ExportUserData(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	data, err := sh.privacyEnhancer.GetUserDataExport(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to export user data"})
		return
	}

	c.JSON(200, data)
}

func (sh *SecurityHandler) AnonymizeUserData(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	if err := sh.privacyEnhancer.AnonymizeUserData(userID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to anonymize user data"})
		return
	}

	c.JSON(200, gin.H{"status": "data_anonymized"})
}

func (sh *SecurityHandler) ValidateInput(c *gin.Context) {
	var request struct {
		Input string `json:"input"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	isValid := sh.securityAudit.ValidateInput(request.Input)

	c.JSON(200, gin.H{
		"valid": isValid,
		"input": request.Input,
	})
}
