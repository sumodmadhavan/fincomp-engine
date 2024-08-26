// File: internal/financials/compute_engine.go

package financials

type ComputeEngine interface {
    Initialize(params interface{}) error
    Validate() error
    Compute() error
    GetResult() interface{}
}