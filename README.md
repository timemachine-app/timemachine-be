# TimeMachine Service

## Overview

TimeMachine is a basic open-source service designed to process data by calling a language model (LLM) and returning the response. It does not log any data, ensuring user privacy and trust. This repository is open source to demonstrate our commitment to transparency and to gain user trust by showing that we do not log any data.

## Repository Structure

The project structure is as follows:

TIMEMACHINE-BE/
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   ├── eventHandler.go
│   │   └── healthHandler.go
├── openai/
│   └── client.go
├── util/
│   ├── openai.go
│   └── ratelimit.go
├── app.yaml
├── default.yaml
├── go.mod
├── go.sum
└── main.go

## Getting Started

### Prerequisites

- Go 1.22.0 or higher

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/timemachine.git
cd timemachine
```

2. Install the dependencies:

```bash
go mod tidy
```

### Configuration

1. Update the configuration files (`app.yaml`, `default.yaml`) with your specific settings.

### Running the Service

To start the service, run:

```bash
go run main.go
```

## Code Structure

- `main.go`: The entry point of the application.
- `internal/config/config.go`: Contains configuration-related code.
- `internal/handlers/eventHandler.go`: Handles event processing by calling the OpenAI LLM.
- `internal/handlers/healthHandler.go`: Provides a health check endpoint.
- `openai/client.go`: Manages interactions with the OpenAI API.
- `util/openai.go`: Contains utility functions related to OpenAI.
- `util/ratelimit.go`: Implements rate limiting middleware.

## Endpoints

### Health Check

- **Endpoint**: `/health`
- **Method**: `GET`
- **Description**: Checks if the service is running and healthy.
- **Response**: 
  ```json
  {
      "status": "healthy"
  }
  ```

### Event Processing

- **Endpoint**: `/event`
- **Method**: `POST`
- **Description**: Processes an event by calling the OpenAI LLM and returns the response.
- **Request Body**: JSON object containing event data.
- **Response**: JSON object containing the response from OpenAI.

## Example

To test the health check endpoint, you can use `curl`:

```bash
curl http://localhost:8080/health
```

To test the event processing endpoint, you can use `curl`:

```bash
curl -X POST http://localhost:8080/event -H "Content-Type: application/json" -d '{"your":"event data"}'
```

## Contributing

We welcome contributions! Please read our [contributing guidelines](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
