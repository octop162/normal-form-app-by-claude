package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CSRFTokenStore stores CSRF tokens with expiration
type CSRFTokenStore struct {
	tokens map[string]time.Time
	mutex  sync.RWMutex
}

// NewCSRFTokenStore creates a new CSRF token store
func NewCSRFTokenStore() *CSRFTokenStore {
	store := &CSRFTokenStore{
		tokens: make(map[string]time.Time),
	}
	// Start cleanup goroutine
	go store.cleanup()
	return store
}

// GenerateToken generates a new CSRF token
func (s *CSRFTokenStore) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(bytes)
	
	s.mutex.Lock()
	s.tokens[token] = time.Now().Add(4 * time.Hour) // 4 hour expiration
	s.mutex.Unlock()
	
	return token, nil
}

// ValidateToken validates a CSRF token
func (s *CSRFTokenStore) ValidateToken(token string) bool {
	s.mutex.RLock()
	expiration, exists := s.tokens[token]
	s.mutex.RUnlock()
	
	if !exists || time.Now().After(expiration) {
		return false
	}
	
	// Remove token after use (single use)
	s.mutex.Lock()
	delete(s.tokens, token)
	s.mutex.Unlock()
	
	return true
}

// cleanup removes expired tokens
func (s *CSRFTokenStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for token, expiration := range s.tokens {
			if now.After(expiration) {
				delete(s.tokens, token)
			}
		}
		s.mutex.Unlock()
	}
}

var csrfStore = NewCSRFTokenStore()

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", 
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data: https:; "+
			"font-src 'self'; "+
			"connect-src 'self'; "+
			"frame-ancestors 'none'")
		
		// HTTPS headers (for production)
		if gin.Mode() == gin.ReleaseMode {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		c.Next()
	}
}

// CSRF middleware for CSRF protection
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate token for GET requests to /api/v1/csrf-token
		if c.Request.Method == "GET" && c.Request.URL.Path == "/api/v1/csrf-token" {
			token, err := csrfStore.GenerateToken()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "CSRF_TOKEN_GENERATION_FAILED",
						"message": "Failed to generate CSRF token",
					},
				})
				c.Abort()
				return
			}
			
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"token": token,
				},
			})
			c.Abort()
			return
		}
		
		// Skip CSRF check for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}
		
		// Skip CSRF check for health endpoints
		if strings.HasPrefix(c.Request.URL.Path, "/health") {
			c.Next()
			return
		}
		
		// Get token from header
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "CSRF_TOKEN_MISSING",
					"message": "CSRF token is required",
				},
			})
			c.Abort()
			return
		}
		
		// Validate token
		if !csrfStore.ValidateToken(token) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "CSRF_TOKEN_INVALID",
					"message": "Invalid or expired CSRF token",
				},
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RateLimitStore stores request counts for rate limiting
type RateLimitStore struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
}

// NewRateLimitStore creates a new rate limit store
func NewRateLimitStore() *RateLimitStore {
	store := &RateLimitStore{
		requests: make(map[string][]time.Time),
	}
	// Start cleanup goroutine
	go store.cleanup()
	return store
}

// IsAllowed checks if a request is allowed based on rate limiting
func (s *RateLimitStore) IsAllowed(key string, limit int, window time.Duration) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-window)
	
	// Get existing requests for this key
	requests := s.requests[key]
	
	// Filter out old requests
	validRequests := make([]time.Time, 0)
	for _, req := range requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	
	// Check if limit exceeded
	if len(validRequests) >= limit {
		return false
	}
	
	// Add current request
	validRequests = append(validRequests, now)
	s.requests[key] = validRequests
	
	return true
}

// cleanup removes old request records
func (s *RateLimitStore) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		cutoff := now.Add(-1 * time.Hour) // Keep 1 hour of data
		
		for key, requests := range s.requests {
			validRequests := make([]time.Time, 0)
			for _, req := range requests {
				if req.After(cutoff) {
					validRequests = append(validRequests, req)
				}
			}
			
			if len(validRequests) == 0 {
				delete(s.requests, key)
			} else {
				s.requests[key] = validRequests
			}
		}
		s.mutex.Unlock()
	}
}

var rateLimitStore = NewRateLimitStore()

// RateLimit middleware for rate limiting
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as key
		key := c.ClientIP()
		
		if !rateLimitStore.IsAllowed(key, limit, window) {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
			c.Header("X-RateLimit-Window", window.String())
			c.Header("Retry-After", fmt.Sprintf("%.0f", window.Seconds()))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please try again later.",
				},
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// InputSanitization middleware for input sanitization
func InputSanitization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add sanitization headers
		c.Header("X-Content-Type-Options", "nosniff")
		
		// For JSON requests, ensure content type is correct
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if contentType != "" && !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "UNSUPPORTED_MEDIA_TYPE",
						"message": "Content-Type must be application/json",
					},
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}