package goalseek

import (
	"financialapi/internal/financials"
	"financialapi/pkg/testutils"
	"testing"
)

func TestCalculate(t *testing.T) {
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

	result, err := Calculate(params)
	testutils.AssertNoError(t, err)

	optimalRate, ok := result["optimalWarrantyRate"].(float64)
	if !ok {
		t.Fatalf("optimalWarrantyRate is not a float64")
	}

	if optimalRate < 300 || optimalRate > 600 {
		t.Errorf("Expected optimalWarrantyRate between 300 and 600, got %f", optimalRate)
	}
}
func TestGoalSeekCalculator(t *testing.T) {
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

	calculator := NewGoalSeekCalculator(params)

	// Test CalculateCumulativeProfit
	profit, err := calculator.CalculateCumulativeProfit()
	testutils.AssertNoError(t, err)
	if profit <= 0 {
		t.Errorf("Expected positive profit, got %f", profit)
	}

	// Test GetParams
	retrievedParams := calculator.GetParams()
	if retrievedParams != params {
		t.Errorf("GetParams did not return the correct parameters")
	}

	// Test SetParams
	newParams := financials.FinancialParams{
		NumYears:     5,
		TargetProfit: 1500000,
		InitialRate:  300,
	}
	err = calculator.SetParams(newParams)
	testutils.AssertNoError(t, err)

	updatedParams := calculator.GetParams()
	if updatedParams != newParams {
		t.Errorf("SetParams did not update the parameters correctly")
	}
}
