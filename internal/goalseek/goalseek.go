package goalseek

import (
	"financialapi/internal/financials"
	"fmt"
)

type GoalSeek struct {
	Params financials.FinancialParams
}

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
	if p, ok := params.(financials.FinancialParams); ok {
		gs.Params = p
		return nil
	}
	return fmt.Errorf("invalid params type for GoalSeek")
}

// NewGoalSeekCalculator creates a new GoalSeek instance
func NewGoalSeekCalculator(params financials.FinancialParams) financials.FinancialCalculator {
	return &GoalSeek{Params: params}
}

// Calculate function remains the same, but now uses financials.FinancialParams
func Calculate(params financials.FinancialParams) (map[string]interface{}, error) {
	optimalRate, iterations, err := financials.GoalSeek(params.TargetProfit, params, params.InitialRate)
	if err != nil {
		return nil, err
	}

	finalCumulativeProfit, err := financials.CalculateFinancials(optimalRate, params)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"optimalWarrantyRate":   optimalRate,
		"iterations":            iterations,
		"finalCumulativeProfit": finalCumulativeProfit,
	}, nil
}
