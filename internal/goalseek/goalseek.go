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

// NewGoalSeekCalculator creates a new GoalSeek instance
func NewGoalSeekCalculator(params financials.FinancialParams) financials.ComputeEngine {
	gs := &GoalSeek{}
	gs.Initialize(params)
	return gs
}

