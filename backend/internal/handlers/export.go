/**
 * Export Handlers
 * HTTP handlers for data export and CRM integration management
 */

package handlers

import (
	"net/http"

	"lynkr/internal/services"

	"github.com/gin-gonic/gin"
)

type ExportHandler struct {
	exportService         *services.ExportService
	crmIntegrationService *services.CRMIntegrationService
}

func NewExportHandler(exportService *services.ExportService, crmIntegrationService *services.CRMIntegrationService) *ExportHandler {
	return &ExportHandler{
		exportService:         exportService,
		crmIntegrationService: crmIntegrationService,
	}
}

func (eh *ExportHandler) CreateExportRequest(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	var request struct {
		EventID  string `json:"eventId"`
		DataType string `json:"dataType"`
		Format   string `json:"format"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	exportReq, err := eh.exportService.CreateExportRequest(brandID, request.EventID, request.DataType, request.Format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create export request"})
		return
	}

	c.JSON(http.StatusOK, exportReq)
}

func (eh *ExportHandler) GetExportStatus(c *gin.Context) {
	requestID := c.Param("requestId")

	exportReq, err := eh.exportService.GetExportStatus(requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Export request not found"})
		return
	}

	c.JSON(http.StatusOK, exportReq)
}

func (eh *ExportHandler) CreateCRMIntegration(c *gin.Context) {
	brandID := c.GetString("brandID")
	if brandID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand ID required"})
		return
	}

	var request struct {
		CRMType      string `json:"crmType"`
		APIKey       string `json:"apiKey"`
		APISecret    string `json:"apiSecret"`
		WebhookURL   string `json:"webhookUrl"`
		SyncInterval int    `json:"syncInterval"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	integration, err := eh.crmIntegrationService.CreateIntegration(
		brandID, request.CRMType, request.APIKey, request.APISecret,
		request.WebhookURL, request.SyncInterval,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CRM integration"})
		return
	}

	// Schedule automatic sync
	eh.crmIntegrationService.ScheduleSync(integration.ID)

	c.JSON(http.StatusOK, integration)
}

func (eh *ExportHandler) SyncEventData(c *gin.Context) {
	integrationID := c.Param("integrationId")
	eventID := c.Param("eventId")

	err := eh.crmIntegrationService.SyncEventData(integrationID, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync event data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}

func (eh *ExportHandler) GetExportFormats(c *gin.Context) {
	formats := gin.H{
		"formats": []map[string]string{
			{"id": "csv", "name": "CSV", "description": "Comma-separated values"},
			{"id": "json", "name": "JSON", "description": "JavaScript Object Notation"},
		},
		"dataTypes": []map[string]string{
			{"id": "attendance", "name": "Attendance Data", "description": "Event attendance records"},
			{"id": "content", "name": "Content Data", "description": "User-generated content"},
			{"id": "analytics", "name": "Analytics Data", "description": "Event analytics and metrics"},
			{"id": "feedback", "name": "Feedback Data", "description": "Polls and survey responses"},
		},
	}

	c.JSON(http.StatusOK, formats)
}

func (eh *ExportHandler) GetCRMTypes(c *gin.Context) {
	crmTypes := gin.H{
		"crmTypes": []map[string]string{
			{"id": "salesforce", "name": "Salesforce", "description": "Salesforce CRM integration"},
			{"id": "hubspot", "name": "HubSpot", "description": "HubSpot CRM integration"},
			{"id": "mailchimp", "name": "Mailchimp", "description": "Mailchimp email marketing integration"},
		},
	}

	c.JSON(http.StatusOK, crmTypes)
}
