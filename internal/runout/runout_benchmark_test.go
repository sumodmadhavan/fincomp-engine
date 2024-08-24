package runout

import (
	"financialapi/internal/financials"
	"testing"
)

func BenchmarkCalculate(b *testing.B) {
	params := financials.FinancialParams{
		NumYears:       10,
		AuHours:        450,
		InitialTSN:     100,
		RateEscalation: 5,
		AIC:            10,
		HSITSN:         1000,
		OverhaulTSN:    3000,
		HSICost:        50000,
		OverhaulCost:   100000,
		TargetProfit:   3000000,
		InitialRate:    320,
	}

	b.ResetTimer() // Reset the timer to exclude setup time

	for i := 0; i < b.N; i++ {
		_, err := Calculate(params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCalculateParallel tests the Calculate function with parallel execution
func BenchmarkCalculateParallel(b *testing.B) {
	params := financials.FinancialParams{
		NumYears:       10,
		AuHours:        450,
		InitialTSN:     100,
		RateEscalation: 5,
		AIC:            10,
		HSITSN:         1000,
		OverhaulTSN:    3000,
		HSICost:        50000,
		OverhaulCost:   100000,
		TargetProfit:   3000000,
		InitialRate:    320,
	}

	b.ResetTimer() // Reset the timer to exclude setup time

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Calculate(params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
