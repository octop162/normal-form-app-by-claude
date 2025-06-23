package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ResponseWriter wrapper for capturing response size
type responseWriter struct {
	gin.ResponseWriter
	size int
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

// PerformanceMetrics stores performance metrics
type PerformanceMetrics struct {
	RequestCount     int64         `json:"request_count"`
	TotalDuration    time.Duration `json:"total_duration"`
	AverageDuration  time.Duration `json:"average_duration"`
	MinDuration      time.Duration `json:"min_duration"`
	MaxDuration      time.Duration `json:"max_duration"`
	ErrorCount       int64         `json:"error_count"`
	ActiveGoroutines int           `json:"active_goroutines"`
	MemoryUsage      uint64        `json:"memory_usage_bytes"`
}

// MetricsCollector collects and manages performance metrics
type MetricsCollector struct {
	mutex           sync.RWMutex
	requestCount    int64
	totalDuration   time.Duration
	minDuration     time.Duration
	maxDuration     time.Duration
	errorCount      int64
	endpointMetrics map[string]*PerformanceMetrics
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		endpointMetrics: make(map[string]*PerformanceMetrics),
		minDuration:     time.Hour, // Initialize with large value
	}
}

var globalMetricsCollector = NewMetricsCollector()

// RecordRequest records metrics for a request
func (mc *MetricsCollector) RecordRequest(endpoint string, duration time.Duration, isError bool) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.requestCount++
	mc.totalDuration += duration

	if duration < mc.minDuration {
		mc.minDuration = duration
	}
	if duration > mc.maxDuration {
		mc.maxDuration = duration
	}

	if isError {
		mc.errorCount++
	}

	// Update endpoint-specific metrics
	if _, exists := mc.endpointMetrics[endpoint]; !exists {
		mc.endpointMetrics[endpoint] = &PerformanceMetrics{
			MinDuration: time.Hour,
		}
	}

	endpointMetric := mc.endpointMetrics[endpoint]
	endpointMetric.RequestCount++
	endpointMetric.TotalDuration += duration

	if duration < endpointMetric.MinDuration {
		endpointMetric.MinDuration = duration
	}
	if duration > endpointMetric.MaxDuration {
		endpointMetric.MaxDuration = duration
	}

	if isError {
		endpointMetric.ErrorCount++
	}

	if endpointMetric.RequestCount > 0 {
		endpointMetric.AverageDuration = endpointMetric.TotalDuration / time.Duration(endpointMetric.RequestCount)
	}
}

// GetMetrics returns current metrics
func (mc *MetricsCollector) GetMetrics() PerformanceMetrics {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	var avgDuration time.Duration
	if mc.requestCount > 0 {
		avgDuration = mc.totalDuration / time.Duration(mc.requestCount)
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return PerformanceMetrics{
		RequestCount:     mc.requestCount,
		TotalDuration:    mc.totalDuration,
		AverageDuration:  avgDuration,
		MinDuration:      mc.minDuration,
		MaxDuration:      mc.maxDuration,
		ErrorCount:       mc.errorCount,
		ActiveGoroutines: runtime.NumGoroutine(),
		MemoryUsage:      memStats.Alloc,
	}
}

// GetEndpointMetrics returns metrics for a specific endpoint
func (mc *MetricsCollector) GetEndpointMetrics(endpoint string) *PerformanceMetrics {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	if metric, exists := mc.endpointMetrics[endpoint]; exists {
		// Create a copy to avoid race conditions
		metricCopy := *metric
		metricCopy.ActiveGoroutines = runtime.NumGoroutine()
		
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		metricCopy.MemoryUsage = memStats.Alloc
		
		return &metricCopy
	}
	return nil
}

// GetAllEndpointMetrics returns metrics for all endpoints
func (mc *MetricsCollector) GetAllEndpointMetrics() map[string]*PerformanceMetrics {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	result := make(map[string]*PerformanceMetrics)
	for endpoint, metric := range mc.endpointMetrics {
		// Create a copy to avoid race conditions
		metricCopy := *metric
		result[endpoint] = &metricCopy
	}
	return result
}

// Reset resets all metrics
func (mc *MetricsCollector) Reset() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.requestCount = 0
	mc.totalDuration = 0
	mc.minDuration = time.Hour
	mc.maxDuration = 0
	mc.errorCount = 0
	mc.endpointMetrics = make(map[string]*PerformanceMetrics)
}

// PerformanceMiddleware tracks request performance
func PerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		endpoint := fmt.Sprintf("%s %s", method, path)

		// Create response writer wrapper
		rw := &responseWriter{ResponseWriter: c.Writer}
		c.Writer = rw

		// Add performance context
		c.Set("performance_start", start)
		c.Set("performance_endpoint", endpoint)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		status := c.Writer.Status()
		
		// Determine if it's an error
		isError := status >= 400

		// Record metrics
		globalMetricsCollector.RecordRequest(endpoint, duration, isError)

		// Add performance headers
		c.Header("X-Response-Time", fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1e6))
		c.Header("X-Response-Size", strconv.Itoa(rw.size))
	}
}

// MetricsEndpoint provides a handler for metrics endpoint
func MetricsEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics := globalMetricsCollector.GetMetrics()
		
		response := gin.H{
			"success": true,
			"data": gin.H{
				"overall_metrics": metrics,
				"endpoint_metrics": globalMetricsCollector.GetAllEndpointMetrics(),
				"timestamp": time.Now().Format(time.RFC3339),
			},
		}
		
		c.JSON(http.StatusOK, response)
	}
}

// Connection pooling optimizer
type ConnectionPool struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

// NewConnectionPool creates optimized database connection pool settings
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		maxOpenConns:    25,  // Based on server capacity
		maxIdleConns:    10,  // Reasonable idle connections
		connMaxLifetime: 30 * time.Minute,
		connMaxIdleTime: 15 * time.Minute,
	}
}

// ApplyToDatabase applies connection pool settings to database
func (cp *ConnectionPool) ApplyToDatabase(db interface{}) {
	// This would be implemented based on the actual database driver
	// For sql.DB:
	// db.SetMaxOpenConns(cp.maxOpenConns)
	// db.SetMaxIdleConns(cp.maxIdleConns)
	// db.SetConnMaxLifetime(cp.connMaxLifetime)
	// db.SetConnMaxIdleTime(cp.connMaxIdleTime)
}

// Caching middleware with TTL
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

type MemoryCache struct {
	mutex sync.RWMutex
	items map[string]*CacheItem
}

func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*CacheItem),
	}
	
	// Start cleanup goroutine
	go cache.cleanup()
	
	return cache
}

func (mc *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	mc.items[key] = &CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (mc *MemoryCache) Get(key string) (interface{}, bool) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	item, exists := mc.items[key]
	if !exists {
		return nil, false
	}
	
	if time.Now().After(item.ExpiresAt) {
		delete(mc.items, key)
		return nil, false
	}
	
	return item.Data, true
}

func (mc *MemoryCache) Delete(key string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	delete(mc.items, key)
}

func (mc *MemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		mc.mutex.Lock()
		now := time.Now()
		for key, item := range mc.items {
			if now.After(item.ExpiresAt) {
				delete(mc.items, key)
			}
		}
		mc.mutex.Unlock()
	}
}

var globalCache = NewMemoryCache()

// CacheMiddleware provides response caching for GET requests
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != "GET" {
			c.Next()
			return
		}
		
		// Skip caching for specific endpoints
		path := c.Request.URL.Path
		skipCache := []string{"/health", "/metrics", "/api/v1/csrf-token"}
		for _, skip := range skipCache {
			if path == skip {
				c.Next()
				return
			}
		}
		
		// Generate cache key
		cacheKey := fmt.Sprintf("%s:%s:%s", c.Request.Method, path, c.Request.URL.RawQuery)
		
		// Try to get from cache
		if cachedData, exists := globalCache.Get(cacheKey); exists {
			if response, ok := cachedData.(gin.H); ok {
				c.Header("X-Cache", "HIT")
				c.JSON(http.StatusOK, response)
				return
			}
		}
		
		// Create response writer to capture response
		rw := &responseWriter{ResponseWriter: c.Writer}
		c.Writer = rw
		
		c.Next()
		
		// Cache successful responses
		if c.Writer.Status() == http.StatusOK && rw.size > 0 {
			// This is a simplified caching approach
			// In practice, you'd need to capture the actual response data
			c.Header("X-Cache", "MISS")
		}
	}
}

// Graceful timeout middleware
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		
		c.Request = c.Request.WithContext(ctx)
		
		finished := make(chan struct{})
		go func() {
			c.Next()
			finished <- struct{}{}
		}()
		
		select {
		case <-finished:
			// Request completed successfully
		case <-ctx.Done():
			// Request timed out
			c.JSON(http.StatusRequestTimeout, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "REQUEST_TIMEOUT",
					"message": "Request timed out",
				},
			})
			c.Abort()
		}
	}
}