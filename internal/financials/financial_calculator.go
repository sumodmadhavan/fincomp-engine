// File: internal/financials/financial_calculator.go

package financials

import "time"

// FinancialCalculator defines the interface for financial calculations
type FinancialCalculator interface {
	CalculateCumulativeProfit() (float64, error)
	GetParams() interface{}
	SetParams(params interface{}) error
}

// CommonParams contains shared parameters between goalseek and runout
type CommonParams struct {
	ContractStartDate time.Time
	ContractEndDate   time.Time
	AUHours           float64
	RateEscalation    float64
}
