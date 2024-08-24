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

The API will be available at `http://localhost:8080`.

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
