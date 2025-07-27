package main

import (
	"log"
	// "path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// ginSwagger "github.com/swaggo/gin-swagger"
	"lynkr/internal/handlers"

	// "github.com/lynkr/brand-activations/backend/internal/middleware"
	"lynkr/internal/middleware"
	"lynkr/internal/performance"
	"lynkr/internal/security"
	"lynkr/internal/services"
	"lynkr/internal/services/content"
	"lynkr/internal/services/event"
	"lynkr/internal/services/user"
	"lynkr/internal/ux"

	"lynkr/pkg/database"
	// "lynkr/pkg/geofencing"
	"lynkr/pkg/privacy"

	ginSwagger "github.com/swaggo/gin-swagger"

	_ "lynkr/docs"

	swaggerFiles "github.com/swaggo/files"
)

// @title Brand Activations API
// @version 1.0
// @description API for the Brand Activations feature
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Initialize database
	dbConfig := database.Config{
		DBPath:        "./data/brand_activations.db",
		MigrationsDir: "",
	}

	if err := database.Initialize(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize privacy components
	anonymizer := privacy.NewAnonymizer("brand-activations-salt")
	retentionManager := privacy.NewRetentionManager(database.DB, anonymizer)

	// Schedule data retention job to run daily
	retentionManager.ScheduleRetentionJob(24 * time.Hour)

	// Initialize services
	userService := user.NewUserService(database.DB)
	// userService := services.NewUserService(database.DB)
	eventService := event.NewEventService(database.DB)
	contentService := content.NewContentService(database.DB)
	// content1Service := services.NewContentService(database.DB)
	brandService := services.NewBrandService(database.DB)
	feedbackService := services.NewFeedbackService(database.DB)
	sentimentService := services.NewSentimentService(database.DB)
	analyticsService := services.NewAnalyticsService(database.DB)
	ecommerceService := services.NewEcommerceService(database.DB)
	discountService := services.NewDiscountService(database.DB)
	pixelService := services.NewPixelService(database.DB)
	aiTaggingService := services.NewAITaggingService(database.DB)
	conversionFunnelService := services.NewConversionFunnelService(database.DB)
	rewardsService := services.NewRewardsService(database.DB)
	pulseSurveyService := services.NewPulseSurveyService(database.DB)
	exportService := services.NewExportService(database.DB)
	crmIntegrationService := services.NewCRMIntegrationService(database.DB)
	// geofenceService := geofencing.NewGeofenceService()
	// geofenceService, _ := geofencing.ParseGeofenceData("./data/geofence.json")

	// Initialize performance components
	dbOptimizer := performance.NewDatabaseOptimizer(database.DB)
	cache := performance.NewCache()
	loadTester := performance.NewLoadTester()

	// Initialize security components
	securityAudit := security.NewSecurityAudit(database.DB)
	privacyEnhancer := security.NewPrivacyEnhancer(database.DB)

	// Initialize UX components
	usabilityTester := ux.NewUsabilityTester(database.DB)

	// Initialize handlers
	// userHandler := handlers.NewUserHandler(userService)
	// eventHandler := handlers.NewEventHandler(eventService, geofenceService)
	handler := handlers.NewHandler(userService, eventService, contentService)
	// contentHandler := handlers.NewContentHandler(content1Service)
	brandHandler := handlers.NewBrandHandler(brandService, "brand-activations-secret-key")
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService, sentimentService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	ecommerceHandler := handlers.NewEcommerceHandler(ecommerceService)
	discountHandler := handlers.NewDiscountHandler(discountService)
	pixelHandler := handlers.NewPixelHandler(pixelService)
	advancedAnalyticsHandler := handlers.NewAdvancedAnalyticsHandler(aiTaggingService, conversionFunnelService)
	rewardsHandler := handlers.NewRewardsHandler(rewardsService, pulseSurveyService)
	exportHandler := handlers.NewExportHandler(exportService, crmIntegrationService)
	performanceHandler := handlers.NewPerformanceHandler(dbOptimizer, cache, loadTester)
	securityHandler := handlers.NewSecurityHandler(securityAudit, privacyEnhancer)
	uxHandler := handlers.NewUXHandler(usabilityTester)

	// Setup database optimization
	dbOptimizer.SetupConnectionPool()
	dbOptimizer.OptimizeIndexes()

	// Setup security enhancements
	privacyEnhancer.UpdateConsentFlow()
	privacyEnhancer.ImplementDataRetention()

	// Initialize Gin
	r := gin.Default()

	// Apply global middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(gin.Recovery())
	r.Use(middleware.RateLimitMiddleware(100, time.Minute)) // 100 requests per minute

	// API routes
	api := r.Group("/api/v1")

	//public routes (no authentication required)
	api.POST("/users/register", handler.CreateUser)
	api.POST("/users/login", handler.Login)
	api.POST("/brands/login", brandHandler.Login)
	api.POST("/webhooks/:integrationId", ecommerceHandler.HandleWebhook)
	api.GET("/pixel/track", pixelHandler.TrackPixel)
	api.POST("/conversions/track", advancedAnalyticsHandler.TrackConversion)

	//user only routes
	userRoutes := r.Group("/user/v1")
	userRoutes.Use(middleware.AuthMiddleware())
	userRoutes.Use(middleware.UserOnlyMiddleware())
	userRoutes.PUT("/users/privacy", securityHandler.UpdatePrivacySettings)
	userRoutes.POST("/events/:id/checkin", handler.CheckInEvent)
	userRoutes.GET("/events/:id/tags", handler.GetEventTags)
	userRoutes.POST("/content", handler.CreateContent)
	userRoutes.POST("/content/:id/analytics", handler.TrackContentAnalytics) //not working
	userRoutes.POST("/sentiment/analyze", feedbackHandler.AnalyzeSentiment)
	userRoutes.PUT("/content/:id/permissions", handler.UpdateContentPermissions)
	userRoutes.GET("/users/rewards", rewardsHandler.GetUserRewards)
	userRoutes.GET("/surveys/available", rewardsHandler.GetAvailableSurveys)
	userRoutes.GET("/users/badges", feedbackHandler.GetUserBadges)
	userRoutes.POST("/feedback/polls/vote", feedbackHandler.SubmitPollVote)
	userRoutes.POST("/feedback/sliders", feedbackHandler.SubmitSliderFeedback)
	userRoutes.POST("/feedback/quick", feedbackHandler.SubmitQuickFeedback)
	userRoutes.POST("/ecommerce/purchases", ecommerceHandler.TrackPurchase)
	userRoutes.POST("/surveys/respond", rewardsHandler.SubmitSurveyResponse)
	userRoutes.GET("/discount/codes/:code/validate", discountHandler.ValidateCode)
	userRoutes.POST("/discount/redeem", discountHandler.RedeemCode)
	userRoutes.GET("/events", handler.ListEvents)
	userRoutes.GET("/events/:id", handler.GetEvent)
	userRoutes.GET("/events/:id/content", handler.GetEventContent)
	// api.POST("/analytics/track", analyticsHandler.TrackEvent)
	// userRoutes.POST("/ecommerce/purchases", ecommerceHandler.TrackPurchase)
	userRoutes.POST("/users/data/delete", securityHandler.RequestDataDeletion)
	userRoutes.GET("/users/data/export", securityHandler.ExportUserData)
	userRoutes.POST("/users/data/anonymize", securityHandler.AnonymizeUserData)
	userRoutes.POST("/security/privacy/update", securityHandler.UpdatePrivacySettings)
	userRoutes.POST("/analytics/track", analyticsHandler.TrackEvent) //need userID so to track user
	userRoutes.POST("/pixel/search", pixelHandler.TrackSearch)

	//brand only d
	brandRoutes := r.Group("/brand/v1")
	brandRoutes.Use(middleware.AuthMiddleware1())
	brandRoutes.Use(middleware.BrandOnlyMiddleware())

	// brandRoutes.POST("/events", handler.CreateEvent)  //instead of brand creating the event organization are creating the events
	brandRoutes.GET("/brands/dashboard", brandHandler.GetDashboardStats)
	brandRoutes.GET("/events", handler.ListEvents)
	brandRoutes.GET("/events/:id/content", handler.GetEventContent)
	brandRoutes.GET("/content/:id", handler.GetEventContent)
	brandRoutes.GET("/content/tags/search", handler.SearchTags) //not working
	brandRoutes.GET("/brands/campaigns", brandHandler.GetCampaigns)
	brandRoutes.POST("/brands/campaigns", brandHandler.CreateCampaign)
	brandRoutes.GET("/brands/content", brandHandler.GetBrandContent)
	brandRoutes.GET("/events/:id/sentiment", feedbackHandler.GetEventSentiment)
	brandRoutes.GET("/events/:id/analytics/engagement", analyticsHandler.GetEngagementMetrics)
	brandRoutes.GET("/events/:id/analytics/attendance", analyticsHandler.GetAttendanceAnalytics) //works but got nothing for now
	brandRoutes.GET("/events/:id/analytics/content", analyticsHandler.GetContentPerformance)
	brandRoutes.GET("/events/:id/analytics/realtime", analyticsHandler.GetRealtimeStats)
	brandRoutes.POST("/ecommerce/integrations", ecommerceHandler.CreateIntegration)
	brandRoutes.GET("/ecommerce/integrations", ecommerceHandler.GetIntegration)
	brandRoutes.GET("/events/:id/purchases/analytics", ecommerceHandler.GetPurchaseAnalytics)
	brandRoutes.GET("/events/:id/purchases/top-products", ecommerceHandler.GetTopProducts)
	brandRoutes.POST("/discount/generate", discountHandler.GenerateCode)
	brandRoutes.GET("/events/:id/discount/analytics", discountHandler.GetCodeAnalytics)
	brandRoutes.GET("/brands/discount/codes", discountHandler.GetBrandCodes)
	brandRoutes.GET("/events/:id/pixel/analytics", pixelHandler.GetPixelAnalytics)
	brandRoutes.GET("/pixel/generate", pixelHandler.GeneratePixelURL)
	brandRoutes.POST("/content/:id/ai-process", advancedAnalyticsHandler.ProcessContentAI)
	brandRoutes.GET("/brands/product-analytics", advancedAnalyticsHandler.GetProductAnalytics)
	brandRoutes.GET("/events/:id/conversion-funnel", advancedAnalyticsHandler.GetConversionFunnel)
	brandRoutes.GET("/events/:id/attribution-report", advancedAnalyticsHandler.GetAttributionReport)
	brandRoutes.POST("/rewards/award", rewardsHandler.AwardReward)
	brandRoutes.POST("/rewards/process-quality", rewardsHandler.ProcessQualityRewards)
	brandRoutes.POST("/events/:id/surveys/schedule", rewardsHandler.ScheduleSurveys)
	brandRoutes.GET("/events/:id/surveys/analytics", rewardsHandler.GetSurveyAnalytics)
	brandRoutes.POST("/export/create", exportHandler.CreateExportRequest)
	brandRoutes.GET("/export/:requestId/status", exportHandler.GetExportStatus)
	brandRoutes.GET("/export/formats", exportHandler.GetExportFormats)
	brandRoutes.POST("/crm/integrations", exportHandler.CreateCRMIntegration)
	brandRoutes.POST("/crm/:integrationId/sync/:eventId", exportHandler.SyncEventData)
	brandRoutes.GET("/crm/types", exportHandler.GetCRMTypes)
	brandRoutes.POST("/security/privacy/update", securityHandler.UpdatePrivacySettings)
	brandRoutes.POST("/events", handler.CreateEvent)

	adminRoutes := api.Group("/performance")
	adminRoutes.Use(middleware.AdminOnlyMiddleware()) // New admin role needed
	adminRoutes.POST("/optimize-db", performanceHandler.OptimizeDatabase)
	adminRoutes.GET("/db-metrics", performanceHandler.GetDatabaseMetrics)
	adminRoutes.GET("/query-stats", performanceHandler.GetQueryStats)
	adminRoutes.GET("/cache-stats", performanceHandler.GetCacheStats)
	adminRoutes.DELETE("/cache", performanceHandler.ClearCache)
	adminRoutes.POST("/load-test", performanceHandler.RunLoadTest)
	adminRoutes.POST("/maintenance", performanceHandler.RunMaintenance)
	adminRoutes.POST("/security/scan", securityHandler.RunSecurityScan)
	adminRoutes.GET("/security/events", securityHandler.GetSecurityEvents)
	adminRoutes.POST("/security/validate-input", securityHandler.ValidateInput)
	// adminRoutes.POST("/security/validate-output", securityHandler.ValidateOutput)
	adminRoutes.POST("/security/privacy/update", securityHandler.UpdatePrivacySettings)
	adminRoutes.GET("/ux/metrics", uxHandler.GetUsabilityMetrics)
	adminRoutes.GET("/ux/heatmap", uxHandler.GetHeatmapData)
	adminRoutes.GET("/ux/journey", uxHandler.GetUserJourney)
	adminRoutes.GET("/ux/pain-points", uxHandler.GetPainPoints)

	// User routes
	// api.GET("/users/profile", handler.)
	// api.PUT("/users/privacy", securityHandler.UpdatePrivacySettings)

	// // Event routes
	// // api.POST("/events", handler.CreateEvent)
	// // api.GET("/events/:id", handler.GetEvent)
	// api.POST("/events/:id/checkin", handler.CheckInEvent)
	// // api.GET("/events/:id/feed", handler.GetEventFeed)
	// // api.GET("/events/:eventId/tags", handler.GetEventTags)
	// // api.GET("/events/:eventId/content", handler.GetEventContent)
	// api.GET("/events/:id", handler.GetEvent)
	// api.GET("/events/:id/tags", handler.GetEventTags)
	// api.GET("/events/:id/content", handler.GetEventContent)

	// // Content routes
	// api.POST("/content", handler.CreateContent)
	// api.GET("/content/:id", handler.GetEventContent)
	// api.PUT("/content/:id/permissions", handler.UpdateContentPermissions)
	// api.POST("/content/:id/analytics", handler.TrackContentAnalytics)
	// api.GET("/content/:id/analytics", handler.GetContentAnalytics)
	// api.GET("/content/tags/search", handler.SearchTags)
	// api.POST("/content/tags/suggest", handler.GetSuggestedTags)

	// // Brand routes
	// api.GET("/brands/dashboard", brandHandler.GetDashboardStats)
	// api.GET("/brands/campaigns", brandHandler.GetCampaigns)
	// api.POST("/brands/campaigns", brandHandler.CreateCampaign)
	// api.GET("/brands/content", brandHandler.GetBrandContent)

	// // Feedback routes
	// api.POST("/feedback/polls/vote", feedbackHandler.SubmitPollVote)
	// api.POST("/feedback/sliders", feedbackHandler.SubmitSliderFeedback)
	// api.POST("/feedback/quick", feedbackHandler.SubmitQuickFeedback)
	// api.POST("/sentiment/analyze", feedbackHandler.AnalyzeSentiment)
	// api.GET("/events/:id/sentiment", feedbackHandler.GetEventSentiment)
	// api.GET("/users/badges", feedbackHandler.GetUserBadges)

	// // Analytics routes
	// api.GET("/events/:id/analytics/engagement", analyticsHandler.GetEngagementMetrics)
	// api.GET("/events/:id/analytics/attendance", analyticsHandler.GetAttendanceAnalytics)
	// api.GET("/events/:id/analytics/content", analyticsHandler.GetContentPerformance)
	// api.GET("/events/:id/analytics/realtime", analyticsHandler.GetRealtimeStats)
	// api.POST("/analytics/track", analyticsHandler.TrackEvent)

	// // E-commerce routes
	// api.POST("/ecommerce/integrations", ecommerceHandler.CreateIntegration)
	// api.GET("/ecommerce/integrations", ecommerceHandler.GetIntegration)
	// api.POST("/ecommerce/purchases", ecommerceHandler.TrackPurchase)
	// api.GET("/events/:id/purchases/analytics", ecommerceHandler.GetPurchaseAnalytics)
	// api.GET("/events/:id/purchases/top-products", ecommerceHandler.GetTopProducts)
	// api.POST("/webhooks/:integrationId", ecommerceHandler.HandleWebhook)

	// // Discount code routes
	// api.POST("/discount/generate", discountHandler.GenerateCode)
	// api.GET("/discount/codes/:code/validate", discountHandler.ValidateCode)
	// api.POST("/discount/redeem", discountHandler.RedeemCode)
	// api.GET("/events/:id/discount/analytics", discountHandler.GetCodeAnalytics)
	// api.GET("/brands/discount/codes", discountHandler.GetBrandCodes)

	// // Pixel tracking routes
	// r.GET("/pixel/track", pixelHandler.TrackPixel) // No auth for pixel tracking
	// api.POST("/pixel/search", pixelHandler.TrackSearch)
	// api.GET("/events/:id/pixel/analytics", pixelHandler.GetPixelAnalytics)
	// api.GET("/pixel/generate", pixelHandler.GeneratePixelURL)

	// // Advanced analytics routes
	// api.POST("/content/:id/ai-process", advancedAnalyticsHandler.ProcessContentAI)
	// api.GET("/brands/product-analytics", advancedAnalyticsHandler.GetProductAnalytics)
	// api.GET("/events/:id/conversion-funnel", advancedAnalyticsHandler.GetConversionFunnel)
	// api.GET("/events/:id/attribution-report", advancedAnalyticsHandler.GetAttributionReport)
	// api.POST("/conversions/track", advancedAnalyticsHandler.TrackConversion)

	// // Rewards and survey routes
	// api.GET("/users/rewards", rewardsHandler.GetUserRewards)
	// api.POST("/rewards/award", rewardsHandler.AwardReward)
	// api.POST("/rewards/process-quality", rewardsHandler.ProcessQualityRewards)
	// api.GET("/surveys/available", rewardsHandler.GetAvailableSurveys)
	// api.POST("/surveys/respond", rewardsHandler.SubmitSurveyResponse)
	// api.POST("/events/:id/surveys/schedule", rewardsHandler.ScheduleSurveys)
	// api.GET("/events/:id/surveys/analytics", rewardsHandler.GetSurveyAnalytics)

	// // Export and CRM routes
	// api.POST("/export/create", exportHandler.CreateExportRequest)
	// api.GET("/export/:requestId/status", exportHandler.GetExportStatus)
	// api.GET("/export/formats", exportHandler.GetExportFormats)
	// api.POST("/crm/integrations", exportHandler.CreateCRMIntegration)
	// api.POST("/crm/:integrationId/sync/:eventId", exportHandler.SyncEventData)
	// api.GET("/crm/types", exportHandler.GetCRMTypes)

	// // Performance routes
	// api.POST("/performance/optimize-db", performanceHandler.OptimizeDatabase)
	// api.GET("/performance/db-metrics", performanceHandler.GetDatabaseMetrics)
	// api.GET("/performance/query-stats", performanceHandler.GetQueryStats)
	// api.GET("/performance/cache-stats", performanceHandler.GetCacheStats)
	// api.DELETE("/performance/cache", performanceHandler.ClearCache)
	// api.POST("/performance/load-test", performanceHandler.RunLoadTest)
	// api.POST("/performance/maintenance", performanceHandler.RunMaintenance)

	// // Security routes
	// api.POST("/security/scan", securityHandler.RunSecurityScan)
	// api.GET("/security/events", securityHandler.GetSecurityEvents)
	// api.POST("/security/privacy/update", securityHandler.UpdatePrivacySettings)
	// api.POST("/users/data/delete", securityHandler.RequestDataDeletion)
	// api.GET("/users/data/export", securityHandler.ExportUserData)
	// api.POST("/users/data/anonymize", securityHandler.AnonymizeUserData)
	// api.POST("/security/validate-input", securityHandler.ValidateInput)

	// // UX routes
	// api.POST("/ux/sessions/start", uxHandler.StartUsabilitySession)
	// api.POST("/ux/actions/track", uxHandler.TrackUserAction)
	// api.POST("/ux/errors/track", uxHandler.TrackUserError)
	// api.POST("/ux/sessions/end", uxHandler.EndUsabilitySession)
	// api.GET("/ux/metrics", uxHandler.GetUsabilityMetrics)
	// api.GET("/ux/heatmap", uxHandler.GetHeatmapData)
	// api.GET("/ux/journey", uxHandler.GetUserJourney)
	// api.GET("/ux/pain-points", uxHandler.GetPainPoints)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Println("Starting server on :8080")
	log.Println("API documentation available at http://localhost:8080/swagger/index.html")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
