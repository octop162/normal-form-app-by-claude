// Package main provides a mock server for external APIs used in development and testing.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort            = "8081"
	shutdownTimeoutSeconds = 30
)

// Mock data structures
type InventoryCheckRequest struct {
	OptionIDs []string `json:"option_ids"`
}

type InventoryCheckResponse struct {
	Success bool           `json:"success"`
	Data    map[string]int `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

type RegionCheckRequest struct {
	Prefecture string   `json:"prefecture"`
	City       string   `json:"city"`
	OptionIDs  []string `json:"option_ids"`
}

type RegionCheckResponse struct {
	Success bool            `json:"success"`
	Data    map[string]bool `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type AddressSearchRequest struct {
	PostalCode string `json:"postal_code"`
}

type AddressSearchResponse struct {
	Success bool         `json:"success"`
	Data    *AddressData `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

type AddressData struct {
	PostalCode string `json:"postal_code"`
	Prefecture string `json:"prefecture"`
	City       string `json:"city"`
	Town       string `json:"town,omitempty"`
}

func main() {
	port := getEnv("MOCK_PORT", defaultPort)
	
	// Set Gin to release mode for production-like behavior
	gin.SetMode(gin.ReleaseMode)
	
	r := setupRouter()
	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Mock API server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down mock server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Mock server exited")
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "mock-api-server",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Inventory API
		api.POST("/inventory/check", handleInventoryCheck)
		
		// Region API
		api.POST("/region/check", handleRegionCheck)
		
		// Address API
		api.POST("/address/search", handleAddressSearch)
	}

	return r
}

// handleInventoryCheck handles inventory check requests
func handleInventoryCheck(c *gin.Context) {
	var req InventoryCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, InventoryCheckResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	if len(req.OptionIDs) == 0 {
		c.JSON(http.StatusBadRequest, InventoryCheckResponse{
			Success: false,
			Error:   "Option IDs cannot be empty",
		})
		return
	}

	// Mock inventory data
	inventory := make(map[string]int)
	for _, optionID := range req.OptionIDs {
		switch optionID {
		case "AA":
			inventory[optionID] = 15 // Good stock
		case "BB":
			inventory[optionID] = 0 // Out of stock
		case "AB":
			inventory[optionID] = 5 // Low stock
		case "TEST":
			inventory[optionID] = 100 // Test option
		default:
			inventory[optionID] = 3 // Default stock
		}
	}

	// Simulate occasional API failures (5% chance)
	if shouldSimulateError(5) {
		c.JSON(http.StatusInternalServerError, InventoryCheckResponse{
			Success: false,
			Error:   "Temporary inventory service unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, InventoryCheckResponse{
		Success: true,
		Data:    inventory,
	})
}

// handleRegionCheck handles region restriction check requests
func handleRegionCheck(c *gin.Context) {
	var req RegionCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, RegionCheckResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	if req.Prefecture == "" || req.City == "" {
		c.JSON(http.StatusBadRequest, RegionCheckResponse{
			Success: false,
			Error:   "Prefecture and city are required",
		})
		return
	}

	if len(req.OptionIDs) == 0 {
		c.JSON(http.StatusBadRequest, RegionCheckResponse{
			Success: false,
			Error:   "Option IDs cannot be empty",
		})
		return
	}

	// Mock region restrictions
	restrictions := make(map[string]bool)
	for _, optionID := range req.OptionIDs {
		allowed := checkMockRegionRestriction(req.Prefecture, req.City, optionID)
		restrictions[optionID] = allowed
	}

	// Simulate occasional API failures (3% chance)
	if shouldSimulateError(3) {
		c.JSON(http.StatusInternalServerError, RegionCheckResponse{
			Success: false,
			Error:   "Temporary region service unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, RegionCheckResponse{
		Success: true,
		Data:    restrictions,
	})
}

// handleAddressSearch handles address search requests
func handleAddressSearch(c *gin.Context) {
	var req AddressSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, AddressSearchResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	if len(req.PostalCode) != 7 {
		c.JSON(http.StatusBadRequest, AddressSearchResponse{
			Success: false,
			Error:   "Postal code must be 7 digits",
		})
		return
	}

	// Mock address data
	address := getMockAddressData(req.PostalCode)
	if address == nil {
		c.JSON(http.StatusOK, AddressSearchResponse{
			Success: false,
			Error:   "Address not found for postal code: " + req.PostalCode,
		})
		return
	}

	// Simulate occasional API failures (2% chance)
	if shouldSimulateError(2) {
		c.JSON(http.StatusInternalServerError, AddressSearchResponse{
			Success: false,
			Error:   "Temporary address service unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, AddressSearchResponse{
		Success: true,
		Data:    address,
	})
}

// checkMockRegionRestriction simulates region restriction logic
func checkMockRegionRestriction(prefecture, _ string, optionID string) bool {
	switch optionID {
	case "AA":
		// AA option restricted in Hokkaido
		return prefecture != "北海道"
	case "BB":
		// BB option restricted in major metropolitan areas
		restrictedPrefectures := []string{"東京都", "大阪府", "愛知県"}
		return !slices.Contains(restrictedPrefectures, prefecture)
	case "AB":
		// AB option available everywhere
		return true
	case "TEST":
		// Test option for health checks
		return true
	default:
		// Default: allow most options except in Okinawa (for testing)
		return prefecture != "沖縄県"
	}
}

// getMockAddressData returns mock address data for testing
func getMockAddressData(postalCode string) *AddressData {
	addressMap := map[string]*AddressData{
		"1000001": {
			PostalCode: "100-0001",
			Prefecture: "東京都",
			City:       "千代田区",
			Town:       "千代田",
		},
		"1000005": {
			PostalCode: "100-0005",
			Prefecture: "東京都",
			City:       "千代田区",
			Town:       "丸の内",
		},
		"1500002": {
			PostalCode: "150-0002",
			Prefecture: "東京都",
			City:       "渋谷区",
			Town:       "渋谷",
		},
		"5410041": {
			PostalCode: "541-0041",
			Prefecture: "大阪府",
			City:       "大阪市中央区",
			Town:       "北浜",
		},
		"2310023": {
			PostalCode: "231-0023",
			Prefecture: "神奈川県",
			City:       "横浜市中区",
			Town:       "山下町",
		},
		"4600008": {
			PostalCode: "460-0008",
			Prefecture: "愛知県",
			City:       "名古屋市中区",
			Town:       "栄",
		},
		"8100001": {
			PostalCode: "810-0001",
			Prefecture: "福岡県",
			City:       "福岡市中央区",
			Town:       "天神",
		},
		"0600001": {
			PostalCode: "060-0001",
			Prefecture: "北海道",
			City:       "札幌市中央区",
			Town:       "北一条西",
		},
		"9000006": {
			PostalCode: "900-0006",
			Prefecture: "沖縄県",
			City:       "那覇市",
			Town:       "おもろまち",
		},
	}

	return addressMap[postalCode]
}

// shouldSimulateError returns true if an error should be simulated based on percentage
func shouldSimulateError(percentage int) bool {
	if percentage <= 0 || percentage >= 100 {
		return false
	}
	
	// Simple pseudo-random based on current time
	now := time.Now().UnixNano()
	return int(now%100) < percentage
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

