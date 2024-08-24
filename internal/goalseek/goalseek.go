package goalseek

import (
	"financialapi/internal/financials"
)

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
