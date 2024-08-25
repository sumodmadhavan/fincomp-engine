package runout

import (
	"fmt"
	"time"
)

type EngineParams struct {
	WarrantyExpDate         time.Time `json:"warrantyExpDate"`
	WarrantyExpHours        float64   `json:"warrantyExpHours"`
	FirstRunRateSwitchDate  time.Time `json:"firstRunRateSwitchDate"`
	SecondRunRateSwitchDate time.Time `json:"secondRunRateSwitchDate"`
	ThirdRunRateSwitchDate  time.Time `json:"thirdRunRateSwitchDate"`
}

type RunoutParams struct {
	ContractStartDate  time.Time      `json:"contractStartDate"`
	ContractEndDate    time.Time      `json:"contractEndDate"`
	AUHours            float64        `json:"auHours"`
	WarrantyRate       float64        `json:"warrantyRate"`
	FirstRunRate       float64        `json:"firstRunRate"`
	SecondRunRate      float64        `json:"secondRunRate"`
	ThirdRunRate       float64        `json:"thirdRunRate"`
	ManagementFees     float64        `json:"managementFees"`
	AICFees            float64        `json:"aicFees"`
	TrustLoadFees      float64        `json:"trustLoadFees"`
	BuyIn              float64        `json:"buyIn"`
	RateEscalation     float64        `json:"rateEscalation"`
	FlightHoursMinimum float64        `json:"flightHoursMinimum"`
	NumOfDaysInYear    float64        `json:"numOfDaysInYear"`
	NumOfDaysInMonth   float64        `json:"numOfDaysInMonth"`
	EnrollmentFees     float64        `json:"enrollmentFees"`
	NumEngines         int            `json:"numEngines"`
	EngineParams       []EngineParams `json:"engineParams"`
}

func (p RunoutParams) Validate() error {
	if p.ContractEndDate.Before(p.ContractStartDate) {
		return fmt.Errorf("contract end date must be after start date")
	}
	if p.AUHours <= 0 {
		return fmt.Errorf("AUHours must be positive")
	}
	if p.WarrantyRate < 0 {
		return fmt.Errorf("WarrantyRate cannot be negative")
	}
	if p.FirstRunRate < 0 {
		return fmt.Errorf("FirstRunRate cannot be negative")
	}
	if p.SecondRunRate < 0 {
		return fmt.Errorf("SecondRunRate cannot be negative")
	}
	if p.ThirdRunRate < 0 {
		return fmt.Errorf("ThirdRunRate cannot be negative")
	}
	if p.ManagementFees < 0 || p.ManagementFees > 100 {
		return fmt.Errorf("ManagementFees must be between 0 and 100")
	}
	if p.AICFees < 0 || p.AICFees > 100 {
		return fmt.Errorf("AICFees must be between 0 and 100")
	}
	if p.TrustLoadFees < 0 || p.TrustLoadFees > 100 {
		return fmt.Errorf("TrustLoadFees must be between 0 and 100")
	}
	if p.BuyIn < 0 {
		return fmt.Errorf("BuyIn cannot be negative")
	}
	if p.RateEscalation < 0 {
		return fmt.Errorf("RateEscalation cannot be negative")
	}
	if p.FlightHoursMinimum < 0 {
		return fmt.Errorf("FlightHoursMinimum cannot be negative")
	}
	if p.NumOfDaysInYear <= 0 {
		return fmt.Errorf("NumOfDaysInYear must be positive")
	}
	if p.NumOfDaysInMonth <= 0 {
		return fmt.Errorf("NumOfDaysInMonth must be positive")
	}
	if p.EnrollmentFees < 0 {
		return fmt.Errorf("EnrollmentFees cannot be negative")
	}
	if p.NumEngines <= 0 {
		return fmt.Errorf("NumEngines must be positive")
	}
	if len(p.EngineParams) != p.NumEngines {
		return fmt.Errorf("number of EngineParams must match NumEngines")
	}

	for i, ep := range p.EngineParams {
		if ep.WarrantyExpDate.Before(p.ContractStartDate) {
			return fmt.Errorf("WarrantyExpDate for engine %d must be after ContractStartDate", i+1)
		}
		if ep.WarrantyExpHours < 0 {
			return fmt.Errorf("WarrantyExpHours for engine %d cannot be negative", i+1)
		}
		if ep.FirstRunRateSwitchDate.Before(p.ContractStartDate) {
			return fmt.Errorf("FirstRunRateSwitchDate for engine %d must be after ContractStartDate", i+1)
		}
		if ep.SecondRunRateSwitchDate.Before(ep.FirstRunRateSwitchDate) {
			return fmt.Errorf("SecondRunRateSwitchDate for engine %d must be after FirstRunRateSwitchDate", i+1)
		}
		if ep.ThirdRunRateSwitchDate.Before(ep.SecondRunRateSwitchDate) {
			return fmt.Errorf("ThirdRunRateSwitchDate for engine %d must be after SecondRunRateSwitchDate", i+1)
		}
	}

	return nil
}
