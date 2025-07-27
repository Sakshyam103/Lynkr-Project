/**
 * Performance Handlers
 * HTTP handlers for performance monitoring and optimization
 */

package handlers

import (
	"time"

	"lynkr/internal/performance"

	"github.com/gin-gonic/gin"
)

type PerformanceHandler struct {
	dbOptimizer *performance.DatabaseOptimizer
	cache       *performance.Cache
	loadTester  *performance.LoadTester
}

func NewPerformanceHandler(dbOptimizer *performance.DatabaseOptimizer, cache *performance.Cache, loadTester *performance.LoadTester) *PerformanceHandler {
	return &PerformanceHandler{
		dbOptimizer: dbOptimizer,
		cache:       cache,
		loadTester:  loadTester,
	}
}

func (ph *PerformanceHandler) OptimizeDatabase(c *gin.Context) {
	if err := ph.dbOptimizer.OptimizeIndexes(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to optimize database"})
		return
	}
	if err := ph.dbOptimizer.OptimizeQueries(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to optimize queries"})
		return
	}

	c.JSON(200, gin.H{"status": "optimized"})
}

func (ph *PerformanceHandler) GetDatabaseMetrics(c *gin.Context) {
	metrics, err := ph.dbOptimizer.GetPerformanceMetrics()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get database metrics"})
		return
	}
	c.JSON(200, metrics)
}

func (ph *PerformanceHandler) GetQueryStats(c *gin.Context) {
	stats, err := ph.dbOptimizer.AnalyzeQueries()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to analyze queries"})
		return
	}
	c.JSON(200, gin.H{"queries": stats})
}

func (ph *PerformanceHandler) GetCacheStats(c *gin.Context) {
	stats := ph.cache.GetStats()
	c.JSON(200, stats)
}

func (ph *PerformanceHandler) ClearCache(c *gin.Context) {
	pattern := c.Query("pattern")
	if pattern != "" {
		ph.cache.InvalidatePattern(pattern)
	} else {
		ph.cache.Clear()
	}
	c.JSON(200, gin.H{"status": "cleared"})
}

func (ph *PerformanceHandler) RunLoadTest(c *gin.Context) {
	var request struct {
		URL         string `json:"url"`
		Concurrency int    `json:"concurrency"`
		Duration    int    `json:"duration"` // seconds
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if request.Concurrency == 0 {
		request.Concurrency = 10
	}
	if request.Duration == 0 {
		request.Duration = 30
	}

	result := ph.loadTester.RunLoadTest(
		request.URL,
		request.Concurrency,
		time.Duration(request.Duration)*time.Second,
	)

	c.JSON(200, result)
}

func (ph *PerformanceHandler) RunMaintenance(c *gin.Context) {
	if err := ph.dbOptimizer.RunMaintenance(); err != nil {
		c.JSON(500, gin.H{"error": "Failed to run maintenance"})
		return
	}
	c.JSON(200, gin.H{"status": "maintenance_completed"})
}
