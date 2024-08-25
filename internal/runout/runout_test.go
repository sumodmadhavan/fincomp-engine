package runout

import (
	"math"
	"testing"
	"time"
)

func TestCalculate(t *testing.T) {
	params := RunoutParams{
		ContractStartDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		ContractEndDate:    time.Date(2034, 12, 31, 23, 59, 59, 0, time.UTC),
		AUHours:            480,
		WarrantyRate:       243.6,
		FirstRunRate:       255.13,
		SecondRunRate:      255.13,
		ThirdRunRate:       255.13,
		ManagementFees:     15.0,
		AICFees:            20.0,
		TrustLoadFees:      2.98,
		BuyIn:              1352291,
		RateEscalation:     8.75,
		FlightHoursMinimum: 150,
		NumOfDaysInYear:    365,
		NumOfDaysInMonth:   30,
		EnrollmentFees:     25000,
		NumEngines:         2,
		EngineParams: []EngineParams{
			{
				WarrantyExpDate:         time.Date(2025, 10, 31, 23, 59, 59, 0, time.UTC),
				WarrantyExpHours:        1000,
				FirstRunRateSwitchDate:  time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
				SecondRunRateSwitchDate: time.Date(2027, 5, 1, 0, 0, 0, 0, time.UTC),
				ThirdRunRateSwitchDate:  time.Date(2028, 7, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				WarrantyExpDate:         time.Date(2025, 10, 31, 23, 59, 59, 0, time.UTC),
				WarrantyExpHours:        1000,
				FirstRunRateSwitchDate:  time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
				SecondRunRateSwitchDate: time.Date(2027, 5, 1, 0, 0, 0, 0, time.UTC),
				ThirdRunRateSwitchDate:  time.Date(2028, 7, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	result, err := Calculate(params)
	if err != nil {
		t.Fatalf("Calculate returned an error: %v", err)
	}

	// Check the number of periods
	expectedPeriods := 12
	if len(result.Periods) != expectedPeriods {
		t.Errorf("Expected %d periods, got %d", expectedPeriods, len(result.Periods))
	}

	// Check the first period
	firstPeriod := result.Periods[0]
	if !firstPeriod.StartDate.Equal(params.ContractStartDate) {
		t.Errorf("First period start date incorrect. Expected %v, got %v", params.ContractStartDate, firstPeriod.StartDate)
	}
	if firstPeriod.NumOfDays != 365 {
		t.Errorf("First period number of days incorrect. Expected 365, got %d", firstPeriod.NumOfDays)
	}
	if !almostEqual(firstPeriod.RateTrend, 1.0, 0.0001) {
		t.Errorf("First period rate trend incorrect. Expected 1.0, got %f", firstPeriod.RateTrend)
	}
	if !almostEqual(firstPeriod.TotalFHRevenue, 1586147, 0.01) {
		t.Errorf("First period total FH revenue incorrect. Expected 1586147, got %f", firstPeriod.TotalFHRevenue)
	}

	// Check engine data for the first period
	engine1 := firstPeriod.Engines[0]
	if engine1.WarrantyRateDays != 365 {
		t.Errorf("First period, engine 1 warranty rate days incorrect. Expected 365, got %d", engine1.WarrantyRateDays)
	}
	if !almostEqual(engine1.FHRevenue, 116928, 0.01) {
		t.Errorf("First period, engine 1 FH revenue incorrect. Expected 116928, got %f", engine1.FHRevenue)
	}

	// Check the last period
	lastPeriod := result.Periods[len(result.Periods)-1]
	if !lastPeriod.EndDate.Equal(params.ContractEndDate) {
		t.Errorf("Last period end date incorrect. Expected %v, got %v", params.ContractEndDate, lastPeriod.EndDate)
	}
	if !almostEqual(lastPeriod.RateTrend, 2.516065, 0.0001) {
		t.Errorf("Last period rate trend incorrect. Expected 2.516065, got %f", lastPeriod.RateTrend)
	}

	// Check overall totals
	if !almostEqual(result.TotalFHRevenue, 6179552.40, 0.01) {
		t.Errorf("Total FH revenue incorrect. Expected 6179552.40, got %f", result.TotalFHRevenue)
	}
	if !almostEqual(result.MgmtFeeRevenue, 724089.21, 0.01) {
		t.Errorf("Management fee revenue incorrect. Expected 724089.21, got %f", result.MgmtFeeRevenue)
	}
	if !almostEqual(result.AICRevenue, 820634.44, 0.01) {
		t.Errorf("AIC revenue incorrect. Expected 820634.44, got %f", result.AICRevenue)
	}
	if !almostEqual(result.TrustLoadRevenue, 122274.53, 0.01) {
		t.Errorf("Trust load revenue incorrect. Expected 122274.53, got %f", result.TrustLoadRevenue)
	}
	if !almostEqual(result.TrustRevenue, 4512554.22, 0.01) {
		t.Errorf("Trust revenue incorrect. Expected 4512554.22, got %f", result.TrustRevenue)
	}

	// Check Buy-In and Enrollment Fees
	if result.BuyIn != params.BuyIn {
		t.Errorf("Buy-In incorrect. Expected %f, got %f", params.BuyIn, result.BuyIn)
	}
	if result.EnrollmentFees != params.EnrollmentFees {
		t.Errorf("Enrollment Fees incorrect. Expected %f, got %f", params.EnrollmentFees, result.EnrollmentFees)
	}
}

// Helper function to compare float64 values with a tolerance
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
