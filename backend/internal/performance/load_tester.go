/**
 * Load Tester
 * Performance testing utilities for load testing
 */

package performance

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type LoadTestResult struct {
	TotalRequests   int           `json:"totalRequests"`
	SuccessfulReqs  int           `json:"successfulRequests"`
	FailedReqs      int           `json:"failedRequests"`
	AvgResponseTime time.Duration `json:"avgResponseTime"`
	MinResponseTime time.Duration `json:"minResponseTime"`
	MaxResponseTime time.Duration `json:"maxResponseTime"`
	RequestsPerSec  float64       `json:"requestsPerSecond"`
	TestDuration    time.Duration `json:"testDuration"`
}

type LoadTester struct {
	client *http.Client
}

func NewLoadTester() *LoadTester {
	return &LoadTester{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (lt *LoadTester) RunLoadTest(url string, concurrency int, duration time.Duration) *LoadTestResult {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	result := &LoadTestResult{
		MinResponseTime: time.Hour, // Initialize with high value
	}

	startTime := time.Now()
	// endTime := startTime.Add(duration)

	// Channel to signal workers to stop
	stopChan := make(chan struct{})

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lt.worker(url, stopChan, result, &mutex)
		}()
	}

	// Stop workers after duration
	time.AfterFunc(duration, func() {
		close(stopChan)
	})

	wg.Wait()

	result.TestDuration = time.Since(startTime)
	if result.TotalRequests > 0 {
		result.RequestsPerSec = float64(result.TotalRequests) / result.TestDuration.Seconds()
	}

	return result
}

func (lt *LoadTester) worker(url string, stopChan <-chan struct{}, result *LoadTestResult, mutex *sync.Mutex) {
	var totalResponseTime time.Duration

	for {
		select {
		case <-stopChan:
			return
		default:
			start := time.Now()
			resp, err := lt.client.Get(url)
			responseTime := time.Since(start)

			mutex.Lock()
			result.TotalRequests++
			totalResponseTime += responseTime

			if err != nil || resp.StatusCode >= 400 {
				result.FailedReqs++
			} else {
				result.SuccessfulReqs++
			}

			if resp != nil {
				resp.Body.Close()
			}

			// Update min/max response times
			if responseTime < result.MinResponseTime {
				result.MinResponseTime = responseTime
			}
			if responseTime > result.MaxResponseTime {
				result.MaxResponseTime = responseTime
			}

			result.AvgResponseTime = totalResponseTime / time.Duration(result.TotalRequests)
			mutex.Unlock()
		}
	}
}

func (lt *LoadTester) TestEndpoints(endpoints []string, concurrency int, duration time.Duration) map[string]*LoadTestResult {
	results := make(map[string]*LoadTestResult)

	for _, endpoint := range endpoints {
		fmt.Printf("Testing endpoint: %s\n", endpoint)
		results[endpoint] = lt.RunLoadTest(endpoint, concurrency, duration)
	}

	return results
}

func (lt *LoadTester) GenerateReport(results map[string]*LoadTestResult) string {
	report := "Load Test Report\n"
	report += "================\n\n"

	for endpoint, result := range results {
		report += fmt.Sprintf("Endpoint: %s\n", endpoint)
		report += fmt.Sprintf("Total Requests: %d\n", result.TotalRequests)
		report += fmt.Sprintf("Successful: %d (%.2f%%)\n", result.SuccessfulReqs, float64(result.SuccessfulReqs)/float64(result.TotalRequests)*100)
		report += fmt.Sprintf("Failed: %d (%.2f%%)\n", result.FailedReqs, float64(result.FailedReqs)/float64(result.TotalRequests)*100)
		report += fmt.Sprintf("Avg Response Time: %v\n", result.AvgResponseTime)
		report += fmt.Sprintf("Min Response Time: %v\n", result.MinResponseTime)
		report += fmt.Sprintf("Max Response Time: %v\n", result.MaxResponseTime)
		report += fmt.Sprintf("Requests/sec: %.2f\n", result.RequestsPerSec)
		report += fmt.Sprintf("Test Duration: %v\n\n", result.TestDuration)
	}

	return report
}
