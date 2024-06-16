# TimeMachine Service

## Overview

TimeMachine is a basic open-source service designed to process data by calling a language model (LLM) and returning the response. It does not log any data, ensuring user privacy and trust. This repository is open source to demonstrate our commitment to transparency and to gain user trust by showing that we do not log any data.

## Getting Started

### Prerequisites

- Go 1.22.0 or higher

### Installation

1. Clone the repository:

```bash
git clone https://github.com/timemachine-app/timemachine-be
cd timemachine-be
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
- **Request Body**: Input Form containing event data.
- **Response**: JSON object containing the processed response from OpenAI.

## Example

To test the health check endpoint, you can use `curl`:

```bash
curl http://localhost:8080/health
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/git/git-scm.com/blob/main/MIT-LICENSE.txt) file for details.