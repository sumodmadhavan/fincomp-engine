package runout

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

	message, ok := result["message"].(string)
	if !ok {
		t.Fatalf("message is not a string")
	}

	expectedMessage := "Runout calculation is under development"
	testutils.AssertEqual(t, expectedMessage, message)
}
