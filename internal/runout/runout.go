package runout

import (
	"financialapi/internal/financials"
	"fmt"
	"time"
)

// Ensure RunoutCalculator implements ComputeEngine
var _ financials.ComputeEngine = (*RunoutCalculator)(nil)

type RunoutCalculator struct {
	Params RunoutParams
	result RunoutResult
}

func (r *RunoutCalculator) Initialize(params interface{}) error {
	if p, ok := params.(RunoutParams); ok {
		r.Params = p
		return nil
	}
	return fmt.Errorf("invalid params type for Runout")
}

func (r *RunoutCalculator) Validate() error {
	return r.Params.Validate()
}

func (r *RunoutCalculator) Compute() error {
	result, err := Calculate(r.Params)
	if err != nil {
		return err
	}
	r.result = result
	return nil
}

func (r *RunoutCalculator) GetResult() interface{} {
	return r.result
}

// Existing methods are kept for backwards compatibility
func (r *RunoutCalculator) CalculateCumulativeProfit() (float64, error) {
	result, err := Calculate(r.Params)
	if err != nil {
		return 0, err
	}
	return result.CumulativeTotalRevenue, nil
}

func (r *RunoutCalculator) GetParams() interface{} {
	return r.Params
}

func (r *RunoutCalculator) SetParams(params interface{}) error {
	return r.Initialize(params)
}

// NewRunoutCalculator creates a new RunoutCalculator instance
func NewRunoutCalculator(params RunoutParams) financials.FinancialCalculator {
	rc := &RunoutCalculator{}
	rc.Initialize(params)
	return rc
}


type EngineData struct {
	EngineID          int
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
	StartDate              time.Time
	EndDate                time.Time
	NumOfDays              int
	RunoutStartDate        time.Time
	RunoutEndDate          time.Time
	NumOfRunoutDays        int
	ContractYearNumber     int
	RateTrend              float64
	Engines                []EngineData
	TotalFHRevenue         float64
	MgmtFeeRevenue         float64
	AICRevenue             float64
	TrustLoadRevenue       float64
	TrustRevenue           float64
	TotalRevenue           float64
	BuyIn                  float64
	CumulativeTotalRevenue float64
}

type RunoutResult struct {
	Periods                []ContractPeriod
	TotalFHRevenue         float64
	MgmtFeeRevenue         float64
	AICRevenue             float64
	TrustLoadRevenue       float64
	TrustRevenue           float64
	TotalRevenue           float64
	EnrollmentFees         float64
	BuyIn                  float64
	CumulativeTotalRevenue float64
}

var rateTrendValues = []float64{
	1, 1.0875, 1.18265625, 1.286138671875, 1.39867580566406, 1.52105993865967,
	1.65415268329239, 1.79889104308047, 1.95629400935001, 2.12746973516814,
	2.31362333699535, 2.51606537898244, 2.73622109964341, 2.97564044586221,
	3.23600898487515, 3.51915977105172, 3.82708625101875, 4.16195629798289,
	4.52612747405639, 4.92216362803633, 5.3528529454895, 5.82122757821983,
	6.33058499131407, 6.88451117805405, 7.48690590613378, 8.14201017292048,
	8.85443606305102, 9.62919921856799, 10.4717541501927, 11.3880326383345,
	12.3844854941888, 13.4681279749303, 14.6465891727367, 15.9281657253512,
	17.3218802263194, 18.8375447461224,
}

var engineValues = []int{1085718, 1085719}

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

	numPeriods := len(periods)
	if numPeriods > len(rateTrendValues) {
		numPeriods = len(rateTrendValues)
	}

	for i := 0; i < numPeriods; i++ {
		result.Periods[i].Engines = make([]EngineData, len(engineValues))
		calculatePeriodDetails(&result.Periods[i], params, i+1, rateTrendValues[i])
		for e, engineValue := range engineValues {
			calculateEngineRevenue(&result.Periods[i], params, e, engineValue)
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
			periodEnd = periodEnd.AddDate(0, 0, -periodEnd.Day())
		}

		if periodEnd.After(endDate) {
			periodEnd = endDate
		}

		runoutStart := periodEnd.AddDate(-1, 0, 1)
		if runoutStart.Before(currentDate) {
			runoutStart = currentDate
		}

		days := int(periodEnd.Sub(currentDate).Hours()/24) + 1
		if days >= 300 {
			period := ContractPeriod{
				StartDate:          currentDate,
				EndDate:            periodEnd,
				NumOfDays:          days,
				RunoutStartDate:    runoutStart,
				RunoutEndDate:      periodEnd,
				NumOfRunoutDays:    int(periodEnd.Sub(runoutStart).Hours()/24) + 1,
				ContractYearNumber: yearNumber,
			}
			periods = append(periods, period)
			yearNumber++
		}

		currentDate = periodEnd.AddDate(0, 0, 1)
	}

	return periods
}

func calculatePeriodDetails(period *ContractPeriod, params RunoutParams, yearNumber int, rateTrend float64) {
	period.RateTrend = rateTrend

	for e, engineValue := range engineValues {
		engine := &period.Engines[e]
		engine.EngineID = engineValue
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

func calculateEngineRevenue(period *ContractPeriod, params RunoutParams, engineIndex, engineValue int) {
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
	cumulativeRevenue := 0.0
	for i := range result.Periods {
		period := &result.Periods[i]
		period.MgmtFeeRevenue = period.TotalFHRevenue * (params.ManagementFees / 100)
		period.AICRevenue = period.TotalFHRevenue - (1-params.ManagementFees/100)*(params.AICFees/100)
		period.TrustLoadRevenue = period.TotalFHRevenue - (1-params.ManagementFees/100)*(params.TrustLoadFees/100)

		if i == 0 {
			period.BuyIn = params.BuyIn
		}

		period.TrustRevenue = period.TotalFHRevenue - (period.MgmtFeeRevenue + period.AICRevenue + period.TrustLoadRevenue + period.BuyIn)
		period.TotalRevenue = period.MgmtFeeRevenue + period.AICRevenue + period.TrustLoadRevenue + period.BuyIn + period.TrustRevenue

		cumulativeRevenue += period.TotalRevenue
		period.CumulativeTotalRevenue = cumulativeRevenue

		result.MgmtFeeRevenue += period.MgmtFeeRevenue
		result.AICRevenue += period.AICRevenue
		result.TrustLoadRevenue += period.TrustLoadRevenue
		result.TrustRevenue += period.TrustRevenue
		result.TotalRevenue += period.TotalRevenue
	}

	result.CumulativeTotalRevenue = cumulativeRevenue
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
