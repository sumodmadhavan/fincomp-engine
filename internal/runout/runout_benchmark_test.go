package runout

import (
	"testing"
	"time"
)

func BenchmarkCalculate(b *testing.B) {
	params := RunoutParams{
		ContractStartDate:  time.Date(2022, 1, 14, 0, 0, 0, 0, time.UTC),
		ContractEndDate:    time.Date(2034, 2, 14, 0, 0, 0, 0, time.UTC),
		AUHours:            480,
		WarrantyRate:       243.6,
		FirstRunRate:       255.13,
		SecondRunRate:      255.13,
		ThirdRunRate:       255.13,
		ManagementFees:     15.0,
		AICFees:            20.0,
		TrustLoadFees:      2.98,
		BuyIn:              1352291.05,
		RateEscalation:     8.75,
		FlightHoursMinimum: 150,
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
	params := RunoutParams{
		ContractStartDate:  time.Date(2022, 1, 14, 0, 0, 0, 0, time.UTC),
		ContractEndDate:    time.Date(2034, 2, 14, 0, 0, 0, 0, time.UTC),
		AUHours:            480,
		WarrantyRate:       243.6,
		FirstRunRate:       255.13,
		SecondRunRate:      255.13,
		ThirdRunRate:       255.13,
		ManagementFees:     15.0,
		AICFees:            20.0,
		TrustLoadFees:      2.98,
		BuyIn:              1352291.05,
		RateEscalation:     8.75,
		FlightHoursMinimum: 150,
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
