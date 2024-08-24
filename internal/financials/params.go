package financials

import "fmt"

type FinancialParams struct {
	NumYears       int     `json:"numYears"`
	AuHours        float64 `json:"auHours"`
	InitialTSN     float64 `json:"initialTSN"`
	RateEscalation float64 `json:"rateEscalation"`
	AIC            float64 `json:"aic"`
	HSITSN         float64 `json:"hsitsn"`
	OverhaulTSN    float64 `json:"overhaulTSN"`
	HSICost        float64 `json:"hsiCost"`
	OverhaulCost   float64 `json:"overhaulCost"`
	TargetProfit   float64 `json:"targetProfit"`
	InitialRate    float64 `json:"initialRate"`
}

func (p FinancialParams) Validate() error {
	if p.NumYears <= 0 {
		return fmt.Errorf("NumYears must be positive")
	}
	if p.AuHours <= 0 {
		return fmt.Errorf("AuHours must be positive")
	}
	if p.InitialTSN < 0 {
		return fmt.Errorf("InitialTSN cannot be negative")
	}
	if p.RateEscalation < 0 {
		return fmt.Errorf("RateEscalation cannot be negative")
	}
	if p.AIC < 0 || p.AIC > 100 {
		return fmt.Errorf("AIC must be between 0 and 100")
	}
	if p.HSITSN <= 0 {
		return fmt.Errorf("HSITSN must be positive")
	}
	if p.OverhaulTSN <= 0 {
		return fmt.Errorf("OverhaulTSN must be positive")
	}
	if p.HSICost < 0 {
		return fmt.Errorf("HSICost cannot be negative")
	}
	if p.OverhaulCost < 0 {
		return fmt.Errorf("OverhaulCost cannot be negative")
	}
	if p.TargetProfit <= 0 {
		return fmt.Errorf("TargetProfit must be positive")
	}
	if p.InitialRate <= 0 {
		return fmt.Errorf("InitialRate must be positive")
	}
	return nil
}
