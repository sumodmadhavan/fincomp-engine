package runout

import (
	"testing"
	"time"
)

func TestCalculate(t *testing.T) {
	params := RunoutParams{
		ContractStartDate:  time.Date(2022, 1, 14, 0, 0, 0, 0, time.UTC),
		ContractEndDate:    time.Date(2034, 2, 14, 0, 0, 0, 0, time.UTC),
		AUHours:            480,
		WarrantyRate:       243.6,
		FirstRunRate:       255.13,
		SecondRunRate:      255.13,
		ThirdRunRate:       255.13,
		ManagementFees:     15.0,
		AICFees:            20.0,
		TrustLoadFees:      2.98,
		BuyIn:              1352291.05,
		RateEscalation:     8.75,
		FlightHoursMinimum: 150,
	}

	result, err := Calculate(params)
	if err != nil {
		t.Fatalf("Calculate returned an error: %v", err)
	}

	// Add assertions to check the correctness of the result
	// Compare with expected values from the Python script output
}
