# Chatear Backend

This is the backend service for Chatear, a real-time chat application built with Clean Architecture principles.

## Tech Stack

*   **Go 1.25.3:** The primary programming language
*   **Gin:** A web framework for building APIs
*   **GraphQL:** A query language for your API (using gqlgen)
*   **PostgreSQL:** Primary database
*   **Redis:** Caching and session storage
*   **NATS:** Message queue for real-time features
*   **JWT:** Authentication and authorization
*   **SMTP:** Email notifications

## Architecture

This project follows **Clean Architecture** principles with **Domain-Driven Design (DDD)** and **SOLID** principles:

- **Domain Layer:** Core business logic and entities
- **Application Layer:** Use cases and application services
- **Infrastructure Layer:** External concerns (database, external APIs)
- **Presentation Layer:** API endpoints and GraphQL resolvers
- **Shared Layer:** Common utilities and cross-cutting concerns

For detailed architecture documentation, see [docs/architecture.md](docs/architecture.md).

## Prerequisites

- Go 1.25.3 or later
- PostgreSQL 12 or later
- Redis 6 or later
- NATS Server 2.0 or later

## Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/jefersonprimer/chatear-backend.git
    cd chatear-backend
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Create a `.env` file:**
    ```bash
    cp env.example .env
    ```

4.  **Update the `.env` file with your credentials:**
    - Configure your database connection
    - Set up Redis connection
    - Configure NATS server
    - Set up SMTP for email notifications
    - Generate a secure JWT secret

5.  **Run database migrations:**
    ```bash
    # Run the migration scripts in the migrations/ directory
    # This depends on your migration tool setup
    ```

## Run

### Development Mode

```bash
go run cmd/api/main.go
```

### Production Mode

```bash
go build -o chatear-backend cmd/api/main.go
./chatear-backend
```

## Project Structure

```
├── cmd/                    # Application entry points
│   ├── api/               # Main API server
│   └── worker/            # Background workers
├── internal/              # Feature modules
│   ├── user/             # User management module
│   └── notification/     # Notification module
├── domain/               # Domain layer (entities, repositories)
├── application/          # Application layer (use cases, services)
├── infrastructure/       # Infrastructure layer (database, external services)
├── presentation/         # Presentation layer (HTTP, GraphQL)
├── shared/              # Shared utilities and constants
├── docs/                # Architecture and API documentation
├── migrations/          # Database migrations
└── env.example         # Environment variables template
```

## API Documentation

- **GraphQL API:** Available at `/graphql` endpoint
- **Health Check:** Available at `/health` endpoint
- **API Documentation:** See [docs/graphql_api.md](docs/graphql_api.md)

## Development

### Running Tests

```bash
go test ./...
```

### Code Generation

```bash
# Generate GraphQL code
go generate ./...

# Or run gqlgen directly
go run github.com/99designs/gqlgen generate
```

### Database Migrations

Database migrations are located in the `migrations/` directory. Run them using your preferred migration tool.

## Contributing

1. Follow Clean Architecture principles
2. Write tests for new features
3. Update documentation as needed
4. Follow Go coding standards

## License

This project is licensed under the MIT License.
