package runout

import (
	"fmt"
	"testing"
	"time"
)

func getTestParams() RunoutParams {
	return RunoutParams{
		ContractStartDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		ContractEndDate:    time.Date(2034, 12, 31, 23, 59, 59, 0, time.UTC),
		AUHours:            480,
		WarrantyRate:       243.6,
		FirstRunRate:       255.13,
		SecondRunRate:      255.13,
		ThirdRunRate:       255.13,
		ManagementFees:     15.0,
		AICFees:            20.0,
		TrustLoadFees:      2.98,
		BuyIn:              1352291,
		RateEscalation:     8.75,
		FlightHoursMinimum: 150,
		NumOfDaysInYear:    365,
		NumOfDaysInMonth:   30,
		EnrollmentFees:     25000,
		NumEngines:         2,
		EngineParams: []EngineParams{
			{
				WarrantyExpDate:         time.Date(2025, 10, 31, 23, 59, 59, 0, time.UTC),
				WarrantyExpHours:        1000,
				FirstRunRateSwitchDate:  time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
				SecondRunRateSwitchDate: time.Date(2027, 5, 1, 0, 0, 0, 0, time.UTC),
				ThirdRunRateSwitchDate:  time.Date(2028, 7, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				WarrantyExpDate:         time.Date(2025, 10, 31, 23, 59, 59, 0, time.UTC),
				WarrantyExpHours:        1000,
				FirstRunRateSwitchDate:  time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
				SecondRunRateSwitchDate: time.Date(2027, 5, 1, 0, 0, 0, 0, time.UTC),
				ThirdRunRateSwitchDate:  time.Date(2028, 7, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
}

func BenchmarkCalculate(b *testing.B) {
	params := getTestParams()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Calculate(params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalculateParallel(b *testing.B) {
	params := getTestParams()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Calculate(params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCalculateWithVaryingContractLengths(b *testing.B) {
	baseParams := getTestParams()
	contractLengths := []int{1, 5, 10, 15, 20} // years

	for _, years := range contractLengths {
		b.Run(fmt.Sprintf("ContractLength_%dYears", years), func(b *testing.B) {
			params := baseParams
			params.ContractEndDate = params.ContractStartDate.AddDate(years, 0, 0)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := Calculate(params)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkCalculateWithVaryingEngineNumbers(b *testing.B) {
	baseParams := getTestParams()
	engineNumbers := []int{1, 2, 4, 8, 16}

	for _, numEngines := range engineNumbers {
		b.Run(fmt.Sprintf("Engines_%d", numEngines), func(b *testing.B) {
			params := baseParams
			params.NumEngines = numEngines
			params.EngineParams = make([]EngineParams, numEngines)
			for i := 0; i < numEngines; i++ {
				params.EngineParams[i] = baseParams.EngineParams[0] // Use the same params for all engines
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, err := Calculate(params)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
