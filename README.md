# Financial API

This project is a REST API for financial calculations, specifically focused on GoalSeek and Runout functionalities. It uses the Newton-Raphson method for numerical computations.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Building and Running](#building-and-running)
  - [Using Go](#using-go)
  - [Using Docker](#using-docker)
- [API Endpoints](#api-endpoints)
- [Calculation Engine](#calculation-engine)
- [Sample Requests and Responses](#sample-requests-and-responses)

## Prerequisites

- Go 1.21 or later
- Docker (optional, for containerized deployment)

## Core Structure 
![image](https://github.com/user-attachments/assets/69f6ef9d-6746-48e4-b46a-8e0e994b2a37)


## Building and Running

### Using Go

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/financialapi.git
   cd financialapi
   ```

2. Build the project:
   ```
   go build ./cmd/server
   ```

3. Run the server:
   ```
   ./server
   ```

The API will be available at `http://localhost:8080`.

### Using Docker

1. Build the Docker image:
   ```
   docker build -t financialapi .
   ```

2. Run the container:
   ```
   docker run -p 8080:8080 financialapi
   ```

### Using Docker with Benchmark

You can run tests and benchmarks inside a Docker container without installing Go on your local machine.

1. First, ensure you have built the Docker image as described in the "Building and Running" section.

2. To run all tests:
   ```
   docker run --rm financialapi go test ./...
   ```

3. To run tests with coverage:
   ```
   docker run --rm financialapi sh -c "go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out"
   ```

4. To run all benchmarks:
   ```
   docker run --rm financialapi go test -bench=. ./...
   ```

5. To run benchmarks with memory allocation statistics:
   ```
   docker run --rm financialapi go test -bench=. -benchmem ./...
   ```

6. To run tests or benchmarks for a specific package (e.g., goalseek):
   ```
   docker run --rm financialapi go test -bench=. ./internal/goalseek
   ```

Note: When running tests or benchmarks in Docker, you won't be able to generate an HTML coverage report or use the Go profiling tools directly. If you need these features, consider running the tests locally with Go installed.

The API will be available at `http://localhost:8080`.
## Testing and Benchmarking

### Running Tests

To run all tests in the project:

```
go test ./...
```

To run tests for a specific package (e.g., goalseek):

```
go test ./internal/goalseek
```

To run tests with coverage:

```
go test -cover ./...
```

For a detailed HTML coverage report:

```
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Running Benchmarks

To run all benchmarks:

```
go test -bench=. ./...
```

To run benchmarks for a specific package (e.g., goalseek):

```
go test -bench=. ./internal/goalseek
```

To run benchmarks with memory allocation statistics:

```
go test -bench=. -benchmem ./...
```

For more detailed profiling, you can use the `-cpuprofile` and `-memprofile` flags:

```
go test -bench=. -cpuprofile=cpu.out -memprofile=mem.out ./...
```

You can then analyze these profiles using:

```
go tool pprof cpu.out
go tool pprof mem.out
```

## API Endpoints

- POST `/goalseek`: Performs GoalSeek calculation
- POST `/runout`: Performs Runout calculation (under development)

## Calculation Engine

The calculation engine uses the Newton-Raphson method for numerical computations. This method is used to find roots of a function, which in our case, helps in finding the optimal warranty rate for a given target profit.

The Newton-Raphson method works as follows:

1. Start with an initial guess for the root.
2. Calculate the function value and its derivative at this point.
3. Use these values to calculate a better approximation of the root.
4. Repeat steps 2-3 until the approximation is close enough to the actual root.

In our GoalSeek implementation, we use this method to find the warranty rate that achieves a target profit. The function we're finding the root of is the difference between the calculated profit and the target profit.

The Newton-Raphson method is a powerful numerical technique for finding roots of real-valued functions. For a detailed mathematical treatment and analysis of the method.

Ben-Israel, A. (2001). Newton's method with modified functions. Contemporary Mathematics, 204, 39-50.

[Refer to this article on ScienceDirect](https://www.sciencedirect.com/science/article/pii/S0377042700004350)

This paper provides insights into the convergence properties and various modifications of the Newton-Raphson method, which can be particularly useful for understanding the theoretical foundations of our implementation.

## GoalSeek and Newton-Raphson Method

### GoalSeek Functionality

The GoalSeek function in our API aims to find the optimal warranty rate that achieves a target profit. It uses the Newton-Raphson method to iteratively approximate the solution.

### Newton-Raphson Method

The Newton-Raphson method is an efficient algorithm for finding roots of a real-valued function. In our case, we use it to find the warranty rate that results in a specific target profit.

The basic formula for the Newton-Raphson method is:

x_{n+1} = x_n - f(x_n) / f'(x_n)

Where:
- x_n is the current approximation
- f(x_n) is the function value at x_n
- f'(x_n) is the derivative of the function at x_n

### Implementation

Here's a simplified version of our Newton-Raphson implementation:

```go
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
```

In this implementation:

1. We define an `objective` function that calculates the difference between the computed profit and the target profit.
2. We approximate the derivative using a finite difference method.
3. The `NewtonRaphson` function implements the iterative process, updating the guess until it converges or reaches the maximum number of iterations.

The GoalSeek functionality uses this method to find the warranty rate that achieves the target profit, making it a powerful tool for financial modeling and decision-making.

For a detailed mathematical treatment of the Newton-Raphson method, refer to:

Ben-Israel, A. (2001). Newton's method with modified functions. Contemporary Mathematics, 204, 39-50.
https://www.sciencedirect.com/science/article/pii/S0377042700004350

## Sample Requests and Responses

### GoalSeek Endpoint

Request:
```json
POST /goalseek
Content-Type: application/json

{
  "numYears": 10,
  "auHours": 450,
  "initialTSN": 100,
  "rateEscalation": 5,
  "aic": 10,
  "hsitsn": 1000,
  "overhaulTSN": 3000,
  "hsiCost": 50000,
  "overhaulCost": 100000,
  "targetProfit": 3000000,
  "initialRate": 320
}
```

Response:
```json
{
  "optimalWarrantyRate": 505.93820432563325,
  "iterations": 3,
  "finalCumulativeProfit": 2999999.9999999986
}
```

### Runout Endpoint (Under Development)

Request:
```json
POST /runout
Content-Type: application/json

{
  "numYears": 10,
  "auHours": 450,
  "initialTSN": 100,
  "rateEscalation": 5,
  "aic": 10,
  "hsitsn": 1000,
  "overhaulTSN": 3000,
  "hsiCost": 50000,
  "overhaulCost": 100000,
  "targetProfit": 3000000,
  "initialRate": 320
}
```

Response:
```json
{
  "message": "Runout calculation is under development"
}
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
