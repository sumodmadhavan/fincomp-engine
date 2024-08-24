package runout

import (
	"math"
	"time"
)

type EngineData struct {
	WarrantyRateDays  int
	FirstRunRateDays  int
	SecondRunRateDays int
	ThirdRunRateDays  int
	TotalDays         int
	FHUtilization     float64
	FHRevenue         float64
	WarrantyCalc      float64
	FirstRunRateCalc  float64
	SecondRunRateCalc float64
	ThirdRunRateCalc  float64
	Rates             float64
	EscalatedRate     float64
	Shortfall         float64
}

type ContractPeriod struct {
	StartDate          time.Time
	EndDate            time.Time
	NumOfDays          int
	RunoutStartDate    time.Time
	RunoutEndDate      time.Time
	NumOfRunoutDays    int
	ContractYearNumber int
	RateTrend          float64
	Engines            []EngineData
	TotalFHRevenue     float64
}

type RunoutResult struct {
	Periods          []ContractPeriod
	TotalFHRevenue   float64
	MgmtFeeRevenue   float64
	AICRevenue       float64
	TrustLoadRevenue float64
	TrustRevenue     float64
	TotalRevenue     float64
	EnrollmentFees   float64
	BuyIn            float64
}

func Calculate(params RunoutParams) (RunoutResult, error) {
	if err := params.Validate(); err != nil {
		return RunoutResult{}, err
	}

	periods := calculateContractPeriods(params.ContractStartDate, params.ContractEndDate)

	result := RunoutResult{
		Periods:        periods,
		EnrollmentFees: params.EnrollmentFees,
		BuyIn:          params.BuyIn,
	}

	for i := range result.Periods {
		result.Periods[i].Engines = make([]EngineData, params.NumEngines)
		calculatePeriodDetails(&result.Periods[i], params, i+1)
		for e := 0; e < params.NumEngines; e++ {
			calculateEngineRevenue(&result.Periods[i], params, e)
		}
		result.Periods[i].TotalFHRevenue = sumEngineFHRevenue(result.Periods[i].Engines)
		result.TotalFHRevenue += result.Periods[i].TotalFHRevenue
	}

	calculateTotalRevenues(&result, params)

	return result, nil
}

func calculateContractPeriods(startDate, endDate time.Time) []ContractPeriod {
	periods := []ContractPeriod{}
	currentDate := startDate
	yearNumber := 1

	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		var periodEnd time.Time
		if currentDate.Day() <= 14 {
			periodEnd = time.Date(currentDate.Year()+1, currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location()).Add(-time.Second)
		} else {
			periodEnd = time.Date(currentDate.Year()+1, currentDate.Month()+1, 1, 0, 0, 0, 0, currentDate.Location()).Add(-time.Second)
		}

		if periodEnd.After(endDate) {
			periodEnd = endDate
		}

		runoutStart := periodEnd.AddDate(-1, 0, 1)
		if runoutStart.Before(currentDate) {
			runoutStart = currentDate
		}

		period := ContractPeriod{
			StartDate:          currentDate,
			EndDate:            periodEnd,
			NumOfDays:          int(periodEnd.Sub(currentDate).Hours()/24) + 1,
			RunoutStartDate:    runoutStart,
			RunoutEndDate:      periodEnd,
			NumOfRunoutDays:    int(periodEnd.Sub(runoutStart).Hours()/24) + 1,
			ContractYearNumber: yearNumber,
		}

		periods = append(periods, period)
		currentDate = periodEnd.AddDate(0, 0, 1)
		yearNumber++
	}

	return periods
}

func calculatePeriodDetails(period *ContractPeriod, params RunoutParams, yearNumber int) {
	period.RateTrend = math.Pow(1+params.RateEscalation/100, float64(yearNumber-1))

	for e := range period.Engines {
		engine := &period.Engines[e]
		engineParams := params.EngineParams[e]

		engine.WarrantyRateDays = calculateDaysWithinPeriod(period.RunoutStartDate, engineParams.WarrantyExpDate, period.RunoutStartDate, period.RunoutEndDate)
		engine.FirstRunRateDays = calculateDaysWithinPeriod(engineParams.WarrantyExpDate.AddDate(0, 0, 1), engineParams.FirstRunRateSwitchDate, period.RunoutStartDate, period.RunoutEndDate)
		engine.SecondRunRateDays = calculateDaysWithinPeriod(engineParams.FirstRunRateSwitchDate.AddDate(0, 0, 1), engineParams.SecondRunRateSwitchDate, period.RunoutStartDate, period.RunoutEndDate)

		if period.RunoutEndDate.After(engineParams.ThirdRunRateSwitchDate) {
			engine.ThirdRunRateDays = calculateDaysWithinPeriod(engineParams.SecondRunRateSwitchDate.AddDate(0, 0, 1), period.RunoutEndDate, period.RunoutStartDate, period.RunoutEndDate)
		} else {
			engine.ThirdRunRateDays = calculateDaysWithinPeriod(engineParams.SecondRunRateSwitchDate.AddDate(0, 0, 1), engineParams.ThirdRunRateSwitchDate, period.RunoutStartDate, period.RunoutEndDate)
		}

		engine.TotalDays = engine.WarrantyRateDays + engine.FirstRunRateDays + engine.SecondRunRateDays + engine.ThirdRunRateDays
	}
}
func calculateEngineRevenue(period *ContractPeriod, params RunoutParams, engineIndex int) {
	engine := &period.Engines[engineIndex]

	engine.WarrantyCalc = float64(engine.WarrantyRateDays) * params.WarrantyRate
	engine.FirstRunRateCalc = float64(engine.FirstRunRateDays) * params.FirstRunRate
	engine.SecondRunRateCalc = float64(engine.SecondRunRateDays) * params.SecondRunRate
	engine.ThirdRunRateCalc = float64(engine.ThirdRunRateDays) * params.ThirdRunRate

	engine.Rates = engine.WarrantyCalc + engine.FirstRunRateCalc + engine.SecondRunRateCalc + engine.ThirdRunRateCalc
	engine.EscalatedRate = engine.Rates * period.RateTrend
	engine.FHUtilization = params.AUHours / params.NumOfDaysInYear * float64(engine.TotalDays)

	if engine.FHUtilization < params.FlightHoursMinimum {
		engine.Shortfall = params.FlightHoursMinimum - engine.FHUtilization
	}

	engine.FHRevenue = engine.EscalatedRate * (params.AUHours / params.NumOfDaysInYear)
}

func calculateTotalRevenues(result *RunoutResult, params RunoutParams) {
	result.MgmtFeeRevenue = result.TotalFHRevenue * (params.ManagementFees / 100)
	result.AICRevenue = result.TotalFHRevenue * (1 - params.ManagementFees/100) * (params.AICFees / 100)
	result.TrustLoadRevenue = result.TotalFHRevenue * (1 - params.ManagementFees/100) * (params.TrustLoadFees / 100)

	result.TrustRevenue = result.TotalFHRevenue - (result.MgmtFeeRevenue + result.AICRevenue + result.TrustLoadRevenue)

	// Add BuyIn to the first period's FHRevenue and TrustRevenue
	if len(result.Periods) > 0 {
		result.Periods[0].TotalFHRevenue += params.BuyIn
		result.TotalFHRevenue += params.BuyIn
		result.TrustRevenue += params.BuyIn
	}

	result.TotalRevenue = result.MgmtFeeRevenue + result.AICRevenue + result.TrustLoadRevenue + result.TrustRevenue
}

func calculateDaysWithinPeriod(start, end, periodStart, periodEnd time.Time) int {
	if start.After(periodEnd) || end.Before(periodStart) {
		return 0
	}
	actualStart := max(start, periodStart)
	actualEnd := min(end, periodEnd)
	return int(actualEnd.Sub(actualStart).Hours()/24) + 1
}

func sumEngineFHRevenue(engines []EngineData) float64 {
	total := 0.0
	for _, engine := range engines {
		total += engine.FHRevenue
	}
	return total
}

func max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
