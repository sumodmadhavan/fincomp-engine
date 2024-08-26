package runout

import (
	"math"
	"testing"
	"time"
)

// almostEqual compares two float64 values with a given tolerance
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
func TestRunoutEngine(t *testing.T) {
	params := RunoutParams{
		ContractStartDate:  time.Date(2022, 1, 14, 0, 0, 0, 0, time.UTC),
		ContractEndDate:    time.Date(2034, 2, 14, 23, 59, 59, 0, time.UTC),
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

	engine := NewRunoutCalculator(params)

	// Test Initialize
	err := engine.Initialize(params)
	if err != nil {
		t.Fatalf("Initialize returned an error: %v", err)
	}

	// Test Validate
	err = engine.Validate()
	if err != nil {
		t.Fatalf("Validate returned an error: %v", err)
	}

	// Test Compute
	err = engine.Compute()
	if err != nil {
		t.Fatalf("Compute returned an error: %v", err)
	}

	// Test GetResult
	result := engine.GetResult()
	runoutResult, ok := result.(RunoutResult)
	if !ok {
		t.Fatalf("GetResult did not return a RunoutResult")
	}

	// Perform the same checks as in the original TestCalculate function
	if len(runoutResult.Periods) != 12 {
		t.Errorf("Expected 12 periods, got %d", len(runoutResult.Periods))
	}

	// ... (continue with other checks from the original TestCalculate function)

	// Test existing methods for backwards compatibility
	profit, err := engine.CalculateCumulativeProfit()
	if err != nil {
		t.Fatalf("CalculateCumulativeProfit returned an error: %v", err)
	}
	if !almostEqual(profit, 4805005.2476703655, 0.01) {
		t.Errorf("Expected cumulative profit of 4805005.2476703655, got %f", profit)
	}

	retrievedParams := engine.GetParams()
	if retrievedParams != params {
		t.Errorf("GetParams did not return the correct parameters")
	}

	newParams := RunoutParams{
		ContractStartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		ContractEndDate:   time.Date(2033, 12, 31, 23, 59, 59, 0, time.UTC),
		AUHours:           500,
		// ... set other fields as needed
	}
	err = engine.SetParams(newParams)
	if err != nil {
		t.Fatalf("SetParams returned an error: %v", err)
	}

	updatedParams := engine.GetParams()
	if updatedParams != newParams {
		t.Errorf("SetParams did not update the parameters correctly")
	}
}
func TestCalculate(t *testing.T) {
	params := RunoutParams{
		ContractStartDate:  time.Date(2022, 1, 14, 0, 0, 0, 0, time.UTC),
		ContractEndDate:    time.Date(2034, 2, 14, 23, 59, 59, 0, time.UTC),
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
	if firstPeriod.NumOfDays != 352 {
		t.Errorf("First period number of days incorrect. Expected 352, got %d", firstPeriod.NumOfDays)
	}
	if !almostEqual(firstPeriod.RateTrend, 1.0, 0.0001) {
		t.Errorf("First period rate trend incorrect. Expected 1.0, got %f", firstPeriod.RateTrend)
	}
	if !almostEqual(firstPeriod.TotalFHRevenue, 225526.8821917808, 0.01) {
		t.Errorf("First period total FH revenue incorrect. Expected 225526.8821917808, got %f", firstPeriod.TotalFHRevenue)
	}

	// Check engine data for the first period
	engine1 := firstPeriod.Engines[0]
	if engine1.WarrantyRateDays != 352 {
		t.Errorf("First period, engine 1 warranty rate days incorrect. Expected 352, got %d", engine1.WarrantyRateDays)
	}
	if !almostEqual(engine1.FHRevenue, 112763.4410958904, 0.01) {
		t.Errorf("First period, engine 1 FH revenue incorrect. Expected 112763.4410958904, got %f", engine1.FHRevenue)
	}

	// Check the last period
	lastPeriod := result.Periods[len(result.Periods)-1]
	if !lastPeriod.EndDate.Equal(time.Date(2033, 12, 31, 23, 59, 59, 0, time.UTC)) {
		t.Errorf("Last period end date incorrect. Expected %v, got %v", time.Date(2033, 12, 31, 23, 59, 59, 0, time.UTC), lastPeriod.EndDate)
	}
	if !almostEqual(lastPeriod.RateTrend, 2.51606537898244, 0.0001) {
		t.Errorf("Last period rate trend incorrect. Expected 2.51606537898244, got %f", lastPeriod.RateTrend)
	}
	if !almostEqual(lastPeriod.TotalFHRevenue, 616246.8097341983, 0.01) {
		t.Errorf("Last period total FH revenue incorrect. Expected 616246.8097341983, got %f", lastPeriod.TotalFHRevenue)
	}

	// Check overall totals
	if !almostEqual(result.TotalFHRevenue, 4805005.2476703655, 0.01) {
		t.Errorf("Total FH revenue incorrect. Expected 4805005.2476703655, got %f", result.TotalFHRevenue)
	}
	if !almostEqual(result.MgmtFeeRevenue, 720750.7871505547, 0.01) {
		t.Errorf("Management fee revenue incorrect. Expected 720750.7871505547, got %f", result.MgmtFeeRevenue)
	}
	if !almostEqual(result.AICRevenue, 4805003.2076703645, 0.01) {
		t.Errorf("AIC revenue incorrect. Expected 4805003.2076703645, got %f", result.AICRevenue)
	}
	if !almostEqual(result.TrustLoadRevenue, 4805004.943710365, 0.01) {
		t.Errorf("Trust load revenue incorrect. Expected 4805004.943710365, got %f", result.TrustLoadRevenue)
	}
	if !almostEqual(result.TrustRevenue, -6878044.74086092, 0.01) {
		t.Errorf("Trust revenue incorrect. Expected -6878044.74086092, got %f", result.TrustRevenue)
	}
	if !almostEqual(result.TotalRevenue, 4805005.2476703655, 0.01) {
		t.Errorf("Total revenue incorrect. Expected 4805005.2476703655, got %f", result.TotalRevenue)
	}

	// Check Buy-In and Enrollment Fees
	if !almostEqual(result.BuyIn, params.BuyIn, 0.01) {
		t.Errorf("Buy-In incorrect. Expected %f, got %f", params.BuyIn, result.BuyIn)
	}
	if !almostEqual(result.EnrollmentFees, params.EnrollmentFees, 0.01) {
		t.Errorf("Enrollment Fees incorrect. Expected %f, got %f", params.EnrollmentFees, result.EnrollmentFees)
	}

	// Check CumulativeTotalRevenue
	if !almostEqual(result.CumulativeTotalRevenue, 4805005.2476703655, 0.01) {
		t.Errorf("Cumulative Total Revenue incorrect. Expected 4805005.2476703655, got %f", result.CumulativeTotalRevenue)
	}
}

func TestRunoutCalculator(t *testing.T) {
	params := RunoutParams{
		ContractStartDate:  time.Date(2022, 1, 14, 0, 0, 0, 0, time.UTC),
		ContractEndDate:    time.Date(2034, 2, 14, 23, 59, 59, 0, time.UTC),
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

	calculator := NewRunoutCalculator(params)

	// Test CalculateCumulativeProfit
	profit, err := calculator.CalculateCumulativeProfit()
	if err != nil {
		t.Fatalf("CalculateCumulativeProfit returned an error: %v", err)
	}
	if !almostEqual(profit, 4805005.2476703655, 0.01) {
		t.Errorf("Expected cumulative profit of 4805005.2476703655, got %f", profit)
	}

	// Test GetParams
	retrievedParams, ok := calculator.GetParams().(RunoutParams)
	if !ok {
		t.Fatalf("GetParams did not return RunoutParams")
	}
	if retrievedParams.ContractStartDate != params.ContractStartDate {
		t.Errorf("GetParams returned incorrect ContractStartDate")
	}
	// Add more checks for other fields as needed

	// Test SetParams
	newParams := RunoutParams{
		ContractStartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		ContractEndDate:   time.Date(2033, 12, 31, 23, 59, 59, 0, time.UTC),
		AUHours:           500,
		// ... set other fields as needed
	}
	err = calculator.SetParams(newParams)
	if err != nil {
		t.Fatalf("SetParams returned an error: %v", err)
	}

	updatedParams, ok := calculator.GetParams().(RunoutParams)
	if !ok {
		t.Fatalf("GetParams did not return RunoutParams after SetParams")
	}
	if updatedParams.ContractStartDate != newParams.ContractStartDate {
		t.Errorf("SetParams did not update ContractStartDate correctly")
	}
	// Add more checks for other fields as needed
}
