package goalseek

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
		engine := NewGoalSeekCalculator(params)

		// Test Initialize
		if err := engine.Initialize(params); err != nil {
			b.Fatalf("Initialize error: %v", err)
		}

		// Test Validate
		if err := engine.Validate(); err != nil {
			b.Fatalf("Validate error: %v", err)
		}

		// Test Compute
		if err := engine.Compute(); err != nil {
			b.Fatalf("Compute error: %v", err)
		}

		// Test GetResult
		result := engine.GetResult()
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			b.Fatalf("GetResult did not return a map[string]interface{}")
		}
		_ = resultMap
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
			engine := NewGoalSeekCalculator(params)

			// Test Initialize
			if err := engine.Initialize(params); err != nil {
				b.Fatalf("Initialize error: %v", err)
			}

			// Test Validate
			if err := engine.Validate(); err != nil {
				b.Fatalf("Validate error: %v", err)
			}

			// Test Compute
			if err := engine.Compute(); err != nil {
				b.Fatalf("Compute error: %v", err)
			}

			// Test GetResult
			result := engine.GetResult()
			resultMap, ok := result.(map[string]interface{})
			if !ok {
				b.Fatalf("GetResult did not return a map[string]interface{}")
			}
			_ = resultMap
		}
	})
}
