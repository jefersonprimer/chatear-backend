# Project Structure Overview

## Root Level Structure

```
chatear-backend/
├── cmd/                          # Application entry points
│   ├── api/                     # Main API server
│   │   └── main.go             # API server entry point
│   └── worker/                  # Background workers
│       ├── notification_worker.go
│       └── user_delete_worker.go
├── internal/                     # Feature modules (existing)
│   ├── user/                    # User management module
│   │   ├── domain/             # User domain logic
│   │   ├── application/        # User use cases
│   │   ├── infrastructure/     # User infrastructure
│   │   └── presentation/       # User presentation layer
│   └── notification/            # Notification module
│       ├── domain/             # Notification domain logic
│       ├── application/        # Notification use cases
│       ├── infrastructure/     # Notification infrastructure
│       └── worker/             # Notification workers
├── domain/                      # Clean Architecture - Domain Layer
│   ├── entities/               # Domain entities
│   │   └── user.go            # User entity example
│   ├── value_objects/          # Value objects
│   ├── repositories/           # Repository interfaces
│   │   └── user_repository.go # User repository interface
│   └── services/               # Domain services
├── application/                 # Clean Architecture - Application Layer
│   ├── usecases/              # Use case implementations
│   │   └── user_usecases.go   # User use cases example
│   ├── services/              # Application services
│   └── dtos/                  # Data Transfer Objects
├── infrastructure/             # Clean Architecture - Infrastructure Layer
│   ├── database/              # Database implementations
│   │   └── postgres_user_repository.go
│   ├── cache/                 # Cache implementations
│   ├── messaging/             # Message queue implementations
│   └── external/              # External service integrations
├── presentation/               # Clean Architecture - Presentation Layer
│   ├── http/                  # HTTP handlers
│   │   └── user_handlers.go   # User HTTP handlers example
│   ├── graphql/               # GraphQL resolvers
│   └── middleware/            # HTTP middleware
├── shared/                     # Clean Architecture - Shared Layer
│   ├── auth/                  # Authentication utilities
│   ├── constants/             # Application constants
│   │   └── app_constants.go
│   ├── errors/                # Error handling
│   │   └── app_errors.go
│   └── util/                  # Utility functions
├── docs/                       # Documentation
│   ├── architecture.md        # Architecture documentation
│   ├── database.md            # Database documentation
│   ├── graphql_api.md         # GraphQL API documentation
│   ├── notification_domain.md # Notification domain docs
│   ├── project_structure.md   # This file
│   ├── security.md            # Security documentation
│   ├── user_domain.md         # User domain documentation
│   └── workers.md             # Workers documentation
├── migrations/                 # Database migrations
│   └── postgres/              # PostgreSQL migrations
├── graph/                      # GraphQL schema and generated code
│   ├── model/                 # Generated GraphQL models
│   ├── resolver.go            # GraphQL resolvers
│   ├── schema.graphqls        # GraphQL schema
│   └── schema.resolvers.go    # Generated resolvers
├── go.mod                     # Go module file
├── go.sum                     # Go module checksums
├── gqlgen.yml                 # GraphQL code generation config
├── env.example                # Environment variables template
└── README.md                  # Project documentation
```

## Architecture Layers Explained

### 1. Domain Layer (`domain/`)
- **Purpose:** Contains the core business logic and rules
- **Dependencies:** None (pure business logic)
- **Contains:** Entities, Value Objects, Repository Interfaces, Domain Services

### 2. Application Layer (`application/`)
- **Purpose:** Orchestrates domain objects and implements use cases
- **Dependencies:** Domain layer only
- **Contains:** Use Cases, Application Services, DTOs

### 3. Infrastructure Layer (`infrastructure/`)
- **Purpose:** Implements external concerns (database, external APIs)
- **Dependencies:** Domain and Application layers
- **Contains:** Database implementations, External service integrations

### 4. Presentation Layer (`presentation/`)
- **Purpose:** Handles user interface and API endpoints
- **Dependencies:** Application layer
- **Contains:** HTTP handlers, GraphQL resolvers, Middleware

### 5. Shared Layer (`shared/`)
- **Purpose:** Common utilities used across all layers
- **Dependencies:** None (or minimal)
- **Contains:** Utilities, Constants, Common interfaces

## Module Organization

The project uses a **modular monolith** approach where each feature is organized as a self-contained module within the `internal/` directory. Each module follows the same Clean Architecture pattern:

```
internal/{module}/
├── domain/        # Module-specific domain logic
├── application/   # Module-specific use cases
├── infrastructure/ # Module-specific infrastructure
└── presentation/  # Module-specific presentation
```

This approach provides:
- **Modularity:** Each feature is self-contained
- **Scalability:** Easy to extract modules to microservices later
- **Maintainability:** Clear separation of concerns
- **Testability:** Each layer can be tested independently

## Key Files

- **`cmd/api/main.go`:** Application entry point
- **`go.mod`:** Go module dependencies
- **`gqlgen.yml`:** GraphQL code generation configuration
- **`env.example`:** Environment variables template
- **`README.md`:** Project setup and usage instructions