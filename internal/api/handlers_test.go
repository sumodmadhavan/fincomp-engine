package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"financialapi/internal/financials"
	"financialapi/pkg/testutils"

	"github.com/gin-gonic/gin"
)

func TestGoalSeekHandler(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup
	router := gin.Default()
	server := &Server{router: router}
	server.setupRoutes()

	// Test data
	params := financials.FinancialParams{
		NumYears:       10,
		AuHours:        450,
		InitialTSN:     100,
		RateEscalation: 5,
		AIC:            10,
		HSITSN:         1000,
		OverhaulTSN:    3000,
		HSICost:        50000,
		OverhaulCost:   100000,
		TargetProfit:   3000000,
		InitialRate:    320,
	}

	// Convert params to JSON
	jsonParams, _ := json.Marshal(params)

	// Create request
	req, _ := http.NewRequest("POST", "/goalseek", bytes.NewBuffer(jsonParams))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check response
	testutils.AssertEqual(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	_, exists := response["optimalWarrantyRate"]
	testutils.AssertEqual(t, true, exists)

	_, exists = response["iterations"]
	testutils.AssertEqual(t, true, exists)

	_, exists = response["finalCumulativeProfit"]
	testutils.AssertEqual(t, true, exists)
}
