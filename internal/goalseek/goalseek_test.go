// File: internal/goalseek/goalseek_test.go

package goalseek

import (
	"financialapi/internal/financials"
	"financialapi/pkg/testutils"
	"testing"
)

func TestGoalSeekEngine(t *testing.T) {
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

	engine := NewGoalSeekCalculator(params)

	// Test Initialize
	err := engine.Initialize(params)
	testutils.AssertNoError(t, err)

	// Test Validate
	err = engine.Validate()
	testutils.AssertNoError(t, err)

	// Test Compute
	err = engine.Compute()
	testutils.AssertNoError(t, err)

	// Test GetResult
	result := engine.GetResult()
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("GetResult did not return a map[string]interface{}")
	}

	optimalRate, ok := resultMap["optimalWarrantyRate"].(float64)
	if !ok {
		t.Fatalf("optimalWarrantyRate is not a float64")
	}

	if optimalRate < 300 || optimalRate > 600 {
		t.Errorf("Expected optimalWarrantyRate between 300 and 600, got %f", optimalRate)
	}
}