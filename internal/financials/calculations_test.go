package financials

import (
	"financialapi/pkg/testutils"
	"testing"
)

func TestCalculateFinancials(t *testing.T) {
	params := FinancialParams{
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

	profit, err := CalculateFinancials(params.InitialRate, params)
	testutils.AssertNoError(t, err)

	// This is a very basic check. You'd want to calculate the expected
	// profit based on the input parameters and check against that.
	if profit <= 0 {
		t.Errorf("Expected positive profit, got %f", profit)
	}
}

func TestGoalSeek(t *testing.T) {
	params := FinancialParams{
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

	optimalRate, iterations, err := GoalSeek(params.TargetProfit, params, params.InitialRate)
	testutils.AssertNoError(t, err)

	if optimalRate <= params.InitialRate {
		t.Errorf("Expected optimal rate > %f, got %f", params.InitialRate, optimalRate)
	}

	if iterations <= 0 {
		t.Errorf("Expected positive number of iterations, got %d", iterations)
	}
}
