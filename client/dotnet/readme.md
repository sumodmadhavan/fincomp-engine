# .NET Client for Go Computation Engine

This project is a .NET client that interacts with a Go-based computation engine through a RESTful API. The client provides pre-processing and scenario handling functionalities, while the Go computation engine performs advanced calculations such as Goal Seek and Run Out.

## Architecture Overview
![Architecture Diagram](https://github.com/user-attachments/assets/280f17ec-c28a-447f-ad2c-4b4306971c54)
- **Pre-Processing/Processing API** (C#): Handles scenarios and output processing. It sends HTTP requests to the Go computation engine with the necessary parameters for calculations.
  
- **Computation Engine** (Go): Receives JSON requests from the .NET client, performs calculations (e.g., Goal Seek, Run Out), and returns the results as JSON responses.

## Features

- **Scenario Handling**: Define and manage different scenarios for financial modeling.
- **Goal Seek**: Utilize the Go engine to find the optimal solution for your financial goals.
- **Run Out**: Predict long-term financial outcomes with high precision.

## Prerequisites

- [.NET SDK](https://dotnet.microsoft.com/download) - Make sure you have the .NET SDK installed on your machine.
- [Go](https://golang.org/dl/) - Ensure the Go runtime is installed and the Go API server is running.

## Setup

### Clone the Repository

```bash
git clone https://github.com/sumodmadhavan/fincomp-engine.git
cd fincomp-engine


