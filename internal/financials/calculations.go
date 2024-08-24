package financials

import (
	"fmt"
	"math"
)

func CalculateFinancials(rate float64, params FinancialParams) (float64, error) {
	var cumulativeProfit float64

	for year := 1; year <= params.NumYears; year++ {
		tsn := params.InitialTSN + params.AuHours*float64(year)
		escalatedRate := rate * math.Pow(1+params.RateEscalation/100, float64(year-1))

		engineRevenue := params.AuHours * escalatedRate
		aicRevenue := engineRevenue * params.AIC / 100
		totalRevenue := engineRevenue + aicRevenue

		hsi := tsn >= params.HSITSN && (year == 1 || tsn-params.AuHours < params.HSITSN)
		overhaul := tsn >= params.OverhaulTSN && (year == 1 || tsn-params.AuHours < params.OverhaulTSN)

		hsiCost := 0.0
		if hsi {
			hsiCost = params.HSICost
		}
		overhaulCost := 0.0
		if overhaul {
			overhaulCost = params.OverhaulCost
		}
		totalCost := hsiCost + overhaulCost
		totalProfit := totalRevenue - totalCost
		cumulativeProfit += totalProfit
	}

	return cumulativeProfit, nil
}

func GoalSeek(targetProfit float64, params FinancialParams, initialGuess float64) (float64, int, error) {
	objective := func(rate float64) (float64, error) {
		profit, err := CalculateFinancials(rate, params)
		if err != nil {
			return 0, err
		}
		return profit - targetProfit, nil
	}

	derivative := func(rate float64) (float64, error) {
		epsilon := 1e-6
		f1, err1 := objective(rate + epsilon)
		f2, err2 := objective(rate)
		if err1 != nil || err2 != nil {
			return 0, fmt.Errorf("error calculating derivative")
		}
		return (f1 - f2) / epsilon, nil
	}

	return NewtonRaphson(objective, derivative, initialGuess, 1e-8, 100)
}

func NewtonRaphson(f, df func(float64) (float64, error), x0, xtol float64, maxIter int) (float64, int, error) {
	for i := 0; i < maxIter; i++ {
		fx, err := f(x0)
		if err != nil {
			return 0, i, err
		}
		if math.Abs(fx) < xtol {
			return x0, i + 1, nil
		}

		dfx, err := df(x0)
		if err != nil {
			return 0, i, err
		}
		if dfx == 0 {
			return 0, i, fmt.Errorf("derivative is zero, can't proceed with Newton-Raphson")
		}

		x0 = x0 - fx/dfx
	}
	return 0, maxIter, fmt.Errorf("Newton-Raphson method did not converge within %d iterations", maxIter)
}
