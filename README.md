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
