# Chatear Backend

This is the backend service for the Chatear application, built with Go. It provides a robust and scalable foundation for real-time communication features, user management, and notifications.

## Table of Contents

*   [Prerequisites](#prerequisites)
*   [Setup](#setup)
*   [Running the API Server](#running-the-api-server)
*   [Running the Notification Worker](#running-the-notification-worker)
*   [Running the User Delete Worker](#running-the-user-delete-worker)
*   [Running Tests](#running-tests)
*   [Documentation](#documentation)

## Prerequisites

Before you begin, ensure you have the following installed:

*   **Go**: Version 1.25.3 or higher.
*   **Supabase**: A running Supabase instance for data persistence.
*   **Redis**: A running Redis instance for caching, session management, and rate limiting.
*   **NATS**: A running NATS server for inter-service communication and event streaming.

## Setup

1.  **Clone the repository**:

    ```bash
    git clone https://github.com/jefersonprimer/chatear-backend.git
    cd chatear-backend
    ```

2.  **Download Go modules**:

    ```bash
    go mod download
    ```

3.  **Configure environment variables**:

    Copy the example environment file and update it with your specific configurations for database, Redis, NATS, JWT secrets, and SMTP settings.

    ```bash
    cp env.example .env
    # Open .env in your editor and fill in the details
    ```

    Refer to `docs/env.md` for a detailed explanation of each environment variable.

## Running the API Server

The API server exposes the GraphQL and HTTP endpoints for the application.

```bash
go run cmd/api/main.go
```

The server will typically start on the port specified in your `.env` file (default: `8080`).

## Running the Notification Worker

The notification worker processes email sending tasks from the NATS queue.

```bash
go run cmd/worker/notification_worker.go
```

## Running the User Delete Worker

The user delete worker handles asynchronous user account deletion processes.

```bash
go run cmd/worker/user_delete_worker.go
```

## Running Tests

To run all unit and integration tests for the project:

```bash
go test ./...
```

To run tests with verbose output:

```bash
go test -v ./...
```

## Documentation

Detailed documentation for the project's architecture, dependencies, environment variables, and domain-specific workflows can be found in the `docs/` directory:

*   [Architecture Overview](docs/architecture.md)
*   [Project Dependencies](docs/dependencies.md)
*   [Environment Variables](docs/env.md)
*   [User Domain](docs/user_domain.md)
*   [Notification Domain](docs/notification_domain.md)
