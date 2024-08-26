// File: internal/goalseek/goalseek.go

package goalseek

import (
	"financialapi/internal/financials"
	"fmt"
)

// Ensure GoalSeek implements ComputeEngine
var _ financials.ComputeEngine = (*GoalSeek)(nil)

type GoalSeek struct {
	Params financials.FinancialParams
	result map[string]interface{}
}

func (gs *GoalSeek) Initialize(params interface{}) error {
	if p, ok := params.(financials.FinancialParams); ok {
		gs.Params = p
		return nil
	}
	return fmt.Errorf("invalid params type for GoalSeek")
}

func (gs *GoalSeek) Validate() error {
	return gs.Params.Validate()
}

func (gs *GoalSeek) Compute() error {
	optimalRate, iterations, err := financials.GoalSeek(gs.Params.TargetProfit, gs.Params, gs.Params.InitialRate)
	if err != nil {
		return err
	}

	finalCumulativeProfit, err := financials.CalculateFinancials(optimalRate, gs.Params)
	if err != nil {
		return err
	}

	gs.result = map[string]interface{}{
		"optimalWarrantyRate":   optimalRate,
		"iterations":            iterations,
		"finalCumulativeProfit": finalCumulativeProfit,
	}

	return nil
}

func (gs *GoalSeek) GetResult() interface{} {
	return gs.result
}

// Existing methods are kept for backwards compatibility
func (gs *GoalSeek) CalculateCumulativeProfit() (float64, error) {
	optimalRate, _, err := financials.GoalSeek(gs.Params.TargetProfit, gs.Params, gs.Params.InitialRate)
	if err != nil {
		return 0, err
	}

	return financials.CalculateFinancials(optimalRate, gs.Params)
}

func (gs *GoalSeek) GetParams() interface{} {
	return gs.Params
}

func (gs *GoalSeek) SetParams(params interface{}) error {
	return gs.Initialize(params)
}

// NewGoalSeekCalculator creates a new GoalSeek instance
func NewGoalSeekCalculator(params financials.FinancialParams) financials.FinancialCalculator {
	gs := &GoalSeek{}
	gs.Initialize(params)
	return gs
}

// Calculate function remains the same for backwards compatibility
func Calculate(params financials.FinancialParams) (map[string]interface{}, error) {
	gs := &GoalSeek{}
	err := gs.Initialize(params)
	if err != nil {
		return nil, err
	}

	err = gs.Compute()
	if err != nil {
		return nil, err
	}

	return gs.result, nil
}